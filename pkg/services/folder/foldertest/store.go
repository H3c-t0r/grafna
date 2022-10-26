package foldertest

import (
	"context"

	"github.com/grafana/grafana/pkg/services/folder"
)

type FakeStore struct {
	ExpectedFolders []*folder.Folder
	ExpectedFolder  *folder.Folder
	ExpectedError   error
}

var _ folder.Store = (*FakeStore)(nil)

func (f *FakeStore) Create(ctx context.Context, cmd *folder.CreateFolderCommand) (*folder.Folder, error) {
	return f.ExpectedFolder, f.ExpectedError
}

func (f *FakeStore) Delete(ctx context.Context, uid string, orgID int64) error {
	return f.ExpectedError
}

func (f *FakeStore) Update(ctx context.Context, cmd *folder.UpdateFolderCommand) (*folder.Folder, error) {
	return f.ExpectedFolder, f.ExpectedError
}

func (f *FakeStore) Move(ctx context.Context, cmd *folder.MoveFolderCommand) (*folder.Folder, error) {
	return f.ExpectedFolder, f.ExpectedError
}

func (f *FakeStore) Get(ctx context.Context, cmd *folder.GetFolderCommand) (*folder.Folder, error) {
	return f.ExpectedFolder, f.ExpectedError
}

func (f *FakeStore) GetParents(ctx context.Context, cmd *folder.GetParentsCommand) ([]*folder.Folder, error) {
	return f.ExpectedFolders, f.ExpectedError
}

func (f *FakeStore) GetChildren(ctx context.Context, cmd *folder.GetTreeCommand) ([]*folder.Folder, error) {
	return f.ExpectedFolders, f.ExpectedError
}
