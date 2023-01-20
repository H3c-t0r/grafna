package folderimpl

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/accesscontrol/actest"
	acmock "github.com/grafana/grafana/pkg/services/accesscontrol/mock"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/folder"
	"github.com/grafana/grafana/pkg/services/guardian"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

var orgID = int64(1)
var usr = &user.SignedInUser{UserID: 1, OrgID: orgID}

func TestIntegrationProvideFolderService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("should register scope resolvers", func(t *testing.T) {
		cfg := setting.NewCfg()
		ac := acmock.New()
		ProvideService(ac, bus.ProvideBus(tracing.InitializeTracerForTest()), cfg, nil, nil, nil, &featuremgmt.FeatureManager{}, nil, nil)

		require.Len(t, ac.Calls.RegisterAttributeScopeResolver, 2)
	})
}

func TestIntegrationFolderService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Run("Folder service tests", func(t *testing.T) {
		dashStore := &dashboards.FakeDashboardStore{}
		db := sqlstore.InitTestDB(t)
		nestedFolderStore := ProvideStore(db, db.Cfg, featuremgmt.WithFeatures([]interface{}{"nestedFolders"}))

		folderStore := dashboards.NewFakeFolderStore(t)

		cfg := setting.NewCfg()
		cfg.RBACEnabled = false
		features := featuremgmt.WithFeatures()
		folderPermissions := acmock.NewMockedPermissionsService()

		service := &Service{
			cfg:                  cfg,
			log:                  log.New("test-folder-service"),
			dashboardStore:       dashStore,
			dashboardFolderStore: folderStore,
			store:                nestedFolderStore,
			searchService:        nil,
			features:             features,
			permissions:          folderPermissions,
			bus:                  bus.ProvideBus(tracing.InitializeTracerForTest()),
		}

		t.Run("Given user has no permissions", func(t *testing.T) {
			origNewGuardian := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{})

			folderId := rand.Int63()
			folderUID := util.GenerateShortUID()

			f := folder.NewFolder("Folder", "")
			f.ID = folderId
			f.UID = folderUID

			folderStore.On("GetFolderByID", mock.Anything, orgID, folderId).Return(f, nil)
			folderStore.On("GetFolderByUID", mock.Anything, orgID, folderUID).Return(f, nil)

			t.Run("When get folder by id should return access denied error", func(t *testing.T) {
				_, err := service.Get(context.Background(), &folder.GetFolderQuery{
					ID:           &folderId,
					OrgID:        orgID,
					SignedInUser: usr,
				})
				require.Equal(t, err, dashboards.ErrFolderAccessDenied)
			})

			var zeroInt int64 = 0
			t.Run("When get folder by id, with id = 0 should return default folder", func(t *testing.T) {
				foldr, err := service.Get(context.Background(), &folder.GetFolderQuery{
					ID:           &zeroInt,
					OrgID:        orgID,
					SignedInUser: usr,
				})
				require.NoError(t, err)
				require.Equal(t, foldr, &folder.Folder{ID: 0, Title: "General"})
			})

			t.Run("When get folder by uid should return access denied error", func(t *testing.T) {
				_, err := service.Get(context.Background(), &folder.GetFolderQuery{
					UID:          &folderUID,
					OrgID:        orgID,
					SignedInUser: usr,
				})
				require.Equal(t, err, dashboards.ErrFolderAccessDenied)
			})

			t.Run("When creating folder should return access denied error", func(t *testing.T) {
				dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil).Times(2)
				_, err := service.Create(context.Background(), &folder.CreateFolderCommand{
					OrgID:        orgID,
					Title:        f.Title,
					UID:          folderUID,
					SignedInUser: usr,
				})
				require.Equal(t, err, dashboards.ErrFolderAccessDenied)
			})

			title := "Folder-TEST"
			t.Run("When updating folder should return access denied error", func(t *testing.T) {
				dashStore.On("GetDashboard", mock.Anything, mock.AnythingOfType("*dashboards.GetDashboardQuery")).Run(func(args mock.Arguments) {
					folder := args.Get(1).(*dashboards.GetDashboardQuery)
					folder.Result = dashboards.NewDashboard("dashboard-test")
					folder.Result.IsFolder = true
				}).Return(&dashboards.Dashboard{}, nil)
				_, err := service.Update(context.Background(), &folder.UpdateFolderCommand{
					UID:          folderUID,
					OrgID:        orgID,
					NewTitle:     &title,
					SignedInUser: usr,
				})
				require.Equal(t, err, dashboards.ErrFolderAccessDenied)
			})

			t.Run("When deleting folder by uid should return access denied error", func(t *testing.T) {
				newFolder := folder.NewFolder("Folder", "")
				newFolder.UID = folderUID

				folderStore.On("GetFolderByID", mock.Anything, orgID, folderId).Return(newFolder, nil)
				folderStore.On("GetFolderByUID", mock.Anything, orgID, folderUID).Return(newFolder, nil)

				err := service.Delete(context.Background(), &folder.DeleteFolderCommand{
					UID:              folderUID,
					OrgID:            orgID,
					ForceDeleteRules: false,
					SignedInUser:     usr,
				})
				require.Error(t, err)
				require.Equal(t, err, dashboards.ErrFolderAccessDenied)
			})

			t.Cleanup(func() {
				guardian.New = origNewGuardian
			})
		})

		t.Run("Given user has permission to save", func(t *testing.T) {
			origNewGuardian := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})

			t.Run("When creating folder should not return access denied error", func(t *testing.T) {
				dash := dashboards.NewDashboardFolder("Test-Folder")
				dash.ID = rand.Int63()
				f := dashboards.FromDashboard(dash)

				dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
				dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(dash, nil).Once()
				folderStore.On("GetFolderByID", mock.Anything, orgID, dash.ID).Return(f, nil)

				actualFolder, err := service.Create(context.Background(), &folder.CreateFolderCommand{
					OrgID:        orgID,
					Title:        dash.Title,
					UID:          "someuid",
					SignedInUser: usr,
				})
				require.NoError(t, err)
				require.Equal(t, f, actualFolder)
			})

			t.Run("When creating folder should return error if uid is general", func(t *testing.T) {
				dash := dashboards.NewDashboardFolder("Test-Folder")
				dash.ID = rand.Int63()

				_, err := service.Create(context.Background(), &folder.CreateFolderCommand{
					OrgID:        orgID,
					Title:        dash.Title,
					UID:          "general",
					SignedInUser: usr,
				})
				require.ErrorIs(t, err, dashboards.ErrFolderInvalidUID)
			})

			t.Run("When updating folder should not return access denied error", func(t *testing.T) {
				dashboardFolder := dashboards.NewDashboardFolder("Folder")
				dashboardFolder.ID = rand.Int63()
				dashboardFolder.UID = util.GenerateShortUID()
				f := dashboards.FromDashboard(dashboardFolder)

				dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
				dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(dashboardFolder, nil)
				folderStore.On("GetFolderByID", mock.Anything, orgID, dashboardFolder.ID).Return(f, nil)

				title := "TEST-Folder"
				req := &folder.UpdateFolderCommand{
					UID:          dashboardFolder.UID,
					OrgID:        orgID,
					NewTitle:     &title,
					SignedInUser: usr,
				}

				reqResult, err := service.Update(context.Background(), req)
				require.NoError(t, err)
				require.Equal(t, f, reqResult)
			})

			t.Run("When deleting folder by uid should not return access denied error", func(t *testing.T) {
				f := folder.NewFolder(util.GenerateShortUID(), "")
				f.ID = rand.Int63()
				f.UID = util.GenerateShortUID()
				folderStore.On("GetFolderByUID", mock.Anything, orgID, f.UID).Return(f, nil)

				var actualCmd *dashboards.DeleteDashboardCommand
				dashStore.On("DeleteDashboard", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					actualCmd = args.Get(1).(*dashboards.DeleteDashboardCommand)
				}).Return(nil).Once()

				expectedForceDeleteRules := rand.Int63()%2 == 0
				err := service.Delete(context.Background(), &folder.DeleteFolderCommand{
					UID:              f.UID,
					OrgID:            orgID,
					ForceDeleteRules: expectedForceDeleteRules,
					SignedInUser:     usr,
				})
				require.NoError(t, err)
				require.NotNil(t, actualCmd)
				require.Equal(t, f.ID, actualCmd.ID)
				require.Equal(t, orgID, actualCmd.OrgID)
				require.Equal(t, expectedForceDeleteRules, actualCmd.ForceDeleteFolderRules)
			})

			t.Cleanup(func() {
				guardian.New = origNewGuardian
			})
		})

		t.Run("Given user has permission to view", func(t *testing.T) {
			origNewGuardian := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanViewValue: true})

			t.Run("When get folder by id should return folder", func(t *testing.T) {
				expected := folder.NewFolder(util.GenerateShortUID(), "")
				expected.ID = rand.Int63()

				folderStore.On("GetFolderByID", mock.Anything, orgID, expected.ID).Return(expected, nil)

				actual, err := service.getFolderByID(context.Background(), usr, expected.ID, orgID)
				require.Equal(t, expected, actual)
				require.NoError(t, err)
			})

			t.Run("When get folder by uid should return folder", func(t *testing.T) {
				expected := folder.NewFolder(util.GenerateShortUID(), "")
				expected.UID = util.GenerateShortUID()

				folderStore.On("GetFolderByUID", mock.Anything, orgID, expected.UID).Return(expected, nil)

				actual, err := service.getFolderByUID(context.Background(), usr, orgID, expected.UID)
				require.Equal(t, expected, actual)
				require.NoError(t, err)
			})

			t.Run("When get folder by title should return folder", func(t *testing.T) {
				expected := folder.NewFolder("TEST-"+util.GenerateShortUID(), "")

				folderStore.On("GetFolderByTitle", mock.Anything, orgID, expected.Title).Return(expected, nil)

				actual, err := service.getFolderByTitle(context.Background(), usr, orgID, expected.Title)
				require.Equal(t, expected, actual)
				require.NoError(t, err)
			})

			t.Cleanup(func() {
				guardian.New = origNewGuardian
			})
		})

		t.Run("Should map errors correct", func(t *testing.T) {
			testCases := []struct {
				ActualError   error
				ExpectedError error
			}{
				{ActualError: dashboards.ErrDashboardTitleEmpty, ExpectedError: dashboards.ErrFolderTitleEmpty},
				{ActualError: dashboards.ErrDashboardUpdateAccessDenied, ExpectedError: dashboards.ErrFolderAccessDenied},
				{ActualError: dashboards.ErrDashboardWithSameNameInFolderExists, ExpectedError: dashboards.ErrFolderSameNameExists},
				{ActualError: dashboards.ErrDashboardWithSameUIDExists, ExpectedError: dashboards.ErrFolderWithSameUIDExists},
				{ActualError: dashboards.ErrDashboardVersionMismatch, ExpectedError: dashboards.ErrFolderVersionMismatch},
				{ActualError: dashboards.ErrDashboardNotFound, ExpectedError: dashboards.ErrFolderNotFound},
				{ActualError: dashboards.ErrDashboardFailedGenerateUniqueUid, ExpectedError: dashboards.ErrFolderFailedGenerateUniqueUid},
				{ActualError: dashboards.ErrDashboardInvalidUid, ExpectedError: dashboards.ErrDashboardInvalidUid},
			}

			for _, tc := range testCases {
				actualError := toFolderError(tc.ActualError)
				assert.EqualErrorf(t, actualError, tc.ExpectedError.Error(),
					"For error '%s' expected error '%s', actual '%s'", tc.ActualError, tc.ExpectedError, actualError)
			}
		})
	})
}

func TestNestedFolderServiceFeatureToggle(t *testing.T) {
	g := guardian.New
	guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
	t.Cleanup(func() {
		guardian.New = g
	})

	nestedFolderStore := NewFakeStore()

	dashStore := dashboards.FakeDashboardStore{}
	dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
	dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(&dashboards.Dashboard{}, nil)

	dashboardFolderStore := dashboards.NewFakeFolderStore(t)
	dashboardFolderStore.On("GetFolderByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&folder.Folder{}, nil)

	cfg := setting.NewCfg()
	cfg.RBACEnabled = false
	folderService := &Service{
		cfg:                  cfg,
		store:                nestedFolderStore,
		dashboardStore:       &dashStore,
		dashboardFolderStore: dashboardFolderStore,
		features:             featuremgmt.WithFeatures(featuremgmt.FlagNestedFolders),
		log:                  log.New("test-folder-service"),
	}
	t.Run("create folder", func(t *testing.T) {
		nestedFolderStore.ExpectedFolder = &folder.Folder{ParentUID: util.GenerateShortUID()}
		res, err := folderService.Create(context.Background(), &folder.CreateFolderCommand{SignedInUser: usr, Title: "my folder"})
		require.NoError(t, err)
		require.NotNil(t, res.UID)
		require.NotEmpty(t, res.ParentUID)
	})
}

func TestNestedFolderService(t *testing.T) {
	t.Run("with feature flag unset", func(t *testing.T) {
		t.Run("When create folder, no create in folder table done", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			// dashboard store & service commands that should be called.
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
			dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(&dashboards.Dashboard{}, nil)

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&folder.Folder{}, nil)

			nestedFolderStore := NewFakeStore()

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures(), nil)
			_, err := folderSvc.Create(context.Background(), &folder.CreateFolderCommand{
				OrgID:        orgID,
				Title:        "myFolder",
				UID:          "myFolder",
				SignedInUser: usr,
			})
			require.NoError(t, err)
			// CreateFolder should not call the folder store create if the feature toggle is not enabled.
			require.False(t, nestedFolderStore.CreateCalled)
		})

		t.Run("When delete folder, no delete in folder table done", func(t *testing.T) {
			var actualCmd *dashboards.DeleteDashboardCommand
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("DeleteDashboard", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				actualCmd = args.Get(1).(*dashboards.DeleteDashboardCommand)
			}).Return(nil).Once()

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByUID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(&folder.Folder{}, nil)

			nestedFolderStore := NewFakeStore()

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures(), nil)
			err := folderSvc.Delete(context.Background(), &folder.DeleteFolderCommand{UID: "myFolder", OrgID: orgID, SignedInUser: usr})
			require.NoError(t, err)
			require.NotNil(t, actualCmd)

			require.False(t, nestedFolderStore.DeleteCalled)
		})
	})

	t.Run("with nested folder feature flag on", func(t *testing.T) {
		t.Run("create, no error", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			// dashboard store commands that should be called.
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
			dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(&dashboards.Dashboard{}, nil)

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&folder.Folder{}, nil)

			nestedFolderStore := NewFakeStore()

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			_, err := folderSvc.Create(context.Background(), &folder.CreateFolderCommand{
				OrgID:        orgID,
				Title:        "myFolder",
				UID:          "myFolder",
				SignedInUser: usr,
			})
			require.NoError(t, err)
			// CreateFolder should also call the folder store's create method.
			require.True(t, nestedFolderStore.CreateCalled)
		})

		t.Run("create without UID, no error", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			// dashboard store commands that should be called.
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
			dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(&dashboards.Dashboard{UID: "newUID"}, nil)

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&folder.Folder{}, nil)

			nestedFolderStore := NewFakeStore()

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			f, err := folderSvc.Create(context.Background(), &folder.CreateFolderCommand{
				OrgID:        orgID,
				Title:        "myFolder",
				SignedInUser: usr,
			})
			require.NoError(t, err)
			// CreateFolder should also call the folder store's create method.
			require.True(t, nestedFolderStore.CreateCalled)
			require.Equal(t, "newUID", f.UID)
		})

		t.Run("create failed because of circular reference", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashboardFolder := dashboards.NewDashboardFolder("myFolder")
			dashboardFolder.ID = rand.Int63()
			dashboardFolder.UID = "myFolder"
			f := dashboards.FromDashboard(dashboardFolder)

			// dashboard store commands that should be called.
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
			dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(dashboardFolder, nil)
			var actualCmd *dashboards.DeleteDashboardCommand
			dashStore.On("DeleteDashboard", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				actualCmd = args.Get(1).(*dashboards.DeleteDashboardCommand)
			}).Return(nil).Once()

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByID", mock.Anything, orgID, dashboardFolder.ID).Return(f, nil)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedParentFolders = []*folder.Folder{
				{UID: "newFolder", ParentUID: "newFolder"},
				{UID: "newFolder2", ParentUID: "newFolder2"},
				{UID: "newFolder3", ParentUID: "newFolder3"},
				{UID: "myFolder", ParentUID: "newFolder"},
			}

			cmd := folder.CreateFolderCommand{
				ParentUID:    "myFolder1",
				OrgID:        orgID,
				Title:        "myFolder",
				UID:          "myFolder",
				SignedInUser: usr,
			}

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			_, err := folderSvc.Create(context.Background(), &cmd)
			require.Error(t, err, folder.ErrCircularReference)
			// CreateFolder should not call the folder store's create method.
			require.False(t, nestedFolderStore.CreateCalled)
			require.NotNil(t, actualCmd)
		})

		t.Run("create returns error from nested folder service", func(t *testing.T) {
			// This test creates and deletes the dashboard, so needs some extra setup.
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			// dashboard store commands that should be called.
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil)
			dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(&dashboards.Dashboard{}, nil)
			var actualCmd *dashboards.DeleteDashboardCommand
			dashStore.On("DeleteDashboard", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				actualCmd = args.Get(1).(*dashboards.DeleteDashboardCommand)
			}).Return(nil).Once()

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&folder.Folder{}, nil)

			// return an error from the folder store
			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedError = errors.New("FAILED")

			// the service return success as long as the legacy create succeeds
			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			_, err := folderSvc.Create(context.Background(), &folder.CreateFolderCommand{
				OrgID:        orgID,
				Title:        "myFolder",
				UID:          "myFolder",
				SignedInUser: usr,
			})
			require.Error(t, err, "FAILED")

			// CreateFolder should also call the folder store's create method.
			require.True(t, nestedFolderStore.CreateCalled)
			require.NotNil(t, actualCmd)
		})

		t.Run("move, no view permission should fail", func(t *testing.T) {
			// This test creates and deletes the dashboard, so needs some extra setup.
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanViewValue: false})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			dashboardFolderStore := dashboards.NewFakeFolderStore(t)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			_, err := folderSvc.Move(context.Background(), &folder.MoveFolderCommand{UID: "myFolder", NewParentUID: "newFolder", OrgID: orgID, SignedInUser: usr})
			require.Error(t, err, dashboards.ErrFolderAccessDenied)
		})

		t.Run("move, no save permission should fail", func(t *testing.T) {
			// This test creates and deletes the dashboard, so needs some extra setup.
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: false, CanViewValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			dashboardFolderStore := dashboards.NewFakeFolderStore(t)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			_, err := folderSvc.Move(context.Background(), &folder.MoveFolderCommand{UID: "myFolder", NewParentUID: "newFolder", OrgID: orgID, SignedInUser: usr})
			require.Error(t, err, dashboards.ErrFolderAccessDenied)
		})

		t.Run("move, no error", func(t *testing.T) {
			// This test creates and deletes the dashboard, so needs some extra setup.
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true, CanViewValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			dashboardFolderStore := dashboards.NewFakeFolderStore(t)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}
			nestedFolderStore.ExpectedParentFolders = []*folder.Folder{
				{UID: "newFolder", ParentUID: "newFolder"},
				{UID: "newFolder2", ParentUID: "newFolder2"},
				{UID: "newFolder3", ParentUID: "newFolder3"},
			}

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			f, err := folderSvc.Move(context.Background(), &folder.MoveFolderCommand{UID: "myFolder", NewParentUID: "newFolder", OrgID: orgID, SignedInUser: usr})
			require.NoError(t, err)
			require.NotNil(t, f)
		})

		t.Run("move when parentUID in the current subtree returns error from nested folder service", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true, CanViewValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			dashboardFolderStore := dashboards.NewFakeFolderStore(t)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}
			nestedFolderStore.ExpectedError = folder.ErrCircularReference

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			f, err := folderSvc.Move(context.Background(), &folder.MoveFolderCommand{UID: "myFolder", NewParentUID: "newFolder", OrgID: orgID, SignedInUser: usr})
			require.Error(t, err, folder.ErrCircularReference)
			require.Nil(t, f)
		})

		t.Run("move when new parentUID depth + subTree height bypassed maximum depth returns error", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true, CanViewValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			dashboardFolderStore := dashboards.NewFakeFolderStore(t)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}
			nestedFolderStore.ExpectedParentFolders = []*folder.Folder{
				{UID: "newFolder", ParentUID: "newFolder"},
				{UID: "newFolder2", ParentUID: "newFolder2"},
			}
			nestedFolderStore.ExpectedFolderHeight = 5

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			f, err := folderSvc.Move(context.Background(), &folder.MoveFolderCommand{UID: "myFolder", NewParentUID: "newFolder2", OrgID: orgID, SignedInUser: usr})
			require.Error(t, err, folder.ErrMaximumDepthReached)
			require.Nil(t, f)
		})

		t.Run("move when parentUID in the current subtree returns error from nested folder service", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true, CanViewValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			dashboardFolderStore := dashboards.NewFakeFolderStore(t)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}
			nestedFolderStore.ExpectedParentFolders = []*folder.Folder{{UID: "myFolder", ParentUID: "12345"}, {UID: "12345", ParentUID: ""}}

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			f, err := folderSvc.Move(context.Background(), &folder.MoveFolderCommand{UID: "myFolder", NewParentUID: "newFolder2", OrgID: orgID, SignedInUser: usr})
			require.Error(t, err, folder.ErrCircularReference)
			require.Nil(t, f)
		})

		t.Run("delete with success", func(t *testing.T) {
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true, CanViewValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			dashStore := &dashboards.FakeDashboardStore{}
			var actualCmd *dashboards.DeleteDashboardCommand
			dashStore.On("DeleteDashboard", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				actualCmd = args.Get(1).(*dashboards.DeleteDashboardCommand)
			}).Return(nil).Once()

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByUID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("string")).Return(&folder.Folder{}, nil)

			nestedFolderStore := NewFakeStore()
			nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			err := folderSvc.Delete(context.Background(), &folder.DeleteFolderCommand{UID: "myFolder", OrgID: orgID, SignedInUser: usr})
			require.NoError(t, err)
			require.NotNil(t, actualCmd)

			require.True(t, nestedFolderStore.DeleteCalled)
		})

		t.Run("create returns error if maximum depth reached", func(t *testing.T) {
			// This test creates and deletes the dashboard, so needs some extra setup.
			g := guardian.New
			guardian.MockDashboardGuardian(&guardian.FakeDashboardGuardian{CanSaveValue: true})
			t.Cleanup(func() {
				guardian.New = g
			})

			// dashboard store commands that should be called.
			dashStore := &dashboards.FakeDashboardStore{}
			dashStore.On("ValidateDashboardBeforeSave", mock.Anything, mock.AnythingOfType("*dashboards.Dashboard"), mock.AnythingOfType("bool")).Return(true, nil).Times(2)
			dashStore.On("SaveDashboard", mock.Anything, mock.AnythingOfType("dashboards.SaveDashboardCommand")).Return(&dashboards.Dashboard{}, nil)
			var actualCmd *dashboards.DeleteDashboardCommand
			dashStore.On("DeleteDashboard", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				actualCmd = args.Get(1).(*dashboards.DeleteDashboardCommand)
			}).Return(nil).Once()

			dashboardFolderStore := dashboards.NewFakeFolderStore(t)
			dashboardFolderStore.On("GetFolderByID", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(&folder.Folder{}, nil)

			parents := make([]*folder.Folder, 0, folder.MaxNestedFolderDepth)
			for i := 0; i < folder.MaxNestedFolderDepth; i++ {
				parents = append(parents, &folder.Folder{UID: fmt.Sprintf("folder%d", i)})
			}

			nestedFolderStore := NewFakeStore()
			//nestedFolderStore.ExpectedFolder = &folder.Folder{UID: "myFolder", ParentUID: "newFolder"}
			nestedFolderStore.ExpectedParentFolders = parents

			folderSvc := setup(t, dashStore, dashboardFolderStore, nestedFolderStore, featuremgmt.WithFeatures("nestedFolders"), actest.FakeAccessControl{
				ExpectedEvaluate: true,
			})
			_, err := folderSvc.Create(context.Background(), &folder.CreateFolderCommand{
				Title:        "folder",
				OrgID:        orgID,
				ParentUID:    parents[len(parents)-1].UID,
				UID:          util.GenerateShortUID(),
				SignedInUser: usr,
			})
			assert.ErrorIs(t, err, folder.ErrMaximumDepthReached)
			require.NotNil(t, actualCmd)
		})
	})
}

func setup(t *testing.T, dashStore dashboards.Store, dashboardFolderStore dashboards.FolderStore, nestedFolderStore store, features featuremgmt.FeatureToggles, ac accesscontrol.AccessControl) folder.Service {
	t.Helper()

	// nothing enabled yet
	cfg := setting.NewCfg()
	cfg.RBACEnabled = false
	return &Service{
		cfg:                  cfg,
		log:                  log.New("test-folder-service"),
		dashboardStore:       dashStore,
		dashboardFolderStore: dashboardFolderStore,
		store:                nestedFolderStore,
		features:             features,
		accessControl:        ac,
	}
}
