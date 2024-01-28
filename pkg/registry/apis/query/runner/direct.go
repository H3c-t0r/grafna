package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/grafana/grafana/pkg/apis/query/v0alpha1"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/grafana-apiserver/endpoints/request"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/plugincontext"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/pluginstore"
	"github.com/grafana/grafana/pkg/setting"
)

type directRunner struct {
	pluginClient plugins.Client
	pCtxProvider *plugincontext.Provider
}

type directRegistry struct {
	pluginsMu     sync.Mutex
	plugins       *v0alpha1.DataSourcePluginList
	apis          map[string]schema.GroupVersion
	groupToPlugin map[string]string
	pluginStore   pluginstore.Store

	// called on demand
	dataSourcesService datasources.DataSourceService
}

var _ QueryRunner = (*directRunner)(nil)
var _ DataSourceRegistry = (*directRegistry)(nil)

// NewDummyTestRunner creates a runner that only works with testdata
func NewDirectQueryRunner(
	pluginClient plugins.Client,
	pCtxProvider *plugincontext.Provider) QueryRunner {
	return &directRunner{
		pluginClient: pluginClient,
		pCtxProvider: pCtxProvider,
	}
}

func NewDirectRegistry(pluginStore pluginstore.Store,
	dataSourcesService datasources.DataSourceService,
) DataSourceRegistry {
	return &directRegistry{
		pluginStore:        pluginStore,
		dataSourcesService: dataSourcesService,
	}
}

// ExecuteQueryData implements QueryHelper.
func (d *directRunner) ExecuteQueryData(ctx context.Context,
	// The k8s group for the datasource (pluginId)
	datasource schema.GroupVersion,

	// The datasource name/uid
	name string,

	// The raw backend query objects
	query []backend.DataQuery,
) (*backend.QueryDataResponse, error) {
	// NOTE: this depends on uid unique across datasources
	settings, err := d.pCtxProvider.GetDataSourceInstanceSettings(ctx, name)
	if err != nil {
		return nil, err
	}

	pCtx, err := d.pCtxProvider.PluginContextForDataSource(ctx, settings)
	if err != nil {
		return nil, err
	}

	return d.pluginClient.QueryData(ctx, &backend.QueryDataRequest{
		PluginContext: pCtx,
		Queries:       query,
	})
}

// GetDatasourceAPI implements DataSourceRegistry.
func (d *directRegistry) GetDatasourceAPI(pluginId string) (schema.GroupVersion, error) {
	d.pluginsMu.Lock()
	defer d.pluginsMu.Unlock()

	if d.plugins == nil {
		err := d.updatePlugins()
		if err != nil {
			return schema.GroupVersion{}, err
		}
	}

	var err error
	gv, ok := d.apis[pluginId]
	if !ok {
		err = fmt.Errorf("no API found for id: " + pluginId)
	}
	return gv, err
}

// GetDataSources implements QueryHelper.
func (d *directRegistry) GetDataSources(ctx context.Context, namespace string, options *internalversion.ListOptions) (*v0alpha1.DataSourceList, error) {
	info, err := request.NamespaceInfoFrom(ctx, true)
	if err != nil {
		return nil, err
	}
	datasources, err := d.dataSourcesService.GetDataSources(ctx, &datasources.GetDataSourcesQuery{
		OrgID: info.OrgID,
	})
	if err != nil {
		return nil, err
	}

	results := &v0alpha1.DataSourceList{
		ListMeta: metav1.ListMeta{
			ResourceVersion: fmt.Sprintf("%d", time.Now().UnixMilli()),
		},
	}
	for _, real := range datasources {
		gv := d.apis[real.Type]

		ds := v0alpha1.DataSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:              real.UID,
				CreationTimestamp: metav1.NewTime(real.Created),
			},
			Title: real.Name,
			Group: gv.Group,
			// Health: &v0alpha1.HealthCheck{
			// 	Status:  "OK",
			// 	Checked: time.Now().UnixMilli(),
			// },
		}

		results.Items = append(results.Items, ds)
	}

	return results, nil
}

// GetDatasourcePlugins no namespace? everything that is available
func (d *directRegistry) GetDatasourcePlugins(ctx context.Context, options *internalversion.ListOptions) (*v0alpha1.DataSourcePluginList, error) {
	d.pluginsMu.Lock()
	defer d.pluginsMu.Unlock()

	if d.plugins == nil {
		err := d.updatePlugins()
		if err != nil {
			return nil, err
		}
	}

	return d.plugins, nil
}

// This should be called when plugins change
func (d *directRegistry) updatePlugins() error {
	groupToPlugin := map[string]string{}
	apis := map[string]schema.GroupVersion{}
	result := &v0alpha1.DataSourcePluginList{
		ListMeta: metav1.ListMeta{
			ResourceVersion: fmt.Sprintf("%d", time.Now().UnixMilli()),
		},
	}

	// TODO? only backend plugins
	for _, dsp := range d.pluginStore.Plugins(context.Background(), plugins.TypeDataSource) {
		ts := setting.BuildStamp * 1000
		if dsp.Info.Build.Time > 0 {
			ts = dsp.Info.Build.Time
		}

		group, err := pluginstore.GetDatasourceGroupNameFromPluginID(dsp.ID)
		if err != nil {
			return err
		}
		gv := schema.GroupVersion{Group: group, Version: "v0alpha1"} // default version
		apis[dsp.ID] = gv
		for _, alias := range dsp.AliasIDs {
			apis[alias] = gv
		}
		groupToPlugin[group] = dsp.ID

		ds := v0alpha1.DataSourcePlugin{
			ObjectMeta: metav1.ObjectMeta{
				Name:              dsp.ID,
				CreationTimestamp: metav1.NewTime(time.UnixMilli(ts)),
			},
			Title:        dsp.Name,
			AliasIDs:     dsp.AliasIDs,
			GroupVersion: gv.String(),
			Description:  dsp.Info.Description,
			Capabilities: []string{"query"}, // stream?
		}
		result.Items = append(result.Items, ds)
	}

	d.plugins = result
	d.apis = apis
	d.groupToPlugin = groupToPlugin
	return nil
}
