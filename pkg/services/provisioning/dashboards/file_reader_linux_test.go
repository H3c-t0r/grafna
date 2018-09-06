// +build linux

package dashboards

import (
	"path/filepath"
	"testing"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
)

var (
	symlinkedFolder = "testdata/test-dashboards/symlink"
)

func TestProvsionedSymlinkedFolder(t *testing.T) {
	bus.ClearBusHandlers()
	fakeService = mockDashboardProvisioningService()

	cfg := &DashboardsAsConfig{
		Name:    "Default",
		Type:    "file",
		OrgId:   1,
		Folder:  "",
		Options: map[string]interface{}{"path": symlinkedFolder},
	}

	reader, err := NewDashboardFileReader(cfg, log.New("test-logger"))
	if err != nil {
		t.Error("expected err to be nil")
	}
	err = reader.startWalkingDisk()
	if err != nil {
		t.Errorf("expected err to be nill")
	}

	want, err := filepath.Abs(containingId)

	if err != nil {
		t.Errorf("expected err to be nill")
	}

	if reader.Path != want {
		t.Errorf("got %s want %s", reader.Path, want)
	}
}
