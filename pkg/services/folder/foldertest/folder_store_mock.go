// Code generated by mockery v2.32.0. DO NOT EDIT.

package foldertest

import (
	context "context"

	folder "github.com/grafana/grafana/pkg/services/folder"
	mock "github.com/stretchr/testify/mock"
)

// FakeFolderStore is an autogenerated mock type for the FolderStore type
type FakeFolderStore struct {
	mock.Mock
}

// GetFolderByID provides a mock function with given fields: ctx, orgID, id
func (_m *FakeFolderStore) GetFolderByID(ctx context.Context, orgID int64, id int64) (*folder.Folder, error) {
	ret := _m.Called(ctx, orgID, id)

	var r0 *folder.Folder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (*folder.Folder, error)); ok {
		return rf(ctx, orgID, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) *folder.Folder); ok {
		r0 = rf(ctx, orgID, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*folder.Folder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, orgID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFolderByTitle provides a mock function with given fields: ctx, orgID, title, folderUID
func (_m *FakeFolderStore) GetFolderByTitle(ctx context.Context, orgID int64, title string, folderUID *string) (*folder.Folder, error) {
	ret := _m.Called(ctx, orgID, title, folderUID)

	var r0 *folder.Folder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *string) (*folder.Folder, error)); ok {
		return rf(ctx, orgID, title, folderUID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, *string) *folder.Folder); ok {
		r0 = rf(ctx, orgID, title, folderUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*folder.Folder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string, *string) error); ok {
		r1 = rf(ctx, orgID, title, folderUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFolderByUID provides a mock function with given fields: ctx, orgID, uid
func (_m *FakeFolderStore) GetFolderByUID(ctx context.Context, orgID int64, uid string) (*folder.Folder, error) {
	ret := _m.Called(ctx, orgID, uid)

	var r0 *folder.Folder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) (*folder.Folder, error)); ok {
		return rf(ctx, orgID, uid)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) *folder.Folder); ok {
		r0 = rf(ctx, orgID, uid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*folder.Folder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, orgID, uid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFolders provides a mock function with given fields: ctx, orgID, uids
func (_m *FakeFolderStore) GetFolders(ctx context.Context, orgID int64, uids []string) (map[string]*folder.Folder, error) {
	ret := _m.Called(ctx, orgID, uids)

	var r0 map[string]*folder.Folder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, []string) (map[string]*folder.Folder, error)); ok {
		return rf(ctx, orgID, uids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, []string) map[string]*folder.Folder); ok {
		r0 = rf(ctx, orgID, uids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*folder.Folder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, []string) error); ok {
		r1 = rf(ctx, orgID, uids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFakeFolderStore creates a new instance of FakeFolderStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFakeFolderStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *FakeFolderStore {
	mock := &FakeFolderStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
