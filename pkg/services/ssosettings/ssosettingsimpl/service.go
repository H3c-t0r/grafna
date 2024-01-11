package ssosettingsimpl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/log"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/secrets"
	"github.com/grafana/grafana/pkg/services/ssosettings"
	"github.com/grafana/grafana/pkg/services/ssosettings/api"
	"github.com/grafana/grafana/pkg/services/ssosettings/database"
	"github.com/grafana/grafana/pkg/services/ssosettings/models"
	"github.com/grafana/grafana/pkg/services/ssosettings/strategies"
	"github.com/grafana/grafana/pkg/setting"
)

var _ ssosettings.Service = (*SSOSettingsService)(nil)

type SSOSettingsService struct {
	log     log.Logger
	cfg     *setting.Cfg
	store   ssosettings.Store
	ac      ac.AccessControl
	secrets secrets.Service

	fbStrategies []ssosettings.FallbackStrategy
	reloadables  map[string]ssosettings.Reloadable
}

func ProvideService(cfg *setting.Cfg, sqlStore db.DB, ac ac.AccessControl,
	routeRegister routing.RouteRegister, features featuremgmt.FeatureToggles,
	secrets secrets.Service) *SSOSettingsService {
	strategies := []ssosettings.FallbackStrategy{
		strategies.NewOAuthStrategy(cfg),
		// register other strategies here, for example SAML
	}

	store := database.ProvideStore(sqlStore)

	svc := &SSOSettingsService{
		log:          log.New("ssosettings.service"),
		cfg:          cfg,
		store:        store,
		ac:           ac,
		fbStrategies: strategies,
		secrets:      secrets,
		reloadables:  make(map[string]ssosettings.Reloadable),
	}

	if features.IsEnabledGlobally(featuremgmt.FlagSsoSettingsApi) {
		ssoSettingsApi := api.ProvideApi(svc, routeRegister, ac)
		ssoSettingsApi.RegisterAPIEndpoints()
	}

	return svc
}

var _ ssosettings.Service = (*SSOSettingsService)(nil)

func (s *SSOSettingsService) GetForProvider(ctx context.Context, provider string) (*models.SSOSettings, error) {
	dbSettings, err := s.store.Get(ctx, provider)
	if err != nil && !errors.Is(err, ssosettings.ErrNotFound) {
		return nil, err
	}

	systemSettings, err := s.loadSettingsUsingFallbackStrategy(ctx, provider)
	if err != nil {
		return nil, err
	}

	// storedSettings := map[string]any{}
	// if dbSettings != nil {
	// 	storedSettings = dbSettings.Settings
	// }

	// finalSettings := mergeSettings(storedSettings, systemSettings.Settings)

	// if dataFromDB {
	// 	dbSettings.Settings = finalSettings
	// 	return dbSettings, nil
	// }

	// // if errors.Is(err, ssosettings.ErrNotFound) {
	// // 	settings, err := s.loadSettingsUsingFallbackStrategy(ctx, provider)
	// // 	if err != nil {
	// // 		return nil, err
	// // 	}

	// // 	return settings, nil
	// // }

	// if err != nil {
	// 	return nil, err
	// }

	// storeSettings.Source = models.DB

	return mergeSSOSettings(dbSettings, systemSettings), nil
}

func (s *SSOSettingsService) GetForProviderWithRedactedSecrets(ctx context.Context, provider string) (*models.SSOSettings, error) {
	storeSettings, err := s.GetForProvider(ctx, provider)
	if err != nil {
		return nil, err
	}

	for k, v := range storeSettings.Settings {
		if isSecret(k) && v != "" {
			storeSettings.Settings[k] = "*********"
		}
	}

	return storeSettings, nil
}

func (s *SSOSettingsService) List(ctx context.Context) ([]*models.SSOSettings, error) {
	result := make([]*models.SSOSettings, 0, len(ssosettings.AllOAuthProviders))
	storedSettings, err := s.store.List(ctx)

	if err != nil {
		return nil, err
	}

	for _, provider := range ssosettings.AllOAuthProviders {
		dbSettings := getSettingByProvider(provider, storedSettings)
		fallbackSettings, err := s.loadSettingsUsingFallbackStrategy(ctx, provider)
		if err != nil {
			return nil, err
		}

		result = append(result, mergeSSOSettings(dbSettings, fallbackSettings))
	}

	return result, nil
}

func (s *SSOSettingsService) Upsert(ctx context.Context, settings models.SSOSettings) error {
	if !isProviderConfigurable(settings.Provider) {
		return ssosettings.ErrInvalidProvider.Errorf("provider %s is not configurable", settings.Provider)
	}

	social, ok := s.reloadables[settings.Provider]
	if !ok {
		return ssosettings.ErrInvalidProvider.Errorf("provider %s not found in reloadables", settings.Provider)
	}

	err := social.Validate(ctx, settings)
	if err != nil {
		return err
	}

	// systemSettings, err := s.loadSettingsUsingFallbackStrategy(ctx, settings.Provider)
	// if err != nil {
	// 	return err
	// }

	// add the SSO settings from system that are not available in the user input
	// in order to have a complete set of SSO settings for every provider in the database
	// settings.Settings = mergeSettings(settings.Settings, systemSettings.Settings)

	settings.Settings, err = s.encryptSecrets(ctx, settings.Settings)
	if err != nil {
		return err
	}

	err = s.store.Upsert(ctx, settings)
	if err != nil {
		return err
	}

	go func() {
		err = social.Reload(context.Background(), settings)
		if err != nil {
			s.log.Error("failed to reload the provider", "provider", settings.Provider, "error", err)
		}
	}()

	return nil
}

func (s *SSOSettingsService) Patch(ctx context.Context, provider string, data map[string]any) error {
	panic("not implemented") // TODO: Implement
}

func (s *SSOSettingsService) Delete(ctx context.Context, provider string) error {
	return s.store.Delete(ctx, provider)
}

func (s *SSOSettingsService) Reload(ctx context.Context, provider string) {
	panic("not implemented") // TODO: Implement
}

func (s *SSOSettingsService) RegisterReloadable(provider string, reloadable ssosettings.Reloadable) {
	if s.reloadables == nil {
		s.reloadables = make(map[string]ssosettings.Reloadable)
	}
	s.reloadables[provider] = reloadable
}

func (s *SSOSettingsService) RegisterFallbackStrategy(providerRegex string, strategy ssosettings.FallbackStrategy) {
	s.fbStrategies = append(s.fbStrategies, strategy)
}

func (s *SSOSettingsService) loadSettingsUsingFallbackStrategy(ctx context.Context, provider string) (*models.SSOSettings, error) {
	loadStrategy, ok := s.getFallbackStrategyFor(provider)
	if !ok {
		return nil, errors.New("no fallback strategy found for provider: " + provider)
	}

	settingsFromSystem, err := loadStrategy.GetProviderConfig(ctx, provider)
	if err != nil {
		return nil, err
	}

	return &models.SSOSettings{
		Provider: provider,
		Source:   models.System,
		Settings: settingsFromSystem,
	}, nil
}

func getSettingByProvider(provider string, settings []*models.SSOSettings) *models.SSOSettings {
	for _, item := range settings {
		if item.Provider == provider {
			return item
		}
	}
	return nil
}

func (s *SSOSettingsService) getFallbackStrategyFor(provider string) (ssosettings.FallbackStrategy, bool) {
	for _, strategy := range s.fbStrategies {
		if strategy.IsMatch(provider) {
			return strategy, true
		}
	}
	return nil, false
}

func (s *SSOSettingsService) encryptSecrets(ctx context.Context, settings map[string]any) (map[string]any, error) {
	result := make(map[string]any)
	for k, v := range settings {
		if isSecret(k) {
			strValue, ok := v.(string)
			if !ok {
				return result, fmt.Errorf("failed to encrypt %s setting because it is not a string: %v", k, v)
			}

			encryptedSecret, err := s.secrets.Encrypt(ctx, []byte(strValue), secrets.WithoutScope())
			if err != nil {
				return result, err
			}
			result[k] = string(encryptedSecret)
		} else {
			result[k] = v
		}
	}

	return result, nil
}

func (s *SSOSettingsService) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)

	// start a background process for reloading the SSO settings for all providers at a fixed interval
	// it is useful for high availability setups running multiple Grafana instances
	for {
		select {
		case <-ticker.C:
			s.doReload(ctx)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *SSOSettingsService) doReload(ctx context.Context) {
	s.log.Debug("reloading SSO Settings for all providers")

	settingsList, err := s.List(ctx)
	if err != nil {
		s.log.Error("failed to fetch SSO Settings for all providers", "err", err)
		return
	}

	for provider, connector := range s.reloadables {
		setting := getSettingByProvider(provider, settingsList)

		err = connector.Reload(ctx, *setting)
		if err != nil {
			s.log.Error("failed to reload SSO Settings", "provider", provider, "err", err)
			continue
		}
	}
}

func isSecret(fieldName string) bool {
	secretFieldPatterns := []string{"secret"}

	for _, v := range secretFieldPatterns {
		if strings.Contains(strings.ToLower(fieldName), strings.ToLower(v)) {
			return true
		}
	}
	return false
}

func mergeSettings(apiSettings, systemSettings map[string]any) map[string]any {
	settings := make(map[string]any)

	for k, v := range apiSettings {
		settings[k] = v
	}

	for k, v := range systemSettings {
		if _, ok := settings[k]; !ok {
			settings[k] = v
		}
	}

	return settings
}

func mergeSSOSettings(dbSettings, systemSettings *models.SSOSettings) *models.SSOSettings {
	storedSettings := map[string]any{}
	if dbSettings != nil {
		storedSettings = dbSettings.Settings
	}

	finalSettings := mergeSettings(storedSettings, systemSettings.Settings)

	if dbSettings != nil {
		dbSettings.Settings = finalSettings
		return dbSettings
	}

	return &models.SSOSettings{
		Provider: systemSettings.Provider,
		Settings: finalSettings,
		Source:   models.System,
	}
}

func isProviderConfigurable(provider string) bool {
	for _, configurable := range ssosettings.ConfigurableOAuthProviders {
		if provider == configurable {
			return true
		}
	}

	return false
}
