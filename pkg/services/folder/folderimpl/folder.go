package folderimpl

import (
	"context"
	"errors"
	"strings"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/events"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/folder"

	"github.com/grafana/grafana/pkg/services/guardian"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/search"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

type Service struct {
	store store

	log              log.Logger
	cfg              *setting.Cfg
	dashboardService dashboards.DashboardService
	dashboardStore   dashboards.Store
	searchService    *search.SearchService
	features         *featuremgmt.FeatureManager
	permissions      accesscontrol.FolderPermissionsService

	// bus is currently used to publish events that cause scheduler to update rules.
	bus bus.Bus
}

func ProvideService(
	ac accesscontrol.AccessControl,
	bus bus.Bus,
	cfg *setting.Cfg,
	dashboardService dashboards.DashboardService,
	dashboardStore dashboards.Store,
	db db.DB, // DB for the (new) nested folder store
	features *featuremgmt.FeatureManager,
	folderPermissionsService accesscontrol.FolderPermissionsService,
	searchService *search.SearchService,
) folder.Service {
	ac.RegisterScopeAttributeResolver(dashboards.NewFolderNameScopeResolver(dashboardStore))
	ac.RegisterScopeAttributeResolver(dashboards.NewFolderIDScopeResolver(dashboardStore))
	store := ProvideStore(db, cfg, features)
	return &Service{
		cfg:              cfg,
		log:              log.New("folder-service"),
		dashboardService: dashboardService,
		dashboardStore:   dashboardStore,
		store:            store,
		searchService:    searchService,
		features:         features,
		permissions:      folderPermissionsService,
		bus:              bus,
	}
}

func (s *Service) GetFolders(ctx context.Context, user *user.SignedInUser, orgID int64, limit int64, page int64) ([]*models.Folder, error) {
	searchQuery := search.Query{
		SignedInUser: user,
		DashboardIds: make([]int64, 0),
		FolderIds:    make([]int64, 0),
		Limit:        limit,
		OrgId:        orgID,
		Type:         "dash-folder",
		Permission:   models.PERMISSION_VIEW,
		Page:         page,
	}

	if err := s.searchService.SearchHandler(ctx, &searchQuery); err != nil {
		return nil, err
	}

	folders := make([]*models.Folder, 0)

	for _, hit := range searchQuery.Result {
		folders = append(folders, &models.Folder{
			Id:    hit.ID,
			Uid:   hit.UID,
			Title: hit.Title,
		})
	}

	return folders, nil
}

func (s *Service) GetFolderByID(ctx context.Context, user *user.SignedInUser, id int64, orgID int64) (*models.Folder, error) {
	if id == 0 {
		return &models.Folder{Id: id, Title: "General"}, nil
	}

	dashFolder, err := s.dashboardStore.GetFolderByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}

	g := guardian.New(ctx, dashFolder.Id, orgID, user)
	if canView, err := g.CanView(); err != nil || !canView {
		if err != nil {
			return nil, toFolderError(err)
		}
		return nil, dashboards.ErrFolderAccessDenied
	}

	return dashFolder, nil
}

func (s *Service) GetFolderByUID(ctx context.Context, user *user.SignedInUser, orgID int64, uid string) (*models.Folder, error) {
	dashFolder, err := s.dashboardStore.GetFolderByUID(ctx, orgID, uid)
	if err != nil {
		return nil, err
	}

	g := guardian.New(ctx, dashFolder.Id, orgID, user)
	if canView, err := g.CanView(); err != nil || !canView {
		if err != nil {
			return nil, toFolderError(err)
		}
		return nil, dashboards.ErrFolderAccessDenied
	}

	return dashFolder, nil
}

func (s *Service) GetFolderByTitle(ctx context.Context, user *user.SignedInUser, orgID int64, title string) (*models.Folder, error) {
	dashFolder, err := s.dashboardStore.GetFolderByTitle(ctx, orgID, title)
	if err != nil {
		return nil, err
	}

	g := guardian.New(ctx, dashFolder.Id, orgID, user)
	if canView, err := g.CanView(); err != nil || !canView {
		if err != nil {
			return nil, toFolderError(err)
		}
		return nil, dashboards.ErrFolderAccessDenied
	}

	return dashFolder, nil
}

func (s *Service) CreateFolder(ctx context.Context, user *user.SignedInUser, orgID int64, title, uid string) (*models.Folder, error) {
	dashFolder := models.NewDashboardFolder(title)
	dashFolder.OrgId = orgID

	trimmedUID := strings.TrimSpace(uid)
	if trimmedUID == accesscontrol.GeneralFolderUID {
		return nil, dashboards.ErrFolderInvalidUID
	}

	dashFolder.SetUid(trimmedUID)
	userID := user.UserID
	if userID == 0 {
		userID = -1
	}
	dashFolder.CreatedBy = userID
	dashFolder.UpdatedBy = userID
	dashFolder.UpdateSlug()

	dto := &dashboards.SaveDashboardDTO{
		Dashboard: dashFolder,
		OrgId:     orgID,
		User:      user,
	}

	saveDashboardCmd, err := s.dashboardService.BuildSaveDashboardCommand(ctx, dto, false, false)
	if err != nil {
		return nil, toFolderError(err)
	}

	dash, err := s.dashboardStore.SaveDashboard(ctx, *saveDashboardCmd)
	if err != nil {
		return nil, toFolderError(err)
	}

	var createdFolder *models.Folder
	createdFolder, err = s.dashboardStore.GetFolderByID(ctx, orgID, dash.Id)
	if err != nil {
		return nil, err
	}

	var permissionErr error
	if !accesscontrol.IsDisabled(s.cfg) {
		var permissions []accesscontrol.SetResourcePermissionCommand
		if user.IsRealUser() && !user.IsAnonymous {
			permissions = append(permissions, accesscontrol.SetResourcePermissionCommand{
				UserID: userID, Permission: models.PERMISSION_ADMIN.String(),
			})
		}

		permissions = append(permissions, []accesscontrol.SetResourcePermissionCommand{
			{BuiltinRole: string(org.RoleEditor), Permission: models.PERMISSION_EDIT.String()},
			{BuiltinRole: string(org.RoleViewer), Permission: models.PERMISSION_VIEW.String()},
		}...)

		_, permissionErr = s.permissions.SetPermissions(ctx, orgID, createdFolder.Uid, permissions...)
	} else if s.cfg.EditorsCanAdmin && user.IsRealUser() && !user.IsAnonymous {
		permissionErr = s.MakeUserAdmin(ctx, orgID, userID, createdFolder.Id, true)
	}

	if permissionErr != nil {
		s.log.Error("Could not make user admin", "folder", createdFolder.Title, "user", userID, "error", permissionErr)
	}

	if s.features.IsEnabled(featuremgmt.FlagNestedFolders) {
		var description string
		if dash.Data != nil {
			description = dash.Data.Get("description").MustString()
		}

		_, err := s.store.Create(ctx, folder.CreateFolderCommand{
			// TODO: Today, if a UID isn't specified, the dashboard store
			// generates a new UID. The new folder store will need to do this as
			// well, but for now we take the UID from the newly created folder.
			UID:         dash.Uid,
			OrgID:       orgID,
			Title:       title,
			Description: description,
			ParentUID:   folder.RootFolderUID,
		})
		if err != nil {
			// We'll log the error and also roll back the previously-created
			// (legacy) folder.
			s.log.Error("error saving folder to nested folder store", err)
			_, err = s.DeleteFolder(ctx, user, orgID, createdFolder.Uid, true)
			if err != nil {
				s.log.Error("error deleting folder after failed save to nested folder store", err)
			}
			return createdFolder, err
		}
		// The folder UID is specified (or generated) during creation, so we'll
		// stop here and return the created model.Folder.
	}
	return createdFolder, nil
}

func (s *Service) UpdateFolder(ctx context.Context, user *user.SignedInUser, orgID int64, existingUid string, cmd *models.UpdateFolderCommand) error {
	query := models.GetDashboardQuery{OrgId: orgID, Uid: existingUid}
	if _, err := s.dashboardStore.GetDashboard(ctx, &query); err != nil {
		return toFolderError(err)
	}

	dashFolder := query.Result
	currentTitle := dashFolder.Title

	if !dashFolder.IsFolder {
		return dashboards.ErrFolderNotFound
	}

	cmd.UpdateDashboardModel(dashFolder, orgID, user.UserID)

	dto := &dashboards.SaveDashboardDTO{
		Dashboard: dashFolder,
		OrgId:     orgID,
		User:      user,
		Overwrite: cmd.Overwrite,
	}

	saveDashboardCmd, err := s.dashboardService.BuildSaveDashboardCommand(ctx, dto, false, false)
	if err != nil {
		return toFolderError(err)
	}

	dash, err := s.dashboardStore.SaveDashboard(ctx, *saveDashboardCmd)
	if err != nil {
		return toFolderError(err)
	}

	var folder *models.Folder
	folder, err = s.dashboardStore.GetFolderByID(ctx, orgID, dash.Id)
	if err != nil {
		return err
	}
	cmd.Result = folder

	if currentTitle != folder.Title {
		if err := s.bus.Publish(ctx, &events.FolderTitleUpdated{
			Timestamp: folder.Updated,
			Title:     folder.Title,
			ID:        dash.Id,
			UID:       dash.Uid,
			OrgID:     orgID,
		}); err != nil {
			s.log.Error("failed to publish FolderTitleUpdated event", "folder", folder.Title, "user", user.UserID, "error", err)
		}
	}

	return nil
}

func (s *Service) DeleteFolder(ctx context.Context, user *user.SignedInUser, orgID int64, uid string, forceDeleteRules bool) (*models.Folder, error) {
	dashFolder, err := s.dashboardStore.GetFolderByUID(ctx, orgID, uid)
	if err != nil {
		return nil, err
	}

	guard := guardian.New(ctx, dashFolder.Id, orgID, user)
	if canSave, err := guard.CanDelete(); err != nil || !canSave {
		if err != nil {
			return nil, toFolderError(err)
		}
		return nil, dashboards.ErrFolderAccessDenied
	}

	deleteCmd := models.DeleteDashboardCommand{OrgId: orgID, Id: dashFolder.Id, ForceDeleteFolderRules: forceDeleteRules}

	if err := s.dashboardStore.DeleteDashboard(ctx, &deleteCmd); err != nil {
		return nil, toFolderError(err)
	}

	return dashFolder, nil
}

func (s *Service) Create(ctx context.Context, cmd *folder.CreateFolderCommand) (*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store
	if cmd.UID == "" {
		cmd.UID = util.GenerateShortUID()
	}
	return s.store.Create(ctx, *cmd)
}

func (s *Service) Update(ctx context.Context, cmd *folder.UpdateFolderCommand) (*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store
	return s.store.Update(ctx, *cmd)
}

func (s *Service) Move(ctx context.Context, cmd *folder.MoveFolderCommand) (*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store

	foldr, err := s.Get(ctx, &folder.GetFolderQuery{
		UID:   &cmd.UID,
		OrgID: cmd.OrgID,
	})
	if err != nil {
		return nil, err
	}

	return s.store.Update(ctx, folder.UpdateFolderCommand{
		Folder:       foldr,
		NewParentUID: &cmd.NewParentUID,
	})
}

func (s *Service) Delete(ctx context.Context, cmd *folder.DeleteFolderCommand) (*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store
	// check if dashboard exists

	foldr, err := s.Get(ctx, &folder.GetFolderQuery{
		UID:   &cmd.UID,
		OrgID: cmd.OrgID,
	})
	if err != nil {
		return nil, err
	}
	err = s.store.Delete(ctx, cmd.UID, cmd.OrgID)
	if err != nil {
		return nil, err
	}
	return foldr, nil
}

func (s *Service) Get(ctx context.Context, cmd *folder.GetFolderQuery) (*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store
	return s.store.Get(ctx, *cmd)
}

func (s *Service) GetParents(ctx context.Context, cmd *folder.GetParentsQuery) ([]*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store
	return s.store.GetParents(ctx, *cmd)
}

func (s *Service) GetTree(ctx context.Context, cmd *folder.GetTreeQuery) ([]*folder.Folder, error) {
	// check the flag, if old - do whatever did before
	//  for new only the store
	return s.store.GetChildren(ctx, *cmd)
}

func (s *Service) MakeUserAdmin(ctx context.Context, orgID int64, userID, folderID int64, setViewAndEditPermissions bool) error {
	return s.dashboardService.MakeUserAdmin(ctx, orgID, userID, folderID, setViewAndEditPermissions)
}

func toFolderError(err error) error {
	if errors.Is(err, dashboards.ErrDashboardTitleEmpty) {
		return dashboards.ErrFolderTitleEmpty
	}

	if errors.Is(err, dashboards.ErrDashboardUpdateAccessDenied) {
		return dashboards.ErrFolderAccessDenied
	}

	if errors.Is(err, dashboards.ErrDashboardWithSameNameInFolderExists) {
		return dashboards.ErrFolderSameNameExists
	}

	if errors.Is(err, dashboards.ErrDashboardWithSameUIDExists) {
		return dashboards.ErrFolderWithSameUIDExists
	}

	if errors.Is(err, dashboards.ErrDashboardVersionMismatch) {
		return dashboards.ErrFolderVersionMismatch
	}

	if errors.Is(err, dashboards.ErrDashboardNotFound) {
		return dashboards.ErrFolderNotFound
	}

	if errors.Is(err, dashboards.ErrDashboardFailedGenerateUniqueUid) {
		err = dashboards.ErrFolderFailedGenerateUniqueUid
	}

	return err
}
