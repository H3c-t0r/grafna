package config

import (
	"github.com/grafana/grafana-azure-sdk-go/azsettings"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/log"
	"github.com/grafana/grafana/pkg/setting"
)

type Cfg struct {
	log log.Logger

	DevMode bool

	PluginsPath string

	PluginSettings       setting.PluginSettings
	PluginsAllowUnsigned []string

	// AWS Plugin Auth
	AWSAllowedAuthProviders []string
	AWSAssumeRoleEnabled    bool
	AWSExternalId           string

	// Azure Cloud settings
	Azure *azsettings.AzureSettings

	// Proxy Settings
	ProxySettings setting.SecureSocksDSProxySettings

	BuildVersion string // TODO Remove

	LogDatasourceRequests bool

	PluginsCDNURLTemplate string

	Tracing Tracing

	GrafanaComURL string

	GrafanaAppURL string

	Features plugins.FeatureToggles

	AngularSupportEnabled bool

	LoadExternalPlugins map[string]bool
}

func NewCfg(devMode bool, pluginsPath string, pluginSettings setting.PluginSettings, pluginsAllowUnsigned []string,
	awsAllowedAuthProviders []string, awsAssumeRoleEnabled bool, awsExternalId string, azure *azsettings.AzureSettings, secureSocksDSProxy setting.SecureSocksDSProxySettings,
	grafanaVersion string, logDatasourceRequests bool, pluginsCDNURLTemplate string, appURL string, tracing Tracing, features plugins.FeatureToggles, angularSupportEnabled bool,
	grafanaComURL string, loadExternalPlugins []string) *Cfg {
	loadExternalPluginsMap := make(map[string]bool)
	for _, pluginID := range loadExternalPlugins {
		loadExternalPluginsMap[pluginID] = true
	}
	return &Cfg{
		log:                     log.New("plugin.cfg"),
		PluginsPath:             pluginsPath,
		BuildVersion:            grafanaVersion,
		DevMode:                 devMode,
		PluginSettings:          pluginSettings,
		PluginsAllowUnsigned:    pluginsAllowUnsigned,
		AWSAllowedAuthProviders: awsAllowedAuthProviders,
		AWSAssumeRoleEnabled:    awsAssumeRoleEnabled,
		AWSExternalId:           awsExternalId,
		Azure:                   azure,
		ProxySettings:           secureSocksDSProxy,
		LogDatasourceRequests:   logDatasourceRequests,
		PluginsCDNURLTemplate:   pluginsCDNURLTemplate,
		Tracing:                 tracing,
		GrafanaComURL:           grafanaComURL,
		GrafanaAppURL:           appURL,
		Features:                features,
		AngularSupportEnabled:   angularSupportEnabled,
		LoadExternalPlugins:     loadExternalPluginsMap,
	}
}
