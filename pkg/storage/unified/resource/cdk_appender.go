package resource

import (
	"bytes"
	context "context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/memblob"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewCDKAppendingStore(ctx context.Context, opts CDKOptions) (AppendingStore, error) {
	if opts.Tracer == nil {
		opts.Tracer = noop.NewTracerProvider().Tracer("cdk-appending-store")
	}

	if opts.Bucket == nil {
		return nil, fmt.Errorf("missing bucket")
	}

	found, _, err := opts.Bucket.ListPage(ctx, blob.FirstPageToken, 1, &blob.ListOptions{
		Prefix:    opts.RootFolder,
		Delimiter: "/",
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, fmt.Errorf("the root folder does not exist")
	}

	// This is not totally safe when running in HA
	if opts.NextResourceVersion == nil {
		if opts.NodeID == 0 {
			opts.NodeID = rand.Int63n(1024)
		}
		eventNode, err := snowflake.NewNode(opts.NodeID)
		if err != nil {
			return nil, apierrors.NewInternalError(
				fmt.Errorf("error initializing snowflake id generator :: %w", err))
		}
		opts.NextResourceVersion = func() int64 {
			return eventNode.Generate().Int64()
		}
	}

	return &cdkAppender{
		tracer: opts.Tracer,
		bucket: opts.Bucket,
		root:   opts.RootFolder,
		nextRV: opts.NextResourceVersion,
	}, nil
}

type cdkAppender struct {
	tracer trace.Tracer
	bucket *blob.Bucket
	root   string
	nextRV func() int64
	mutex  sync.Mutex
}

func (s *cdkAppender) getPath(key *ResourceKey, rv int64) string {
	var buffer bytes.Buffer
	buffer.WriteString(s.root)

	if key.Group == "" {
		return buffer.String()
	}
	buffer.WriteString(key.Group)

	if key.Resource == "" {
		return buffer.String()
	}
	buffer.WriteString("/")
	buffer.WriteString(key.Resource)

	if key.Namespace == "" {
		if key.Name == "" {
			return buffer.String()
		}
		buffer.WriteString("/__cluster__")
	} else {
		buffer.WriteString("/")
		buffer.WriteString(key.Namespace)
	}

	if key.Name == "" {
		return buffer.String()
	}
	buffer.WriteString("/")
	buffer.WriteString(key.Name)

	if rv > 0 {
		buffer.WriteString(fmt.Sprintf("/%d.json", rv))
	}
	return buffer.String()
}

func (s *cdkAppender) WriteEvent(ctx context.Context, event WriteEvent) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	body := fsEvent{
		ResourceVersion: s.nextRV(),
		Message:         event.Message,
		Operation:       event.Operation.String(),
		Value:           event.Value,
	}

	bytes, err := json.Marshal(&body)
	if err != nil {
		return 0, err
	}

	err = s.bucket.WriteAll(ctx, s.getPath(event.Key, body.ResourceVersion), bytes, &blob.WriterOptions{
		ContentType: "application/json",
	})
	return body.ResourceVersion, err
}

// Create new name for a given resource
func (s *cdkAppender) GenerateName(ctx context.Context, key *ResourceKey, prefix string) (string, error) {
	// TODO... shorter and make sure it does not exist
	return prefix + "x" + uuid.New().String(), nil
}

func (s *cdkAppender) open(ctx context.Context, p string) (*fsEvent, error) {
	raw, err := s.bucket.ReadAll(ctx, p)
	if err != nil {
		return nil, err
	}

	evt := &fsEvent{}
	err = json.Unmarshal(raw, evt)
	return evt, err
}

// Read implements ResourceStoreServer.
func (s *cdkAppender) Read(ctx context.Context, req *ReadRequest) (*ReadResponse, error) {
	rv := req.ResourceVersion

	path := s.getPath(req.Key, req.ResourceVersion)
	if rv < 1 {
		iter := s.bucket.List(&blob.ListOptions{Prefix: path + "/", Delimiter: "/"})
		for {
			obj, err := iter.Next(ctx)
			if err == io.EOF {
				break
			}
			if strings.HasSuffix(obj.Key, ".json") {
				idx := strings.LastIndex(obj.Key, "/") + 1
				edx := strings.LastIndex(obj.Key, ".")
				if idx > 0 {
					v, err := strconv.ParseInt(obj.Key[idx:edx], 10, 64)
					if err == nil && v > rv {
						path = obj.Key // find the path with biggest resource version
					}
				}
			}
		}
	}

	evt, err := s.open(ctx, path)
	if err != nil || evt.Operation == ResourceOperation_DELETED.String() {
		return nil, apierrors.NewNotFound(schema.GroupResource{
			Group:    req.Key.Group,
			Resource: req.Key.Resource,
		}, req.Key.Name)
	}

	return &ReadResponse{
		ResourceVersion: evt.ResourceVersion,
		Value:           evt.Value,
		Message:         evt.Message,
	}, nil
}

// List implements AppendingStore.
func (s *cdkAppender) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	resources, err := buildTree(ctx, s, req.Options.Key)
	if err != nil {
		return nil, err
	}

	rsp := &ListResponse{}
	for _, item := range resources {
		latest := item.versions[0]
		evt, err := s.open(ctx, latest.key)
		if err != nil {
			return nil, err
		}

		rsp.Items = append(rsp.Items, &ResourceWrapper{
			ResourceVersion: latest.rv,
			Value:           evt.Value,
			// Operation??
		})
	}
	return rsp, nil
}

// Watch implements AppendingStore.
func (s *cdkAppender) Watch(ctx context.Context, since int64, options *ListOptions) (chan *WatchEvent, error) {
	panic("unimplemented")
}

// group > resource > namespace > name > versions
type cdkResource struct {
	prefix   string
	versions []cdkVersion
}
type cdkVersion struct {
	rv  int64
	key string
}

func buildTree(ctx context.Context, s *cdkAppender, key *ResourceKey) ([]cdkResource, error) {
	byPrefix := make(map[string]*cdkResource)

	path := s.getPath(key, 0)
	iter := s.bucket.List(&blob.ListOptions{Prefix: path, Delimiter: ""}) // "" is recursive
	for {
		obj, err := iter.Next(ctx)
		if err == io.EOF {
			break
		}
		if strings.HasSuffix(obj.Key, ".json") {
			idx := strings.LastIndex(obj.Key, "/") + 1
			edx := strings.LastIndex(obj.Key, ".")
			if idx > 0 {
				rv, err := strconv.ParseInt(obj.Key[idx:edx], 10, 64)
				if err == nil {
					prefix := obj.Key[:idx]
					res, ok := byPrefix[prefix]
					if !ok {
						res = &cdkResource{prefix: prefix}
						byPrefix[prefix] = res
					}

					res.versions = append(res.versions, cdkVersion{
						rv:  rv,
						key: obj.Key,
					})
				}
			}
		}
	}

	// Now sort all versions
	resources := make([]cdkResource, 0, len(byPrefix))
	for _, res := range byPrefix {
		sort.Slice(res.versions, func(i, j int) bool {
			return res.versions[i].rv > res.versions[j].rv
		})
		resources = append(resources, *res)
	}
	sort.Slice(resources, func(i, j int) bool {
		a := resources[i].versions[0].rv
		b := resources[j].versions[0].rv
		return a > b
	})

	return resources, nil
}
