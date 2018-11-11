package dashboards

import (
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/services/dashboards"

	"github.com/grafana/grafana/pkg/models"
)

type DashboardsAsConfig struct {
	Editable              bool
	DisableDeletion       bool
	OrgId                 int64
	UpdateIntervalSeconds int64
	Name                  string
	Type                  string
	Folder                string
	Options               map[string]interface{}
}

type DashboardsAsConfigV0 struct {
	Editable              bool                   `json:"editable" yaml:"editable"`
	DisableDeletion       bool                   `json:"disableDeletion" yaml:"disableDeletion"`
	OrgId                 int64                  `json:"org_id" yaml:"org_id"`
	UpdateIntervalSeconds int64                  `json:"updateIntervalSeconds" yaml:"updateIntervalSeconds"`
	Name                  string                 `json:"name" yaml:"name"`
	Type                  string                 `json:"type" yaml:"type"`
	Folder                string                 `json:"folder" yaml:"folder"`
	Options               map[string]interface{} `json:"options" yaml:"options"`
}

type ConfigVersion struct {
	ApiVersion int64 `json:"apiVersion" yaml:"apiVersion"`
}

type DashboardAsConfigV1 struct {
	Providers []*DashboardProviderConfigs `json:"providers" yaml:"providers"`
}

type DashboardProviderConfigs struct {
	Editable              bool                   `json:"editable" yaml:"editable"`
	DisableDeletion       bool                   `json:"disableDeletion" yaml:"disableDeletion"`
	OrgId                 int64                  `json:"orgId" yaml:"orgId"`
	UpdateIntervalSeconds int64                  `json:"updateIntervalSeconds" yaml:"updateIntervalSeconds"`
	Name                  string                 `json:"name" yaml:"name"`
	Type                  string                 `json:"type" yaml:"type"`
	Folder                string                 `json:"folder" yaml:"folder"`
	Options               map[string]interface{} `json:"options" yaml:"options"`
}

func createDashboardJson(data *simplejson.Json, lastModified time.Time, cfg *DashboardsAsConfig, folderId int64) (*dashboards.SaveDashboardDTO, error) {
	dash := &dashboards.SaveDashboardDTO{}
	dash.Dashboard = models.NewDashboardFromJson(data)
	dash.UpdatedAt = lastModified
	dash.Overwrite = true
	dash.OrgId = cfg.OrgId
	dash.Dashboard.OrgId = cfg.OrgId
	dash.Dashboard.FolderId = folderId

	if dash.Dashboard.Title == "" {
		return nil, models.ErrDashboardTitleEmpty
	}

	return dash, nil
}

func mapV0ToDashboardAsConfig(v0 []*DashboardsAsConfigV0) []*DashboardsAsConfig {
	var r []*DashboardsAsConfig

	for _, v := range v0 {
		r = append(r, &DashboardsAsConfig{
			Name:                  v.Name,
			Type:                  v.Type,
			OrgId:                 v.OrgId,
			Folder:                v.Folder,
			Editable:              v.Editable,
			Options:               v.Options,
			DisableDeletion:       v.DisableDeletion,
			UpdateIntervalSeconds: v.UpdateIntervalSeconds,
		})
	}

	return r
}

func (dc *DashboardAsConfigV1) mapToDashboardAsConfig() []*DashboardsAsConfig {
	var r []*DashboardsAsConfig

	for _, v := range dc.Providers {
		r = append(r, &DashboardsAsConfig{
			Name:                  v.Name,
			Type:                  v.Type,
			OrgId:                 v.OrgId,
			Folder:                v.Folder,
			Editable:              v.Editable,
			Options:               v.Options,
			DisableDeletion:       v.DisableDeletion,
			UpdateIntervalSeconds: v.UpdateIntervalSeconds,
		})
	}

	return r
}
