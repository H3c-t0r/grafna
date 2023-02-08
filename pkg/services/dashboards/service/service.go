package service

import (
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/store/entity"
	"github.com/grafana/grafana/pkg/services/store/k8saccess"
)

func ProvideSimpleDashboardService(
	features featuremgmt.FeatureToggles,
	svc *DashboardServiceImpl,
	store entity.EntityStoreServer,
) dashboards.DashboardService {
	if features.IsEnabled(featuremgmt.FlagK8sDashboards) {
		return k8saccess.NewDashboardService(svc, store)
	}
	return svc
}

func ProvideDashboardProvisioningService(
	features featuremgmt.FeatureToggles, orig *DashboardServiceImpl,
) dashboards.DashboardProvisioningService {
	return orig
}

func ProvideDashboardPluginService(
	features featuremgmt.FeatureToggles, orig *DashboardServiceImpl,
) dashboards.PluginService {
	return orig
}
