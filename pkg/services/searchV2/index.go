package searchV2

import (
	"context"

	"github.com/grafana/grafana/pkg/services/store"

	"github.com/blugelabs/bluge"
)

// IndexFactory called by index manager when it requires an index for organization.
// If initialized with non-nil Writer then manager was able to load the index from
// a backup - so the factory can skip initial index building.
type IndexFactory func(ctx context.Context, orgID int64, writer *bluge.Writer) (Index, error)

// Index is responsible for a search inside one organization.
type Index interface {
	Reader() (*bluge.Reader, func(), error)
	ReIndex(ctx context.Context, force bool) error
	ApplyEvent(ctx context.Context, event store.ResourceEvent) error
	BackupTo(ctx context.Context, directory string) error
}
