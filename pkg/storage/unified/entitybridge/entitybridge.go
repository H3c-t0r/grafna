package entitybridge

import (
	"context"
	"fmt"
	"os"

	"gocloud.dev/blob/fileblob"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/store/entity"
	"github.com/grafana/grafana/pkg/services/store/entity/db/dbimpl"
	"github.com/grafana/grafana/pkg/services/store/entity/sqlstash"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/storage/unified/resource"
)

// Creates a ResourceServer using the existing entity tables
// NOTE: most of the field values are ignored
func ProvideResourceServer(db db.DB, cfg *setting.Cfg, features featuremgmt.FeatureToggles, tracer tracing.Tracer) (resource.ResourceServer, error) {
	if true {
		tmp, err := os.MkdirTemp("", "xxx-*")
		if err != nil {
			return nil, err
		}

		bucket, err := fileblob.OpenBucket(tmp, &fileblob.Options{
			CreateDir: true,
			Metadata:  fileblob.MetadataDontWrite, // skip
		})
		if err != nil {
			return nil, err
		}

		fmt.Printf("ROOT: %s\n\n", tmp)
		store, err := resource.NewCDKAppendingStore(context.Background(), resource.CDKAppenderOptions{
			Bucket: bucket,
		})
		if err != nil {
			return nil, err
		}

		return resource.NewResourceServer(resource.ResourceServerOptions{
			Store: store,
		})
	}

	eDB, err := dbimpl.ProvideEntityDB(db, cfg, features, tracer)
	if err != nil {
		return nil, err
	}

	entity, err := sqlstash.ProvideSQLEntityServer(eDB, tracer)
	if err != nil {
		return nil, err
	}

	store := &entityBridge{
		entity: entity,
	}
	return resource.NewResourceServer(resource.ResourceServerOptions{
		Tracer:      tracer,
		Store:       store,
		NodeID:      234, // from config?  used for snowflake ID
		Diagnostics: store,
		Lifecycle:   store,
	})
}

type entityBridge struct {
	entity sqlstash.SqlEntityServer
}

// Init implements ResourceServer.
func (b *entityBridge) Init() error {
	return b.entity.Init()
}

// Stop implements ResourceServer.
func (b *entityBridge) Stop() {
	b.entity.Stop()
}

// Convert resource key to the entity key
func toEntityKey(key *resource.ResourceKey) string {
	e := entity.Key{
		Group:     key.Group,
		Resource:  key.Resource,
		Namespace: key.Namespace,
		Name:      key.Name,
	}
	return e.String()
}

func (b *entityBridge) WriteEvent(ctx context.Context, event resource.WriteEvent) (int64, error) {
	key := toEntityKey(event.Key)

	// Delete does not need to create an entity first
	if event.Type == resource.WatchEvent_DELETED {
		rsp, err := b.entity.Delete(ctx, &entity.DeleteEntityRequest{
			Key:             key,
			PreviousVersion: event.PreviousRV,
		})
		if err != nil {
			return 0, err
		}
		return rsp.Entity.ResourceVersion, err
	}

	obj := event.Object
	msg := &entity.Entity{
		Key:       key,
		Group:     event.Key.Group,
		Resource:  event.Key.Resource,
		Namespace: event.Key.Namespace,
		Name:      event.Key.Name,
		Guid:      string(event.Object.GetUID()),

		//	Key:     fmt.Sprint("%s/%s/%s/%s", ),
		Folder: obj.GetFolder(),
		Body:   event.Value,

		Labels: obj.GetLabels(),
		Size:   int64(len(event.Value)),
	}

	switch event.Type {
	case resource.WatchEvent_ADDED:
		msg.Action = entity.Entity_CREATED
		rsp, err := b.entity.Create(ctx, &entity.CreateEntityRequest{Entity: msg})
		if err != nil {
			return 0, err
		}
		return rsp.Entity.ResourceVersion, err

	case resource.WatchEvent_MODIFIED:
		msg.Action = entity.Entity_UPDATED
		rsp, err := b.entity.Update(ctx, &entity.UpdateEntityRequest{
			Entity:          msg,
			PreviousVersion: event.PreviousRV,
		})
		if err != nil {
			return 0, err
		}
		return rsp.Entity.ResourceVersion, err

	default:
	}

	return 0, fmt.Errorf("unsupported operation: %s", event.Type.String())
}

func (b *entityBridge) WatchWriteEvents(ctx context.Context) (<-chan *resource.WrittenEvent, error) {
	return nil, resource.ErrNotImplementedYet
}

// IsHealthy implements ResourceServer.
func (b *entityBridge) IsHealthy(ctx context.Context, req *resource.HealthCheckRequest) (*resource.HealthCheckResponse, error) {
	rsp, err := b.entity.IsHealthy(ctx, &entity.HealthCheckRequest{
		Service: req.Service, // ??
	})
	if err != nil {
		return nil, err
	}
	return &resource.HealthCheckResponse{
		Status: resource.HealthCheckResponse_ServingStatus(rsp.Status),
	}, nil
}

// Read implements ResourceServer.
func (b *entityBridge) Read(ctx context.Context, req *resource.ReadRequest) (*resource.ReadResponse, error) {
	v, err := b.entity.Read(ctx, &entity.ReadEntityRequest{
		Key:      toEntityKey(req.Key),
		WithBody: true,
	})
	if err != nil {
		return nil, err
	}
	return &resource.ReadResponse{
		ResourceVersion: v.ResourceVersion,
		Value:           v.Body,
	}, nil
}

// List implements ResourceServer.
func (b *entityBridge) List(ctx context.Context, req *resource.ListRequest) (*resource.ListResponse, error) {
	key := req.Options.Key
	query := &entity.EntityListRequest{
		NextPageToken: req.NextPageToken,
		Limit:         req.Limit,
		Key:           []string{toEntityKey(key)},
		WithBody:      true,
	}

	// Assumes everything is equals
	if len(req.Options.Labels) > 0 {
		query.Labels = make(map[string]string)
		for _, q := range req.Options.Labels {
			query.Labels[q.Key] = q.Values[0]
		}
	}

	found, err := b.entity.List(ctx, query)
	if err != nil {
		return nil, err
	}

	rsp := &resource.ListResponse{
		ResourceVersion: found.ResourceVersion,
		NextPageToken:   found.NextPageToken,
	}
	for _, item := range found.Results {
		rsp.Items = append(rsp.Items, &resource.ResourceWrapper{
			ResourceVersion: item.ResourceVersion,
			Value:           item.Body,
		})
	}
	return rsp, nil
}

// GetBlob implements ResourceServer.
func (b *entityBridge) GetBlob(context.Context, *resource.GetBlobRequest) (*resource.GetBlobResponse, error) {
	return nil, resource.ErrNotImplementedYet
}

// History implements ResourceServer.
func (b *entityBridge) History(context.Context, *resource.HistoryRequest) (*resource.HistoryResponse, error) {
	return nil, resource.ErrNotImplementedYet
}

// Origin implements ResourceServer.
func (b *entityBridge) Origin(context.Context, *resource.OriginRequest) (*resource.OriginResponse, error) {
	return nil, resource.ErrNotImplementedYet
}
