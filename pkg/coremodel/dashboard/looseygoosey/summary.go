package looseygoosey

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/grafana/grafana/pkg/models"
)

func NewDashboardSummaryBuilder(lookup DatasourceLookup) models.ObjectSummaryBuilder {
	return func(ctx context.Context, uid string, body []byte) (*models.ObjectSummary, []byte, error) {
		summary := &models.ObjectSummary{
			Labels: make(map[string]string),
			Fields: make(map[string]interface{}),
		}
		stream := bytes.NewBuffer(body)
		dash, err := readDashboard(stream, lookup)
		if err != nil {
			summary.Error = &models.ObjectErrorInfo{
				Message: err.Error(),
			}
			return summary, body, err
		}

		dashboardRefs := NewReferenceAccumulator()
		url := fmt.Sprintf("/d/%s/%s", uid, models.SlugifyTitle(dash.Title))
		summary.Name = dash.Title
		summary.Description = dash.Description
		summary.URL = url
		for _, v := range dash.Tags {
			summary.Labels[v] = ""
		}
		if len(dash.TemplateVars) > 0 {
			summary.Fields["hasTemplateVars"] = true
		}
		summary.Fields["schemaVersion"] = dash.SchemaVersion

		for _, panel := range dash.Panels {
			panelRefs := NewReferenceAccumulator()
			p := &models.ObjectSummary{
				UID:  uid + "#" + strconv.FormatInt(panel.ID, 10),
				Kind: "panel",
			}
			p.Name = panel.Title
			p.Description = panel.Description
			p.URL = fmt.Sprintf("%s?viewPanel=%d", url, panel.ID)
			p.Fields = make(map[string]interface{}, 0)

			panelRefs.Add("panel", panel.Type, "")
			for _, v := range panel.Datasource {
				dashboardRefs.Add(StandardKindDataSource, v.Type, v.UID)
				panelRefs.Add(StandardKindDataSource, v.Type, v.UID)
			}

			for _, v := range panel.Transformer {
				panelRefs.Add(StandardKindTransform, v, "")
			}

			dashboardRefs.Add(StandardKindPanel, panel.Type, "")
			p.References = panelRefs.Get()
			summary.Nested = append(summary.Nested, p)
		}

		summary.References = dashboardRefs.Get()
		return summary, body, nil
	}
}
