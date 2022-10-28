// Copyright 2022 Grafana Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file is autogenerated. DO NOT EDIT.
//
// Run `make gen-cue` from repository root to regenerate.

package corelist

import (
	"fmt"
	"io/fs"

	"github.com/grafana/grafana"
	"github.com/grafana/grafana/pkg/plugins/pfs"
	"github.com/grafana/thema"
)

func makeTreeOrPanic(path string, pkgname string, rt *thema.Runtime) *pfs.Tree {
	sub, err := fs.Sub(grafana.CueSchemaFS, path)
	if err != nil {
		panic("could not create fs sub to " + path)
	}
	tree, err := pfs.ParsePluginFS(sub, rt)
	if err != nil {
		panic(fmt.Sprintf("error parsing plugin metadata for %s: %s", pkgname, err))
	}
	return tree
}

func coreTreeList(rt *thema.Runtime) pfs.TreeList {
	return pfs.TreeList{
		makeTreeOrPanic("public/app/plugins/datasource/alertmanager", "alertmanager", rt),
		makeTreeOrPanic("public/app/plugins/datasource/cloud-monitoring", "stackdriver", rt),
		makeTreeOrPanic("public/app/plugins/datasource/cloudwatch", "cloudwatch", rt),
		makeTreeOrPanic("public/app/plugins/datasource/dashboard", "dashboard", rt),
		makeTreeOrPanic("public/app/plugins/datasource/elasticsearch", "elasticsearch", rt),
		makeTreeOrPanic("public/app/plugins/datasource/grafana", "grafana", rt),
		makeTreeOrPanic("public/app/plugins/datasource/grafana-azure-monitor-datasource", "grafana_azure_monitor_datasource", rt),
		makeTreeOrPanic("public/app/plugins/datasource/graphite", "graphite", rt),
		makeTreeOrPanic("public/app/plugins/datasource/jaeger", "jaeger", rt),
		makeTreeOrPanic("public/app/plugins/datasource/loki", "loki", rt),
		makeTreeOrPanic("public/app/plugins/datasource/mssql", "mssql", rt),
		makeTreeOrPanic("public/app/plugins/datasource/mysql", "mysql", rt),
		makeTreeOrPanic("public/app/plugins/datasource/parca", "parca", rt),
		makeTreeOrPanic("public/app/plugins/datasource/phlare", "phlare", rt),
		makeTreeOrPanic("public/app/plugins/datasource/postgres", "postgres", rt),
		makeTreeOrPanic("public/app/plugins/datasource/prometheus", "prometheus", rt),
		makeTreeOrPanic("public/app/plugins/datasource/tempo", "tempo", rt),
		makeTreeOrPanic("public/app/plugins/datasource/testdata", "testdata", rt),
		makeTreeOrPanic("public/app/plugins/datasource/zipkin", "zipkin", rt),
		makeTreeOrPanic("public/app/plugins/panel/alertGroups", "alertGroups", rt),
		makeTreeOrPanic("public/app/plugins/panel/alertlist", "alertlist", rt),
		makeTreeOrPanic("public/app/plugins/panel/annolist", "annolist", rt),
		makeTreeOrPanic("public/app/plugins/panel/barchart", "barchart", rt),
		makeTreeOrPanic("public/app/plugins/panel/bargauge", "bargauge", rt),
		makeTreeOrPanic("public/app/plugins/panel/dashlist", "dashlist", rt),
		makeTreeOrPanic("public/app/plugins/panel/debug", "debug", rt),
		makeTreeOrPanic("public/app/plugins/panel/flamegraph", "flamegraph", rt),
		makeTreeOrPanic("public/app/plugins/panel/gauge", "gauge", rt),
		makeTreeOrPanic("public/app/plugins/panel/geomap", "geomap", rt),
		makeTreeOrPanic("public/app/plugins/panel/gettingstarted", "gettingstarted", rt),
		makeTreeOrPanic("public/app/plugins/panel/graph", "graph", rt),
		makeTreeOrPanic("public/app/plugins/panel/histogram", "histogram", rt),
		makeTreeOrPanic("public/app/plugins/panel/icon", "icon", rt),
		makeTreeOrPanic("public/app/plugins/panel/live", "live", rt),
		makeTreeOrPanic("public/app/plugins/panel/logs", "logs", rt),
		makeTreeOrPanic("public/app/plugins/panel/news", "news", rt),
		makeTreeOrPanic("public/app/plugins/panel/nodeGraph", "nodeGraph", rt),
		makeTreeOrPanic("public/app/plugins/panel/piechart", "piechart", rt),
		makeTreeOrPanic("public/app/plugins/panel/stat", "stat", rt),
		makeTreeOrPanic("public/app/plugins/panel/table-old", "table_old", rt),
		makeTreeOrPanic("public/app/plugins/panel/text", "text", rt),
		makeTreeOrPanic("public/app/plugins/panel/traces", "traces", rt),
		makeTreeOrPanic("public/app/plugins/panel/welcome", "welcome", rt),
		makeTreeOrPanic("public/app/plugins/panel/xychart", "xychart", rt),
	}
}
