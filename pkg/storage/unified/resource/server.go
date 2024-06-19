package resource

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/grafana/grafana/pkg/apimachinery/identity"
	"github.com/grafana/grafana/pkg/apimachinery/utils"
)

// Package-level errors.
var (
	ErrNotFound                  = errors.New("entity not found")
	ErrOptimisticLockingFailed   = errors.New("optimistic locking failed")
	ErrUserNotFoundInContext     = errors.New("user not found in context")
	ErrUnableToReadResourceJSON  = errors.New("unable to read resource json")
	ErrNextPageTokenNotSupported = errors.New("nextPageToken not yet supported")
	ErrLimitNotSupported         = errors.New("limit not yet supported")
	ErrNotImplementedYet         = errors.New("not implemented yet")
)

// ResourceServer implements all services
type ResourceServer interface {
	ResourceStoreServer
	ResourceSearchServer
	DiagnosticsServer
	LifecycleHooks
}

type AppendingStore interface {
	// Write a Create/Update/Delete,
	// NOTE: the contents of WriteEvent have been validated
	// Return the revisionVersion for this event or error
	WriteEvent(context.Context, WriteEvent) (int64, error)

	// Create new name for a given resource
	GenerateName(ctx context.Context, key *ResourceKey, prefix string) (string, error)

	// Read a value from storage
	Read(context.Context, *ReadRequest) (*ReadResponse, error)

	// Implement List -- this expects the read after write semantics
	List(context.Context, *ListRequest) (*ListResponse, error)

	// Watch for events
	Watch(context.Context, *WatchRequest) (chan *WatchEvent, error)
}

type ResourceServerOptions struct {
	// OTel tracer
	Tracer trace.Tracer

	// When running in a cluster, each node should have a different ID
	// This is used for snowflake generation and log identification
	NodeID int64

	// Get the next EventID.  When not set, this will default to snowflake IDs
	NextEventID func() int64

	// Real storage backend
	Store AppendingStore

	// Real storage backend
	Search ResourceSearchServer

	// Diagnostics
	Diagnostics DiagnosticsServer

	// Check if a user has access to write folders
	// When this is nil, no resources can have folders configured
	WriteAccess WriteAccessHooks

	// Callbacks for startup and shutdown
	Lifecycle LifecycleHooks

	// Get the current time in unix millis
	Now func() int64
}

func NewResourceServer(opts ResourceServerOptions) (ResourceServer, error) {
	if opts.Tracer == nil {
		opts.Tracer = noop.NewTracerProvider().Tracer("resource-server")
	}

	if opts.NextEventID == nil {
		if opts.NodeID == 0 {
			opts.NodeID = rand.Int63n(1024)
		}
		eventNode, err := snowflake.NewNode(opts.NodeID)
		if err != nil {
			return nil, apierrors.NewInternalError(
				fmt.Errorf("error initializing snowflake id generator :: %w", err))
		}
		opts.NextEventID = func() int64 {
			return eventNode.Generate().Int64()
		}
	}

	if opts.Store == nil {
		return nil, fmt.Errorf("missing AppendingStore implementation")
	}
	if opts.Search == nil {
		opts.Search = &noopService{}
	}
	if opts.Diagnostics == nil {
		opts.Search = &noopService{}
	}
	if opts.Now == nil {
		opts.Now = func() int64 {
			return time.Now().UnixMilli()
		}
	}

	return &server{
		tracer:      opts.Tracer,
		nextEventID: opts.NextEventID,
		store:       opts.Store,
		search:      opts.Search,
		diagnostics: opts.Diagnostics,
		access:      opts.WriteAccess,
		lifecycle:   opts.Lifecycle,
		now:         opts.Now,
	}, nil
}

var _ ResourceServer = &server{}

type server struct {
	tracer      trace.Tracer
	nextEventID func() int64
	store       AppendingStore
	search      ResourceSearchServer
	diagnostics DiagnosticsServer
	access      WriteAccessHooks
	lifecycle   LifecycleHooks
	now         func() int64

	// init checking
	once    sync.Once
	initErr error
}

// Init implements ResourceServer.
func (s *server) Init() error {
	s.once.Do(func() {
		// TODO, setup a broadcaster for watch

		// Call lifecycle hooks
		if s.lifecycle != nil {
			err := s.lifecycle.Init()
			if err != nil {
				s.initErr = fmt.Errorf("initialize Resource Server: %w", err)
			}
		}
	})
	return s.initErr
}

func (s *server) Stop() {
	s.initErr = fmt.Errorf("service is stopping")
	if s.lifecycle != nil {
		s.lifecycle.Stop()
	}
	s.initErr = fmt.Errorf("service is stopped")
}

// Old value indicates an update -- otherwise a create
func (s *server) newEventBuilder(ctx context.Context, key *ResourceKey, value, oldValue []byte) (*writeEventBuilder, error) {
	event, err := newEventFromBytes(value, oldValue)
	if err != nil {
		return nil, err
	}
	event.EventID = s.nextEventID()
	event.Key = key
	event.Requester, err = identity.GetRequester(ctx)
	if err != nil {
		return nil, ErrUserNotFoundInContext
	}

	obj := event.Meta
	if key.Name != obj.GetName() {
		return nil, apierrors.NewBadRequest("key/name do not match")
	}
	if key.Namespace != obj.GetNamespace() {
		return nil, apierrors.NewBadRequest("key/namespace do not match")
	}

	gvk := obj.GetGroupVersionKind()
	if gvk.Kind == "" {
		return nil, apierrors.NewBadRequest("expecting resources with a kind in the body")
	}
	if gvk.Version == "" {
		return nil, apierrors.NewBadRequest("expecting resources with an apiVersion")
	}
	if gvk.Group != "" && gvk.Group != key.Group {
		return nil, apierrors.NewBadRequest(
			fmt.Sprintf("group in key does not match group in the body (%s != %s)", key.Group, gvk.Group),
		)
	}

	// This needs to be a create function
	if key.Name == "" {
		prefix := obj.GetGenerateName()
		if prefix == "" {
			return nil, apierrors.NewBadRequest("must have name or generate name set")
		}
		key.Name, err = s.store.GenerateName(ctx, key, prefix)
		if err != nil {
			return nil, err
		}
		obj.SetName(key.Name)
		obj.SetGenerateName("")
	} else if obj.GetGenerateName() != "" {
		return nil, apierrors.NewBadRequest("values with a name must not include generate name")
	}
	err = validateName(obj.GetName())
	if err != nil {
		return nil, err
	}

	if obj.GetName() != key.Name {
		return nil, apierrors.NewBadRequest("key name does not match the name in the body")
	}
	folder := obj.GetFolder()
	if folder != "" {
		err = s.access.CanWriteFolder(ctx, event.Requester, folder)
		if err != nil {
			return nil, err
		}
	}
	origin, err := obj.GetOriginInfo()
	if err != nil {
		return nil, apierrors.NewBadRequest("invalid origin info")
	}
	if origin != nil {
		err = s.access.CanWriteOrigin(ctx, event.Requester, origin.Name)
		if err != nil {
			return nil, err
		}
	}
	obj.SetOriginInfo(origin)

	// Make sure old values do not mutate things they should not
	if event.OldMeta != nil {
		old := event.OldMeta

		if obj.GetUID() != event.OldMeta.GetUID() {
			return nil, apierrors.NewBadRequest(
				fmt.Sprintf("UIDs do not match (old: %s, new: %s)", old.GetUID(), obj.GetUID()))
		}

		// Can not change creation timestamps+user
		if obj.GetCreatedBy() != event.OldMeta.GetCreatedBy() {
			return nil, apierrors.NewBadRequest(
				fmt.Sprintf("created by changed (old: %s, new: %s)", old.GetCreatedBy(), obj.GetCreatedBy()))
		}
		if obj.GetCreationTimestamp() != event.OldMeta.GetCreationTimestamp() {
			return nil, apierrors.NewBadRequest(
				fmt.Sprintf("creation timestamp changed (old:%v, new:%v)", old.GetCreationTimestamp(), obj.GetCreationTimestamp()))
		}
	}
	return event, nil
}

func (s *server) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	ctx, span := s.tracer.Start(ctx, "storage_server.Create")
	defer span.End()

	if err := s.Init(); err != nil {
		return nil, err
	}

	rsp := &CreateResponse{}
	builder, err := s.newEventBuilder(ctx, req.Key, req.Value, nil)
	if err != nil {
		rsp.Status, err = errToStatus(err)
		return rsp, err
	}

	obj := builder.Meta
	obj.SetCreatedBy(builder.Requester.GetUID().String())
	obj.SetUpdatedBy("")
	obj.SetUpdatedTimestamp(nil)
	obj.SetCreationTimestamp(metav1.NewTime(time.UnixMilli(s.now())))
	obj.SetUID(types.UID(uuid.New().String()))

	event, err := builder.toEvent()
	if err != nil {
		rsp.Status, err = errToStatus(err)
		return rsp, err
	}

	rsp.ResourceVersion, err = s.store.WriteEvent(ctx, event)
	if err == nil {
		rsp.Value = event.Value // with mutated fields
	} else {
		rsp.Status, err = errToStatus(err)
	}
	return rsp, err
}

// Convert golang errors to status result errors that can be returned to a client
func errToStatus(err error) (*StatusResult, error) {
	if err != nil {
		apistatus, ok := err.(apierrors.APIStatus)
		if ok {
			s := apistatus.Status()
			return &StatusResult{
				Status:  s.Status,
				Message: s.Message,
				Reason:  string(s.Reason),
				Code:    s.Code,
			}, nil
		}

		// TODO... better conversion!!!
		return &StatusResult{
			Status:  "Failure",
			Message: err.Error(),
			Code:    500,
		}, nil
	}
	return nil, err
}

func (s *server) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	ctx, span := s.tracer.Start(ctx, "storage_server.Update")
	defer span.End()

	if err := s.Init(); err != nil {
		return nil, err
	}

	rsp := &UpdateResponse{}
	if req.ResourceVersion < 0 {
		rsp.Status, _ = errToStatus(apierrors.NewBadRequest("update must include the previous version"))
		return rsp, nil
	}

	latest, err := s.store.Read(ctx, &ReadRequest{
		Key: req.Key,
	})
	if err != nil {
		return nil, err
	}
	if latest.Value == nil {
		return nil, apierrors.NewBadRequest("current value does not exist")
	}

	builder, err := s.newEventBuilder(ctx, req.Key, req.Value, latest.Value)
	if err != nil {
		rsp.Status, err = errToStatus(err)
		return rsp, err
	}

	obj := builder.Meta
	obj.SetUpdatedBy(builder.Requester.GetUID().String())
	obj.SetUpdatedTimestampMillis(time.Now().UnixMilli())

	event, err := builder.toEvent()
	if err != nil {
		rsp.Status, err = errToStatus(err)
		return rsp, err
	}

	event.Event = WatchEvent_MODIFIED
	event.PreviousRV = latest.ResourceVersion

	rsp.ResourceVersion, err = s.store.WriteEvent(ctx, event)
	rsp.Status, err = errToStatus(err)
	if err == nil {
		rsp.Value = event.Value // with mutated fields
	} else {
		rsp.Status, err = errToStatus(err)
	}
	return rsp, err
}

func (s *server) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	ctx, span := s.tracer.Start(ctx, "storage_server.Delete")
	defer span.End()

	if err := s.Init(); err != nil {
		return nil, err
	}

	rsp := &DeleteResponse{}
	if req.ResourceVersion < 0 {
		return nil, apierrors.NewBadRequest("update must include the previous version")
	}

	latest, err := s.store.Read(ctx, &ReadRequest{
		Key: req.Key,
	})
	if err != nil {
		return nil, err
	}
	if latest.ResourceVersion != req.ResourceVersion {
		return nil, ErrOptimisticLockingFailed
	}

	now := metav1.NewTime(time.UnixMilli(s.now()))
	event := WriteEvent{
		EventID:    s.nextEventID(),
		Key:        req.Key,
		Event:      WatchEvent_DELETED,
		PreviousRV: latest.ResourceVersion,
	}
	requester, err := identity.GetRequester(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest("unable to get user")
	}
	marker := &DeletedMarker{}
	err = json.Unmarshal(latest.Value, marker)
	if err != nil {
		return nil, apierrors.NewBadRequest(
			fmt.Sprintf("unable to read previous object, %v", err))
	}
	obj, err := utils.MetaAccessor(marker)
	if err != nil {
		return nil, err
	}
	obj.SetDeletionTimestamp(&now)
	obj.SetUpdatedTimestamp(&now.Time)
	obj.SetManagedFields(nil)
	obj.SetFinalizers(nil)
	obj.SetUpdatedBy(requester.GetUID().String())
	marker.TypeMeta = metav1.TypeMeta{
		Kind:       "DeletedMarker",
		APIVersion: "storage.grafana.app/v0alpha1", // ?? or can we stick this in common?
	}
	marker.Annotations["RestoreResourceVersion"] = fmt.Sprintf("%d", event.PreviousRV)
	event.Value, err = json.Marshal(marker)
	if err != nil {
		return nil, apierrors.NewBadRequest(
			fmt.Sprintf("unable creating deletion marker, %v", err))
	}

	rsp.ResourceVersion, err = s.store.WriteEvent(ctx, event)
	rsp.Status, err = errToStatus(err)
	return rsp, err
}

func (s *server) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}

	if req.Key.Group == "" {
		status, _ := errToStatus(apierrors.NewBadRequest("missing group"))
		return &ReadResponse{Status: status}, nil
	}
	if req.Key.Resource == "" {
		status, _ := errToStatus(apierrors.NewBadRequest("missing resource"))
		return &ReadResponse{Status: status}, nil
	}

	rsp, err := s.store.Read(ctx, req)
	if err != nil {
		if rsp == nil {
			rsp = &ReadResponse{}
		}
		rsp.Status, err = errToStatus(err)
	}
	return rsp, err
}

func (s *server) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}

	rsp, err := s.store.List(ctx, req)
	// Status???
	return rsp, err
}

func (s *server) Watch(req *WatchRequest, srv ResourceStore_WatchServer) error {
	if err := s.Init(); err != nil {
		return err
	}

	// TODO??? can we move any of the common processing here?
	stream, err := s.store.Watch(srv.Context(), req)
	if err != nil {
		return err
	}
	for event := range stream {
		srv.Send(event)
	}
	return nil
}

// GetBlob implements ResourceServer.
func (s *server) GetBlob(ctx context.Context, req *GetBlobRequest) (*GetBlobResponse, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}
	rsp, err := s.search.GetBlob(ctx, req)
	rsp.Status, err = errToStatus(err)
	return rsp, err
}

// History implements ResourceServer.
func (s *server) History(ctx context.Context, req *HistoryRequest) (*HistoryResponse, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s.search.History(ctx, req)
}

// IsHealthy implements ResourceServer.
func (s *server) IsHealthy(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s.diagnostics.IsHealthy(ctx, req)
}

// Origin implements ResourceServer.
func (s *server) Origin(ctx context.Context, req *OriginRequest) (*OriginResponse, error) {
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s.search.Origin(ctx, req)
}
