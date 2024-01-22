package folder

import (
	"context"
)

type Service interface {
	// GetChildren returns an array containing all child folders.
	GetChildren(ctx context.Context, q *GetChildrenQuery) ([]*Folder, error)
	// GetParents returns an array containing add parent folders if nested folders are enabled
	// otherwise it returns an empty array
	GetParents(ctx context.Context, q GetParentsQuery) ([]*Folder, error)
	Create(ctx context.Context, cmd *CreateFolderCommand) (*Folder, error)

	// GetFolder takes a GetFolderCommand and returns a folder matching the
	// request. One of UID, ID or Title must be included. If multiple values
	// are included in the request, Grafana will select one in order of
	// specificity (UID, ID, Title).
	Get(ctx context.Context, q *GetFolderQuery) (*Folder, error)

	// Update is used to update a folder's UID, Title and Description. To change
	// a folder's parent folder, use Move.
	Update(ctx context.Context, cmd *UpdateFolderCommand) (*Folder, error)
	Delete(ctx context.Context, cmd *DeleteFolderCommand) error
	// Move changes a folder's parent folder to the requested new parent.
	Move(ctx context.Context, cmd *MoveFolderCommand) (*Folder, error)
	RegisterService(service RegistryService) error
	GetDescendantCounts(ctx context.Context, q *GetDescendantCountsQuery) (DescendantCounts, error)
	// TODO: double check if needed
	// GetFolders returns all folders for the given orgID and UIDs
	// If IncludeFullPaths is true, the full path of each folder will be included in the response.
	GetFolders(ctx context.Context, q *GetFoldersQuery) ([]*Folder, error)
	// WithFullpath sets the FullPath property of the folder to the full path of the folder.
	WithFullpath(ctx context.Context, folders *Folder, includeFullpath bool) (*Folder, error)
}

// FolderStore is a folder store.
//
//go:generate mockery --name FolderStore --structname FakeFolderStore --outpkg foldertest --output foldertest --filename folder_store_mock.go
type FolderStore interface {
	// GetFolderByTitle retrieves a folder by its title
	GetFolderByTitle(ctx context.Context, orgID int64, title string) (*Folder, error)
	// GetFolderByUID retrieves a folder by its UID
	GetFolderByUID(ctx context.Context, orgID int64, uid string) (*Folder, error)
	// GetFolderByID retrieves a folder by its ID
	GetFolderByID(ctx context.Context, orgID int64, id int64) (*Folder, error)
	// GetFolders returns all folders for the given orgID and UIDs.
	GetFolders(ctx context.Context, orgID int64, uids []string) (map[string]*Folder, error)
}
