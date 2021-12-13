package plugins

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
)

// Store is the storage for plugins.
type Store interface {
	// Plugin finds a plugin by its ID.
	Plugin(ctx context.Context, pluginID string) (PluginDTO, bool)
	// Plugins returns plugins by their requested type.
	Plugins(ctx context.Context, pluginTypes ...Type) []PluginDTO

	// Add adds a plugin from the repository to the store.
	Add(ctx context.Context, pluginID, version string, repo Repository) error
	// Remove removes a plugin from the store.
	Remove(ctx context.Context, pluginID string) error
}

// Loader is responsible for loading plugins from the file system.
type Loader interface {
	// Load will return a list of plugins found in the provided file system paths.
	Load(paths []string, ignore map[string]struct{}) ([]*Plugin, error)
	// LoadWithFactory will return a plugin found in the provided file system path and use the provided factory to
	// construct the plugin backend client.
	LoadWithFactory(path string, factory backendplugin.PluginFactoryFunc) (*Plugin, error)
}

// Repository is responsible for retrieving plugin information from a repository.
type Repository interface {
	// Download downloads the requested plugin archive.
	Download(ctx context.Context, pluginID, version string) (*PluginArchiveInfo, error)
	// GetDownloadOptions provides information for downloading the requested plugin.
	GetDownloadOptions(ctx context.Context, pluginID, version string) (*PluginDownloadOptions, error)
	// DownloadWithURL downloads the requested plugin from the specified URL.
	DownloadWithURL(ctx context.Context, pluginID, archiveURL string) (*PluginArchiveInfo, error)
}

// Client is used to communicate with backend plugin implementations.
type Client interface {
	backend.QueryDataHandler
	backend.CheckHealthHandler
	backend.StreamHandler

	// CallResource calls a plugin resource.
	CallResource(pCtx backend.PluginContext, ctx *models.ReqContext, path string)
	// CollectMetrics collects metrics from a plugin.
	CollectMetrics(ctx context.Context, pluginID string) (*backend.CollectMetricsResult, error)
}

type RendererManager interface {
	// Renderer returns a renderer plugin.
	Renderer() *Plugin
}

type CoreBackendRegistrar interface {
	// LoadAndRegister loads and registers a Core backend plugin
	LoadAndRegister(pluginID string, factory backendplugin.PluginFactoryFunc) error
}

type StaticRouteResolver interface {
	Routes() []*StaticRoute
}

type ErrorResolver interface {
	PluginErrors() []*Error
}

type PluginLoaderAuthorizer interface {
	// CanLoadPlugin confirms if a plugin is authorized to load
	CanLoadPlugin(plugin *Plugin) bool
}

type PluginDashboardManager interface {
	// GetPluginDashboards gets dashboards for a certain org/plugin.
	GetPluginDashboards(ctx context.Context, orgID int64, pluginID string) ([]*PluginDashboardInfoDTO, error)
	// LoadPluginDashboard loads a plugin dashboard.
	LoadPluginDashboard(ctx context.Context, pluginID, path string) (*models.Dashboard, error)
	// ImportDashboard imports a dashboard.
	ImportDashboard(ctx context.Context, pluginID, path string, orgID, folderID int64, dashboardModel *simplejson.Json,
		overwrite bool, inputs []ImportDashboardInput, user *models.SignedInUser) (PluginDashboardInfoDTO,
		*models.Dashboard, error)
}

type ImportDashboardInput struct {
	Type     string `json:"type"`
	PluginId string `json:"pluginId"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}
