package plugins

import (
	"encoding/json"
	"path"

	"github.com/grafana/grafana/pkg/plugins/backendplugin"

	"github.com/grafana/grafana/pkg/util/errutil"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins/datasource/wrapper"
	"github.com/grafana/grafana/pkg/tsdb"
)

// DataSourcePlugin contains all metadata about a datasource plugin
type DataSourcePlugin struct {
	FrontendPluginBase
	Annotations  bool              `json:"annotations"`
	Metrics      bool              `json:"metrics"`
	Alerting     bool              `json:"alerting"`
	Explore      bool              `json:"explore"`
	Table        bool              `json:"tables"`
	Logs         bool              `json:"logs"`
	QueryOptions map[string]bool   `json:"queryOptions,omitempty"`
	BuiltIn      bool              `json:"builtIn,omitempty"`
	Mixed        bool              `json:"mixed,omitempty"`
	Routes       []*AppPluginRoute `json:"routes"`
	Streaming    bool              `json:"streaming"`

	Backend    bool   `json:"backend,omitempty"`
	Executable string `json:"executable,omitempty"`
	SDK        bool   `json:"sdk,omitempty"`
}

func (p *DataSourcePlugin) Load(decoder *json.Decoder, pluginDir string) error {
	if err := decoder.Decode(p); err != nil {
		return errutil.Wrapf(err, "Failed to decode datasource plugin")
	}

	if err := p.registerPlugin(pluginDir); err != nil {
		return errutil.Wrapf(err, "Failed to register plugin")
	}

	if p.Backend {
		cmd := ComposePluginStartCommmand(p.Executable)
		fullpath := path.Join(p.PluginDir, cmd)
		descriptor := backendplugin.NewBackendPluginDescriptor(p.Id, fullpath, backendplugin.PluginCallbacks{
			OnLegacyStart: p.onLegacyPluginStart,
			OnStart:       p.onPluginStart,
		})
		if err := backendplugin.Register(descriptor); err != nil {
			return errutil.Wrapf(err, "Failed to register backend plugin")
		}
	}

	DataSources[p.Id] = p
	return nil
}

func (p *DataSourcePlugin) onLegacyPluginStart(pluginID string, client *backendplugin.LegacyClient, logger log.Logger) error {
	tsdb.RegisterTsdbQueryEndpoint(pluginID, func(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
		return wrapper.NewDatasourcePluginWrapper(logger, client.DatasourcePlugin), nil
	})

	return nil
}

func (p *DataSourcePlugin) onPluginStart(pluginID string, client *backendplugin.Client, logger log.Logger) error {
	if client.BackendPlugin != nil {
		tsdb.RegisterTsdbQueryEndpoint(pluginID, func(dsInfo *models.DataSource) (tsdb.TsdbQueryEndpoint, error) {
			return wrapper.NewDatasourcePluginWrapperV2(logger, client.BackendPlugin), nil
		})
	}

	return nil
}
