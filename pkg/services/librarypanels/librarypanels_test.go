package librarypanels

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	dboards "github.com/grafana/grafana/pkg/dashboards"

	"github.com/grafana/grafana/pkg/services/libraryelements"

	"github.com/stretchr/testify/require"
	"gopkg.in/macaron.v1"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
)

const UserInDbName = "user_in_db"
const UserInDbAvatar = "/avatar/402d08de060496d6b6874495fe20f5ad"

func TestLoadLibraryPanelsForDashboard(t *testing.T) {
	scenarioWithLibraryPanel(t, "When an admin tries to load a dashboard with a library panel, it should copy JSON properties from library panel",
		func(t *testing.T, sc scenarioContext) {
			dashJSON := map[string]interface{}{
				"panels": []interface{}{
					map[string]interface{}{
						"id": int64(1),
						"gridPos": map[string]interface{}{
							"h": 6,
							"w": 6,
							"x": 0,
							"y": 0,
						},
					},
					map[string]interface{}{
						"id": int64(2),
						"gridPos": map[string]interface{}{
							"h": 6,
							"w": 6,
							"x": 6,
							"y": 0,
						},
						"libraryPanel": map[string]interface{}{
							"uid":  sc.initialResult.Result.UID,
							"name": sc.initialResult.Result.Name,
						},
					},
				},
			}
			dash := models.Dashboard{
				Title: "Testing LoadLibraryPanelsForDashboard",
				Data:  simplejson.NewFromAny(dashJSON),
			}
			dashInDB := createDashboard(t, sc.sqlStore, sc.user, &dash, sc.folder.Id)
			err := sc.elementService.ConnectElementsToDashboard(sc.reqContext, []string{sc.initialResult.Result.UID}, dashInDB.Id)
			require.NoError(t, err)

			err = sc.service.LoadLibraryPanelsForDashboard(sc.reqContext, dashInDB)
			require.NoError(t, err)
			expectedJSON := map[string]interface{}{
				"title":   "Testing LoadLibraryPanelsForDashboard",
				"uid":     dashInDB.Uid,
				"version": dashInDB.Version,
				"panels": []interface{}{
					map[string]interface{}{
						"id": int64(1),
						"gridPos": map[string]interface{}{
							"h": 6,
							"w": 6,
							"x": 0,
							"y": 0,
						},
					},
					map[string]interface{}{
						"id": int64(2),
						"gridPos": map[string]interface{}{
							"h": 6,
							"w": 6,
							"x": 6,
							"y": 0,
						},
						"datasource":  "${DS_GDEV-TESTDATA}",
						"description": "A description",
						"libraryPanel": map[string]interface{}{
							"uid":         sc.initialResult.Result.UID,
							"name":        sc.initialResult.Result.Name,
							"type":        sc.initialResult.Result.Type,
							"description": sc.initialResult.Result.Description,
							"version":     sc.initialResult.Result.Version,
							"meta": map[string]interface{}{
								"folderName":          "ScenarioFolder",
								"folderUid":           sc.folder.Uid,
								"connectedDashboards": int64(1),
								"created":             sc.initialResult.Result.Meta.Created,
								"updated":             sc.initialResult.Result.Meta.Updated,
								"createdBy": map[string]interface{}{
									"id":        sc.initialResult.Result.Meta.CreatedBy.ID,
									"name":      UserInDbName,
									"avatarUrl": UserInDbAvatar,
								},
								"updatedBy": map[string]interface{}{
									"id":        sc.initialResult.Result.Meta.UpdatedBy.ID,
									"name":      UserInDbName,
									"avatarUrl": UserInDbAvatar,
								},
							},
						},
						"title": "Text - Library Panel",
						"type":  "text",
					},
				},
			}
			expected := simplejson.NewFromAny(expectedJSON)
			if diff := cmp.Diff(expected.Interface(), dash.Data.Interface(), getCompareOptions()...); diff != "" {
				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
			}
		})

	//	scenarioWithLibraryPanel(t, "When an admin tries to load a dashboard with a library panel without uid, it should fail",
	//		func(t *testing.T, sc scenarioContext) {
	//			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID, ":dashboardId": "1"})
	//			resp := sc.service.connectHandler(sc.reqContext)
	//			require.Equal(t, 200, resp.Status())
	//
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"libraryPanel": map[string]interface{}{
	//							"name": sc.initialResult.Result.Name,
	//						},
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   1,
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.LoadLibraryPanelsForDashboard(sc.reqContext, &dash)
	//			require.EqualError(t, err, errLibraryPanelHeaderUIDMissing.Error())
	//		})
	//
	//	scenarioWithLibraryPanel(t, "When an admin tries to load a dashboard with a library panel that is not connected, it should set correct JSON and continue",
	//		func(t *testing.T, sc scenarioContext) {
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"libraryPanel": map[string]interface{}{
	//							"uid":  sc.initialResult.Result.UID,
	//							"name": sc.initialResult.Result.Name,
	//						},
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   1,
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.LoadLibraryPanelsForDashboard(sc.reqContext, &dash)
	//			require.NoError(t, err)
	//			expectedJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"libraryPanel": map[string]interface{}{
	//							"uid":  sc.initialResult.Result.UID,
	//							"name": sc.initialResult.Result.Name,
	//						},
	//						"type": fmt.Sprintf("Name: \"%s\", UID: \"%s\"", sc.initialResult.Result.Name, sc.initialResult.Result.UID),
	//					},
	//				},
	//			}
	//			expected := simplejson.NewFromAny(expectedJSON)
	//			if diff := cmp.Diff(expected.Interface(), dash.Data.Interface(), getCompareOptions()...); diff != "" {
	//				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
	//			}
	//		})
	//}
	//
	//func TestCleanLibraryPanelsForDashboard(t *testing.T) {
	//	scenarioWithLibraryPanel(t, "When an admin tries to store a dashboard with a library panel, it should just keep the correct JSON properties in library panel",
	//		func(t *testing.T, sc scenarioContext) {
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"datasource": "${DS_GDEV-TESTDATA}",
	//						"libraryPanel": map[string]interface{}{
	//							"uid":  sc.initialResult.Result.UID,
	//							"name": sc.initialResult.Result.Name,
	//						},
	//						"title": "Text - Library Panel",
	//						"type":  "text",
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   1,
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.CleanLibraryPanelsForDashboard(&dash)
	//			require.NoError(t, err)
	//			expectedJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"libraryPanel": map[string]interface{}{
	//							"uid":  sc.initialResult.Result.UID,
	//							"name": sc.initialResult.Result.Name,
	//						},
	//					},
	//				},
	//			}
	//			expected := simplejson.NewFromAny(expectedJSON)
	//			if diff := cmp.Diff(expected.Interface(), dash.Data.Interface(), getCompareOptions()...); diff != "" {
	//				t.Fatalf("Result mismatch (-want +got):\n%s", diff)
	//			}
	//		})
	//
	//	scenarioWithLibraryPanel(t, "When an admin tries to store a dashboard with a library panel without uid, it should fail",
	//		func(t *testing.T, sc scenarioContext) {
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"datasource": "${DS_GDEV-TESTDATA}",
	//						"libraryPanel": map[string]interface{}{
	//							"name": sc.initialResult.Result.Name,
	//						},
	//						"title": "Text - Library Panel",
	//						"type":  "text",
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   1,
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.CleanLibraryPanelsForDashboard(&dash)
	//			require.EqualError(t, err, errLibraryPanelHeaderUIDMissing.Error())
	//		})
	//
	//	scenarioWithLibraryPanel(t, "When an admin tries to store a dashboard with a library panel without name, it should fail",
	//		func(t *testing.T, sc scenarioContext) {
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"datasource": "${DS_GDEV-TESTDATA}",
	//						"libraryPanel": map[string]interface{}{
	//							"uid": sc.initialResult.Result.UID,
	//						},
	//						"title": "Text - Library Panel",
	//						"type":  "text",
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   1,
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.CleanLibraryPanelsForDashboard(&dash)
	//			require.EqualError(t, err, errLibraryPanelHeaderNameMissing.Error())
	//		})
	//}
	//
	//func TestConnectLibraryPanelsForDashboard(t *testing.T) {
	//	scenarioWithLibraryPanel(t, "When an admin tries to store a dashboard with a library panel, it should connect the two",
	//		func(t *testing.T, sc scenarioContext) {
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"datasource": "${DS_GDEV-TESTDATA}",
	//						"libraryPanel": map[string]interface{}{
	//							"uid":  sc.initialResult.Result.UID,
	//							"name": sc.initialResult.Result.Name,
	//						},
	//						"title": "Text - Library Panel",
	//						"type":  "text",
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   int64(1),
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.ConnectLibraryPanelsForDashboard(sc.reqContext, &dash)
	//			require.NoError(t, err)
	//
	//			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
	//			resp := sc.service.getConnectedDashboardsHandler(sc.reqContext)
	//			require.Equal(t, 200, resp.Status())
	//
	//			var dashResult libraryPanelDashboardsResult
	//			err = json.Unmarshal(resp.Body(), &dashResult)
	//			require.NoError(t, err)
	//			require.Len(t, dashResult.Result, 1)
	//			require.Equal(t, int64(1), dashResult.Result[0])
	//		})
	//
	//	scenarioWithLibraryPanel(t, "When an admin tries to store a dashboard with a library panel without uid, it should fail",
	//		func(t *testing.T, sc scenarioContext) {
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"datasource": "${DS_GDEV-TESTDATA}",
	//						"libraryPanel": map[string]interface{}{
	//							"name": sc.initialResult.Result.Name,
	//						},
	//						"title": "Text - Library Panel",
	//						"type":  "text",
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   int64(1),
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.ConnectLibraryPanelsForDashboard(sc.reqContext, &dash)
	//			require.EqualError(t, err, errLibraryPanelHeaderUIDMissing.Error())
	//		})
	//
	//	scenarioWithLibraryPanel(t, "When an admin tries to store a dashboard with unused/removed library panels, it should disconnect unused/removed library panels",
	//		func(t *testing.T, sc scenarioContext) {
	//			command := getCreateCommand(sc.folder.Id, "Unused Libray Panel")
	//			resp := sc.service.createHandler(sc.reqContext, command)
	//			var unused = validateAndUnMarshalResponse(t, resp)
	//			sc.reqContext.ReplaceAllParams(map[string]string{":uid": unused.Result.UID, ":dashboardId": "1"})
	//			resp = sc.service.connectHandler(sc.reqContext)
	//			require.Equal(t, 200, resp.Status())
	//
	//			dashJSON := map[string]interface{}{
	//				"panels": []interface{}{
	//					map[string]interface{}{
	//						"id": int64(1),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 0,
	//							"y": 0,
	//						},
	//					},
	//					map[string]interface{}{
	//						"id": int64(2),
	//						"gridPos": map[string]interface{}{
	//							"h": 6,
	//							"w": 6,
	//							"x": 6,
	//							"y": 0,
	//						},
	//						"datasource": "${DS_GDEV-TESTDATA}",
	//						"libraryPanel": map[string]interface{}{
	//							"uid":  sc.initialResult.Result.UID,
	//							"name": sc.initialResult.Result.Name,
	//						},
	//						"title": "Text - Library Panel",
	//						"type":  "text",
	//					},
	//				},
	//			}
	//			dash := models.Dashboard{
	//				Id:   int64(1),
	//				Data: simplejson.NewFromAny(dashJSON),
	//			}
	//
	//			err := sc.service.ConnectLibraryPanelsForDashboard(sc.reqContext, &dash)
	//			require.NoError(t, err)
	//
	//			sc.reqContext.ReplaceAllParams(map[string]string{":uid": sc.initialResult.Result.UID})
	//			resp = sc.service.getConnectedDashboardsHandler(sc.reqContext)
	//			require.Equal(t, 200, resp.Status())
	//
	//			var existingResult libraryPanelDashboardsResult
	//			err = json.Unmarshal(resp.Body(), &existingResult)
	//			require.NoError(t, err)
	//			require.Len(t, existingResult.Result, 1)
	//			require.Equal(t, int64(1), existingResult.Result[0])
	//
	//			sc.reqContext.ReplaceAllParams(map[string]string{":uid": unused.Result.UID})
	//			resp = sc.service.getConnectedDashboardsHandler(sc.reqContext)
	//			require.Equal(t, 200, resp.Status())
	//
	//			var unusedResult libraryPanelDashboardsResult
	//			err = json.Unmarshal(resp.Body(), &unusedResult)
	//			require.NoError(t, err)
	//			require.Len(t, unusedResult.Result, 0)
	//		})
}

type libraryPanel struct {
	ID          int64
	OrgID       int64
	FolderID    int64
	UID         string
	Name        string
	Type        string
	Description string
	Model       map[string]interface{}
	Version     int64
	Meta        libraryelements.LibraryElementDTOMeta
}

type libraryPanelResult struct {
	Result libraryPanel `json:"result"`
}

//type libraryPanelDashboardsResult struct {
//	Result []int64 `json:"result"`
//}

func overrideLibraryServicesInRegistry(cfg *setting.Cfg) (*LibraryPanelService, *libraryelements.LibraryElementService) {
	les := libraryelements.LibraryElementService{
		SQLStore: nil,
		Cfg:      cfg,
	}

	elementsOverride := func(d registry.Descriptor) (*registry.Descriptor, bool) {
		descriptor := registry.Descriptor{
			Name:     "LibraryElementService",
			Instance: &les,
		}

		return &descriptor, true
	}

	registry.RegisterOverride(elementsOverride)

	lps := LibraryPanelService{
		SQLStore:              nil,
		Cfg:                   cfg,
		LibraryElementService: &les,
	}

	panelsOverride := func(d registry.Descriptor) (*registry.Descriptor, bool) {
		descriptor := registry.Descriptor{
			Name:     "LibraryPanelService",
			Instance: &lps,
		}

		return &descriptor, true
	}

	registry.RegisterOverride(panelsOverride)

	return &lps, &les
}

type scenarioContext struct {
	ctx            *macaron.Context
	service        *LibraryPanelService
	elementService *libraryelements.LibraryElementService
	reqContext     *models.ReqContext
	user           models.SignedInUser
	folder         *models.Folder
	initialResult  libraryPanelResult
	sqlStore       *sqlstore.SQLStore
}

type folderACLItem struct {
	roleType   models.RoleType
	permission models.PermissionType
}

func createDashboard(t *testing.T, sqlStore *sqlstore.SQLStore, user models.SignedInUser, dash *models.Dashboard, folderID int64) *models.Dashboard {
	dash.FolderId = folderID
	dashItem := &dashboards.SaveDashboardDTO{
		Dashboard: dash,
		Message:   "",
		OrgId:     user.OrgId,
		User:      &user,
		Overwrite: false,
	}
	origUpdateAlerting := dashboards.UpdateAlerting
	t.Cleanup(func() {
		dashboards.UpdateAlerting = origUpdateAlerting
	})
	dashboards.UpdateAlerting = func(store dboards.Store, orgID int64, dashboard *models.Dashboard,
		user *models.SignedInUser) error {
		return nil
	}

	dashboard, err := dashboards.NewService(sqlStore).SaveDashboard(dashItem, true)
	require.NoError(t, err)

	return dashboard
}

func createFolderWithACL(t *testing.T, sqlStore *sqlstore.SQLStore, title string, user models.SignedInUser,
	items []folderACLItem) *models.Folder {
	t.Helper()

	s := dashboards.NewFolderService(user.OrgId, &user, sqlStore)
	t.Logf("Creating folder with title and UID %q", title)
	folder, err := s.CreateFolder(title, title)
	require.NoError(t, err)

	updateFolderACL(t, sqlStore, folder.Id, items)

	return folder
}

func updateFolderACL(t *testing.T, sqlStore *sqlstore.SQLStore, folderID int64, items []folderACLItem) {
	t.Helper()

	if len(items) == 0 {
		return
	}

	var aclItems []*models.DashboardAcl
	for _, item := range items {
		role := item.roleType
		permission := item.permission
		aclItems = append(aclItems, &models.DashboardAcl{
			DashboardID: folderID,
			Role:        &role,
			Permission:  permission,
			Created:     time.Now(),
			Updated:     time.Now(),
		})
	}

	err := sqlStore.UpdateDashboardACL(folderID, aclItems)
	require.NoError(t, err)
}

//func validateAndUnMarshalResponse(t *testing.T, resp response.Response) libraryPanelResult {
//	t.Helper()
//
//	require.Equal(t, 200, resp.Status())
//
//	var result = libraryPanelResult{}
//	err := json.Unmarshal(resp.Body(), &result)
//	require.NoError(t, err)
//
//	return result
//}

func scenarioWithLibraryPanel(t *testing.T, desc string, fn func(t *testing.T, sc scenarioContext)) {
	t.Helper()

	testScenario(t, desc, func(t *testing.T, sc scenarioContext) {
		command := libraryelements.CreateLibraryElementCommand{
			FolderID: sc.folder.Id,
			Name:     "Text - Library Panel",
			Model: []byte(`
			{
			  "datasource": "${DS_GDEV-TESTDATA}",
			  "id": 1,
			  "title": "Text - Library Panel",
			  "type": "text",
			  "description": "A description"
			}
		`),
			Kind: int64(libraryelements.Panel),
		}
		resp, err := sc.elementService.CreateElement(sc.reqContext, command)
		require.NoError(t, err)
		var model map[string]interface{}
		err = json.Unmarshal(resp.Model, &model)
		require.NoError(t, err)

		sc.initialResult = libraryPanelResult{
			Result: libraryPanel{
				ID:          resp.ID,
				OrgID:       resp.OrgID,
				FolderID:    resp.FolderID,
				UID:         resp.UID,
				Name:        resp.Name,
				Type:        resp.Type,
				Description: resp.Description,
				Model:       model,
				Version:     resp.Version,
				Meta:        resp.Meta,
			},
		}

		fn(t, sc)
	})
}

// testScenario is a wrapper around t.Run performing common setup for library panel tests.
// It takes your real test function as a callback.
func testScenario(t *testing.T, desc string, fn func(t *testing.T, sc scenarioContext)) {
	t.Helper()

	t.Run(desc, func(t *testing.T) {
		t.Cleanup(registry.ClearOverrides)

		ctx := macaron.Context{
			Req: macaron.Request{Request: &http.Request{}},
		}
		orgID := int64(1)
		role := models.ROLE_ADMIN

		cfg := setting.NewCfg()
		// Everything in this service is behind the feature toggle "panelLibrary"
		cfg.FeatureToggles = map[string]bool{"panelLibrary": true}
		// Because the LibraryPanelService is behind a feature toggle, we need to override the service in the registry
		// with a Cfg that contains the feature toggle so migrations are run properly
		service, elementService := overrideLibraryServicesInRegistry(cfg)

		// We need to assign SQLStore after the override and migrations are done
		sqlStore := sqlstore.InitTestDB(t)
		elementService.SQLStore = sqlStore
		service.SQLStore = sqlStore

		user := models.SignedInUser{
			UserId:     1,
			Name:       "Signed In User",
			Login:      "signed_in_user",
			Email:      "signed.in.user@test.com",
			OrgId:      orgID,
			OrgRole:    role,
			LastSeenAt: time.Now(),
		}

		// deliberate difference between signed in user and user in db to make it crystal clear
		// what to expect in the tests
		// In the real world these are identical
		cmd := models.CreateUserCommand{
			Email: "user.in.db@test.com",
			Name:  "User In DB",
			Login: UserInDbName,
		}
		_, err := sqlStore.CreateUser(context.Background(), cmd)
		require.NoError(t, err)

		sc := scenarioContext{
			user:           user,
			ctx:            &ctx,
			service:        service,
			elementService: elementService,
			sqlStore:       sqlStore,
			reqContext: &models.ReqContext{
				Context:      &ctx,
				SignedInUser: &user,
			},
		}

		sc.folder = createFolderWithACL(t, sc.sqlStore, "ScenarioFolder", sc.user, []folderACLItem{})

		fn(t, sc)
	})
}

func getCompareOptions() []cmp.Option {
	return []cmp.Option{
		cmp.Transformer("Time", func(in time.Time) int64 {
			return in.UTC().Unix()
		}),
	}
}
