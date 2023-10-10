// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     PluginTreeListJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package corelist

import (
	"fmt"
	"io/fs"

	"github.com/grafana/grafana"
	"github.com/grafana/grafana/pkg/plugins/pfs"
	"github.com/grafana/thema"
)

func parsePluginOrPanic(path string, pkgname string, rt *thema.Runtime) pfs.ParsedPlugin {
	sub, err := fs.Sub(grafana.CueSchemaFS, path)
	if err != nil {
		panic("could not create fs sub to " + path)
	}
	pp, err := pfs.ParsePluginFS(sub, rt)
	if err != nil {
		panic(fmt.Sprintf("error parsing plugin metadata for %s: %s", pkgname, err))
	}
	return pp
}

func corePlugins(rt *thema.Runtime) []pfs.ParsedPlugin {
	return []pfs.ParsedPlugin{
		parsePluginOrPanic("public/app/plugins/datasource/alertmanager", "alertmanager", rt),
		parsePluginOrPanic("public/app/plugins/datasource/azuremonitor", "grafana_azure_monitor_datasource", rt),
		parsePluginOrPanic("public/app/plugins/datasource/cloud-monitoring", "stackdriver", rt),
		parsePluginOrPanic("public/app/plugins/datasource/cloudwatch", "cloudwatch", rt),
		parsePluginOrPanic("public/app/plugins/datasource/dashboard", "dashboard", rt),
		parsePluginOrPanic("public/app/plugins/datasource/elasticsearch", "elasticsearch", rt),
		parsePluginOrPanic("public/app/plugins/datasource/grafana", "grafana", rt),
		parsePluginOrPanic("public/app/plugins/datasource/grafana-pyroscope-datasource", "grafana_pyroscope_datasource", rt),
		parsePluginOrPanic("public/app/plugins/datasource/grafana-testdata-datasource", "grafana_testdata_datasource", rt),
		parsePluginOrPanic("public/app/plugins/datasource/graphite", "graphite", rt),
		parsePluginOrPanic("public/app/plugins/datasource/jaeger", "jaeger", rt),
		parsePluginOrPanic("public/app/plugins/datasource/loki", "loki", rt),
		parsePluginOrPanic("public/app/plugins/datasource/mssql", "mssql", rt),
		parsePluginOrPanic("public/app/plugins/datasource/mysql", "mysql", rt),
		parsePluginOrPanic("public/app/plugins/datasource/parca", "parca", rt),
		parsePluginOrPanic("public/app/plugins/datasource/postgres", "postgres", rt),
		parsePluginOrPanic("public/app/plugins/datasource/prometheus", "prometheus", rt),
		parsePluginOrPanic("public/app/plugins/datasource/tempo", "tempo", rt),
		parsePluginOrPanic("public/app/plugins/datasource/zipkin", "zipkin", rt),
		parsePluginOrPanic("public/app/plugins/panel/alertGroups", "alertGroups", rt),
		parsePluginOrPanic("public/app/plugins/panel/alertlist", "alertlist", rt),
		parsePluginOrPanic("public/app/plugins/panel/annolist", "annolist", rt),
		parsePluginOrPanic("public/app/plugins/panel/barchart", "barchart", rt),
		parsePluginOrPanic("public/app/plugins/panel/bargauge", "bargauge", rt),
		parsePluginOrPanic("public/app/plugins/panel/candlestick", "candlestick", rt),
		parsePluginOrPanic("public/app/plugins/panel/canvas", "canvas", rt),
		parsePluginOrPanic("public/app/plugins/panel/dashlist", "dashlist", rt),
		parsePluginOrPanic("public/app/plugins/panel/datagrid", "datagrid", rt),
		parsePluginOrPanic("public/app/plugins/panel/debug", "debug", rt),
		parsePluginOrPanic("public/app/plugins/panel/flamegraph", "flamegraph", rt),
		parsePluginOrPanic("public/app/plugins/panel/gauge", "gauge", rt),
		parsePluginOrPanic("public/app/plugins/panel/geomap", "geomap", rt),
		parsePluginOrPanic("public/app/plugins/panel/gettingstarted", "gettingstarted", rt),
		parsePluginOrPanic("public/app/plugins/panel/graph", "graph", rt),
		parsePluginOrPanic("public/app/plugins/panel/heatmap", "heatmap", rt),
		parsePluginOrPanic("public/app/plugins/panel/histogram", "histogram", rt),
		parsePluginOrPanic("public/app/plugins/panel/live", "live", rt),
		parsePluginOrPanic("public/app/plugins/panel/logs", "logs", rt),
		parsePluginOrPanic("public/app/plugins/panel/news", "news", rt),
		parsePluginOrPanic("public/app/plugins/panel/nodeGraph", "nodeGraph", rt),
		parsePluginOrPanic("public/app/plugins/panel/piechart", "piechart", rt),
		parsePluginOrPanic("public/app/plugins/panel/stat", "stat", rt),
		parsePluginOrPanic("public/app/plugins/panel/state-timeline", "state_timeline", rt),
		parsePluginOrPanic("public/app/plugins/panel/status-history", "status_history", rt),
		parsePluginOrPanic("public/app/plugins/panel/table", "table", rt),
		parsePluginOrPanic("public/app/plugins/panel/table-old", "table_old", rt),
		parsePluginOrPanic("public/app/plugins/panel/text", "text", rt),
		parsePluginOrPanic("public/app/plugins/panel/timeseries", "timeseries", rt),
		parsePluginOrPanic("public/app/plugins/panel/traces", "traces", rt),
		parsePluginOrPanic("public/app/plugins/panel/trend", "trend", rt),
		parsePluginOrPanic("public/app/plugins/panel/welcome", "welcome", rt),
		parsePluginOrPanic("public/app/plugins/panel/xychart", "xychart", rt),
	}
}
