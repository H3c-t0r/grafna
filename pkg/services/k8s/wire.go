package k8s

import (
	"github.com/google/wire"

	"github.com/grafana/grafana/pkg/services/k8s/apiserver"
	"github.com/grafana/grafana/pkg/services/k8s/client"
	"github.com/grafana/grafana/pkg/services/k8s/wrappers"
)

var WireSet = wire.NewSet(
	apiserver.WireSet,
	client.WireSet,
	wrappers.ProvideDashboardStoreWrapper, // Replaces the standard dashboard store
)
