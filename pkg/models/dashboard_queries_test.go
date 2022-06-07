package models

import (
	"testing"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/stretchr/testify/require"
)

const (
	dashboardWithNoQueries = `
{
  "panels": [
    {
      "id": 2,
      "title": "Panel Title",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 35
}`

	dashboardWithQueries = `
{
  "panels": [
    {
      "id": 2,
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "_yxMP8Ynk"
          },
          "exemplar": true,
          "expr": "go_goroutines{job=\"$job\"}",
          "interval": "",
          "legendFormat": "",
          "refId": "A"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "promds2"
          },
          "exemplar": true,
          "expr": "query2",
          "interval": "",
          "legendFormat": "",
          "refId": "B"
        }
      ],
      "title": "Panel Title",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 35
}`

	oldStyleDashboard = `
{
  "panels": [
    {
	  "datasource": "_yxMP8Ynk",
      "id": 2,
      "targets": [
        {
          "exemplar": true,
          "expr": "go_goroutines{job=\"$job\"}",
          "interval": "",
          "legendFormat": "",
          "refId": "A"
        }
      ],
      "title": "Panel Title",
      "type": "timeseries"
    }
  ],
  "schemaVersion": 21
}`
)

func TestGetQueriesFromDashboard(t *testing.T) {
	t.Run("can extract no queries from empty dashboard", func(t *testing.T) {
		json, err := simplejson.NewJson([]byte(`{"panels": {}}`))
		require.NoError(t, err)

		queries := GetQueriesFromDashboard(json)
		require.Len(t, queries, 0)
	})

	t.Run("can extract no queries from empty panel", func(t *testing.T) {
		json, err := simplejson.NewJson([]byte(dashboardWithNoQueries))
		require.NoError(t, err)

		queries := GetQueriesFromDashboard(json)
		require.Len(t, queries, 1)
		require.Contains(t, queries, int64(2))
		require.Len(t, queries[2], 0)
	})

	t.Run("can extract queries from panels", func(t *testing.T) {
		json, err := simplejson.NewJson([]byte(dashboardWithQueries))
		require.NoError(t, err)

		queries := GetQueriesFromDashboard(json)
		require.Len(t, queries, 1)
		require.Contains(t, queries, int64(2))
		require.Len(t, queries[2], 2)
		query, err := queries[2][0].MarshalJSON()
		require.NoError(t, err)
		require.JSONEq(t, `{
            "datasource": {
              "type": "prometheus",
              "uid": "_yxMP8Ynk"
            },
            "exemplar": true,
            "expr": "go_goroutines{job=\"$job\"}",
            "interval": "",
            "legendFormat": "",
            "refId": "A"
		}`, string(query))
		query, err = queries[2][1].MarshalJSON()
		require.NoError(t, err)
		require.JSONEq(t, `{
            "datasource": {
              "type": "prometheus",
              "uid": "promds2"
            },
            "exemplar": true,
            "expr": "query2",
            "interval": "",
            "legendFormat": "",
            "refId": "B"
		}`, string(query))
	})

	t.Run("can extract queries from old-style panels", func(t *testing.T) {
		json, err := simplejson.NewJson([]byte(oldStyleDashboard))
		require.NoError(t, err)

		queries := GetQueriesFromDashboard(json)
		require.Len(t, queries, 1)
		require.Contains(t, queries, int64(2))
		require.Len(t, queries[2], 1)
		query, err := queries[2][0].MarshalJSON()
		require.NoError(t, err)
		require.JSONEq(t, `{
            "datasource": "_yxMP8Ynk",
            "exemplar": true,
            "expr": "go_goroutines{job=\"$job\"}",
            "interval": "",
            "legendFormat": "",
            "refId": "A"
		}`, string(query))
	})
}

func TestGroupQueriesByDataSource(t *testing.T) {
	t.Run("can divide queries by datasource", func(t *testing.T) {
		queries := []*simplejson.Json{
			simplejson.MustJson([]byte(`{
				"datasource": {
					"type": "prometheus",
					"uid": "_yxMP8Ynk"
				},
				"exemplar": true,
				"expr": "go_goroutines{job=\"$job\"}",
				"interval": "",
				"legendFormat": "",
				"refId": "A"
			}`)),
			simplejson.MustJson([]byte(`{
				"datasource": {
					"type": "prometheus",
					"uid": "promds2"
				},
				"exemplar": true,
				"expr": "query2",
				"interval": "",
				"legendFormat": "",
				"refId": "B"
			}`)),
		}

		queriesByDatasource := GroupQueriesByDataSource(queries)
		require.Len(t, queriesByDatasource, 2)
		require.Contains(t, queriesByDatasource, []*simplejson.Json{simplejson.MustJson([]byte(`{
            "datasource": {
              "type": "prometheus",
              "uid": "_yxMP8Ynk"
            },
            "exemplar": true,
            "expr": "go_goroutines{job=\"$job\"}",
            "interval": "",
            "legendFormat": "",
            "refId": "A"
		}`))})
		require.Contains(t, queriesByDatasource, []*simplejson.Json{simplejson.MustJson([]byte(`{
            "datasource": {
              "type": "prometheus",
              "uid": "promds2"
            },
            "exemplar": true,
            "expr": "query2",
            "interval": "",
            "legendFormat": "",
            "refId": "B"
		}`))})
	})
}
