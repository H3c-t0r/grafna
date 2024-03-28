package cloudmigrationimpl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/cloudmigration"
	"github.com/grafana/grafana/pkg/services/cloudmigration/api"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/folder"
	"github.com/grafana/grafana/pkg/services/gcom"
	"github.com/grafana/grafana/pkg/services/secrets"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

// Service Define the cloudmigration.Service Implementation.
type Service struct {
	store store

	log *log.ConcreteLogger
	cfg *setting.Cfg

	features featuremgmt.FeatureToggles

	dsService       datasources.DataSourceService
	gcomService     gcom.Service
	dashboarService dashboards.DashboardService
	folderService   folder.Service
	secretsService  secrets.Service

	api     *api.CloudMigrationAPI
	tracer  tracing.Tracer
	metrics *Metrics
}

var LogPrefix = "cloudmigration.service"

const (
	// nolint:gosec
	cloudMigrationAccessPolicyName = "grafana-cloud-migrations"
	//nolint:gosec
	cloudMigrationTokenName = "grafana-cloud-migrations"
)

var _ cloudmigration.Service = (*Service)(nil)

// ProvideService Factory for method used by wire to inject dependencies.
// builds the service, and api, and configures routes
func ProvideService(
	cfg *setting.Cfg,
	features featuremgmt.FeatureToggles,
	db db.DB,
	dsService datasources.DataSourceService,
	routeRegister routing.RouteRegister,
	prom prometheus.Registerer,
	tracer tracing.Tracer,
) cloudmigration.Service {
	if !features.IsEnabledGlobally(featuremgmt.FlagOnPremToCloudMigrations) {
		return &NoopServiceImpl{}
	}

	s := &Service{
		store:       &sqlStore{db: db},
		log:         log.New(LogPrefix),
		cfg:         cfg,
		features:    features,
		dsService:   dsService,
		gcomService: gcom.New(gcom.Config{ApiURL: cfg.GrafanaComAPIURL, Token: cfg.CloudMigration.GcomAPIToken}),
		tracer:      tracer,
		metrics:     newMetrics(),
	}
	s.api = api.RegisterApi(routeRegister, s, tracer)

	if err := s.registerMetrics(prom, s.metrics); err != nil {
		s.log.Warn("error registering prom metrics", "error", err.Error())
	}

	return s
}

func (s *Service) CreateToken(ctx context.Context) (cloudmigration.CreateAccessTokenResponse, error) {
	ctx, span := s.tracer.Start(ctx, "CloudMigrationService.CreateToken")
	defer span.End()
	logger := s.log.FromContext(ctx)
	requestID := tracing.TraceIDFromContext(ctx, false)

	timeoutCtx, cancel := context.WithTimeout(ctx, s.cfg.CloudMigration.FetchInstanceTimeout)
	defer cancel()
	instance, err := s.gcomService.GetInstanceByID(timeoutCtx, requestID, s.cfg.StackID)
	if err != nil {
		return cloudmigration.CreateAccessTokenResponse{}, fmt.Errorf("fetching instance by id: id=%s %w", s.cfg.StackID, err)
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, s.cfg.CloudMigration.FetchAccessPolicyTimeout)
	defer cancel()
	existingAccessPolicy, err := s.findAccessPolicyByName(timeoutCtx, instance.RegionSlug, cloudMigrationAccessPolicyName)
	if err != nil {
		return cloudmigration.CreateAccessTokenResponse{}, fmt.Errorf("fetching access policy by name: name=%s %w", cloudMigrationAccessPolicyName, err)
	}

	if existingAccessPolicy != nil {
		timeoutCtx, cancel := context.WithTimeout(ctx, s.cfg.CloudMigration.DeleteAccessPolicyTimeout)
		defer cancel()
		if _, err := s.gcomService.DeleteAccessPolicy(timeoutCtx, gcom.DeleteAccessPolicyParams{
			RequestID:      requestID,
			AccessPolicyID: existingAccessPolicy.ID,
			Region:         instance.RegionSlug,
		}); err != nil {
			return cloudmigration.CreateAccessTokenResponse{}, fmt.Errorf("deleting access policy: id=%s region=%s %w", existingAccessPolicy.ID, instance.RegionSlug, err)
		}
		logger.Info("deleted access policy", existingAccessPolicy.ID, "name", existingAccessPolicy.Name)
	}

	timeoutCtx, cancel = context.WithTimeout(ctx, s.cfg.CloudMigration.CreateAccessPolicyTimeout)
	defer cancel()
	accessPolicy, err := s.gcomService.CreateAccessPolicy(timeoutCtx,
		gcom.CreateAccessPolicyParams{
			RequestID: requestID,
			Region:    instance.RegionSlug,
		},
		gcom.CreateAccessPolicyPayload{
			Name:        cloudMigrationAccessPolicyName,
			DisplayName: cloudMigrationAccessPolicyName,
			Realms:      []gcom.Realm{{Type: "stack", Identifier: s.cfg.StackID, LabelPolicies: []gcom.LabelPolicy{}}},
			Scopes:      []string{"cloud-migrations:read", "cloud-migrations:write"},
		})
	if err != nil {
		return cloudmigration.CreateAccessTokenResponse{}, fmt.Errorf("creating access policy: %w", err)
	}
	logger.Info("created access policy", "id", accessPolicy.ID, "name", accessPolicy.Name)

	timeoutCtx, cancel = context.WithTimeout(ctx, s.cfg.CloudMigration.CreateTokenTimeout)
	defer cancel()
	token, err := s.gcomService.CreateToken(timeoutCtx,
		gcom.CreateTokenParams{RequestID: requestID, Region: instance.RegionSlug},
		gcom.CreateTokenPayload{
			AccessPolicyID: accessPolicy.ID,
			DisplayName:    cloudMigrationTokenName,
			Name:           cloudMigrationTokenName,
			ExpiresAt:      time.Now().Add(s.cfg.CloudMigration.TokenExpiresAfter),
		})
	if err != nil {
		return cloudmigration.CreateAccessTokenResponse{}, fmt.Errorf("creating access token: %w", err)
	}
	logger.Info("created access token", "id", token.ID, "name", token.Name)
	s.metrics.accessTokenCreated.With(prometheus.Labels{"slug": s.cfg.Slug}).Inc()

	bytes, err := json.Marshal(map[string]string{
		"token":  token.Token,
		"region": instance.ClusterSlug,
	})
	if err != nil {
		return cloudmigration.CreateAccessTokenResponse{}, fmt.Errorf("encoding token: %w", err)
	}

	return cloudmigration.CreateAccessTokenResponse{Token: base64.StdEncoding.EncodeToString(bytes)}, nil
}

func (s *Service) findAccessPolicyByName(ctx context.Context, regionSlug, accessPolicyName string) (*gcom.AccessPolicy, error) {
	ctx, span := s.tracer.Start(ctx, "CloudMigrationService.findAccessPolicyByName")
	defer span.End()

	accessPolicies, err := s.gcomService.ListAccessPolicies(ctx, gcom.ListAccessPoliciesParams{
		RequestID: tracing.TraceIDFromContext(ctx, false),
		Region:    regionSlug,
		Name:      accessPolicyName,
	})
	if err != nil {
		return nil, fmt.Errorf("listing access policies: name=%s region=%s :%w", accessPolicyName, regionSlug, err)
	}

	for _, accessPolicy := range accessPolicies {
		if accessPolicy.Name == accessPolicyName {
			return &accessPolicy, nil
		}
	}

	return nil, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) error {
	// TODO: Implement method
	return nil
}

func (s *Service) SaveEncryptedToken(ctx context.Context, token string) error {
	// TODO: Implement method
	return nil
}

func (s *Service) GetMigration(ctx context.Context, id int64) (*cloudmigration.CloudMigration, error) {
	migration, err := s.store.GetMigration(ctx, id)
	if err != nil {
		return nil, err
	}
	strValue := migration.AuthToken
	decoded, err := base64.RawStdEncoding.DecodeString(strValue)
	if err != nil {
		s.log.Error("Failed to decode secret string", "err", err, "value")
		return nil, err
	}

	decryptedToken, err := s.secretsService.Decrypt(ctx, decoded)
	if err != nil {
		s.log.Error("Failed to decrypt secret", "err", err)
		return nil, err
	}
	migration.AuthToken = string(decryptedToken)

	return migration, nil
}

func (s *Service) GetMigrationList(ctx context.Context) (*cloudmigration.CloudMigrationListResponse, error) {
	values, err := s.store.GetAllCloudMigrations(ctx)
	if err != nil {
		return nil, err
	}

	migrations := make([]cloudmigration.CloudMigrationResponse, 0)
	for _, v := range values {
		migrations = append(migrations, cloudmigration.CloudMigrationResponse{
			ID:      v.ID,
			Stack:   v.Stack,
			Created: v.Created,
			Updated: v.Updated,
		})
	}
	return &cloudmigration.CloudMigrationListResponse{Migrations: migrations}, nil
}

func (s *Service) CreateMigration(ctx context.Context, cm cloudmigration.CloudMigrationRequest) (*cloudmigration.CloudMigrationResponse, error) {
	// TODO: Implement method
	return nil, nil
}

func (s *Service) UpdateMigration(ctx context.Context, id int64, cm cloudmigration.CloudMigrationRequest) (*cloudmigration.CloudMigrationResponse, error) {
	// TODO: Implement method
	return nil, nil
}

func (s *Service) RunMigration(ctx context.Context, id int64) (*cloudmigration.RunMigrationResponse, error) {
	var result cloudmigration.RunMigrationResponse

	dataSources, err := s.dsService.GetAllDataSources(ctx, &datasources.GetAllDataSourcesQuery{})
	if err != nil {
		s.log.Error("Failed to get all datasources", "err", err)
		return nil, err
	}

	folders, err := s.folderService.GetFolders(ctx, folder.GetFoldersQuery{})
	if err != nil {
		s.log.Error("Failed to get all folders", "err", err)
		return nil, err
	}

	dashs, err := s.dashboarService.GetAllDashboards(ctx)
	if err != nil {
		s.log.Error("Failed to get all dashboards", "err", err)
		return nil, err
	}

	dataSourcesToMarshal := []datasources.AddDataSourceCommand{}
	for _, dataSource := range dataSources {
		// Decrypt secure json to send raw credentials
		decryptedData, err := s.secretsService.DecryptJsonData(ctx, dataSource.SecureJsonData)
		if err != nil {
			s.log.Error("Failed to decrypt secure json data", "err", err)
			return nil, err
		}
		dataSourceCmd := datasources.AddDataSourceCommand{
			OrgID:           dataSource.OrgID,
			Name:            dataSource.Name,
			Type:            dataSource.Type,
			Access:          dataSource.Access,
			URL:             dataSource.URL,
			User:            dataSource.User,
			Database:        dataSource.Database,
			BasicAuth:       dataSource.BasicAuth,
			BasicAuthUser:   dataSource.BasicAuthUser,
			WithCredentials: dataSource.WithCredentials,
			IsDefault:       dataSource.IsDefault,
			JsonData:        dataSource.JsonData,
			SecureJsonData:  decryptedData,
			ReadOnly:        dataSource.ReadOnly,
			UID:             dataSource.UID,
		}
		dataSourcesToMarshal = append(dataSourcesToMarshal, dataSourceCmd)
	}
	result.DataSources, err = json.Marshal(dataSourcesToMarshal)
	if err != nil {
		s.log.Error("Failed to marshal datasources", "err", err)
		return nil, err
	}

	var allFoldersToMarshal []folder.Folder
	for _, folder := range folders {
		allFoldersToMarshal = append(allFoldersToMarshal, *folder)
	}
	result.Folders, err = json.Marshal(allFoldersToMarshal)
	if err != nil {
		s.log.Error("Failed to marshal folders", "err", err)
		return nil, err
	}

	var allDashboardsToMarshal []dashboards.Dashboard
	for _, dashboard := range dashs {
		allDashboardsToMarshal = append(allDashboardsToMarshal, *dashboard)
	}
	result.Dashboards, err = json.Marshal(allDashboardsToMarshal)
	if err != nil {
		s.log.Error("Failed to marshal dashboards", "err", err)
		return nil, err
	}

	// save migration run in DB table
	err = s.SaveMigrationRun(ctx, &cloudmigration.CloudMigrationRun{
		ID: id,
		Result: cloudmigration.MigrationResult{
			Status:  "success",
			Message: "Migration run successful",
		},
	})
	if err != nil {
		s.log.Error("Failed to save migration run", "err", err)
		return &result, err
	}

	return &result, nil
}

func (s *Service) SaveMigrationRun(ctx context.Context, cmr *cloudmigration.CloudMigrationRun) error {
	cmr.CloudMigrationUID = util.GenerateShortUID()
	cmr.Created = time.Now()
	cmr.Updated = time.Now()
	cmr.Finished = time.Now()
	return s.store.SaveMigrationRun(ctx, cmr)
}

func (s *Service) GetMigrationStatus(ctx context.Context, id string, runID string) (*cloudmigration.CloudMigrationRun, error) {
	// TODO: Implement method
	return nil, nil
}

func (s *Service) GetMigrationStatusList(ctx context.Context, id string) ([]cloudmigration.CloudMigrationRun, error) {
	// TODO: Implement method
	return nil, nil
}

func (s *Service) DeleteMigration(ctx context.Context, id string) error {
	// TODO: Implement method
	return nil
}

func (s *Service) ParseCloudMigrationConfig(cfg *setting.Cfg) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("cfg cannot be nil")
	}
	section := cfg.Raw.Section("cloud_migration")
	domain := section.Key("domain").MustString("")
	if domain != "" {
		s.log.Warn("cloudmigration domain not set")
	}
	return "", nil
}
