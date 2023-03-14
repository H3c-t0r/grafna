package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/dashboards/database"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/k8s/client"
	"github.com/grafana/grafana/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE this is how you reset the CRD
//kubectl delete CustomResourceDefinition dashboards.dashboard.core.grafana.com

type StoreWrapper struct {
	database.DashboardSQLStore
	log       log.Logger
	clientset client.ClientSetProvider
	namespace string
}

var _ dashboards.Store = (*StoreWrapper)(nil)

func ProvideStoreWrapper(
	features featuremgmt.FeatureToggles,
	store database.DashboardSQLStore,
	clientset client.ClientSetProvider,
) (dashboards.Store, error) {
	// When feature is disabled, resolve the upstream SQL store
	if !features.IsEnabled(featuremgmt.FlagK8s) {
		return store, nil
	}
	return &StoreWrapper{
		DashboardSQLStore: store,
		log:               log.New("k8s.dashboards.service-wrapper"),
		clientset:         clientset,
		namespace:         "default",
	}, nil
}

// SaveDashboard saves the dashboard to kubernetes
func (s *StoreWrapper) SaveDashboard(ctx context.Context, cmd dashboards.SaveDashboardCommand) (*dashboards.Dashboard, error) {
	// Same save path but with additional metadata
	return s.SaveProvisionedDashboard(ctx, cmd, nil)
}

// SaveDashboard will write the dashboard to k8s then wait for it to exist in the SQL store
func (s *StoreWrapper) SaveProvisionedDashboard(ctx context.Context, cmd dashboards.SaveDashboardCommand, provisioning *dashboards.DashboardProvisioning) (*dashboards.Dashboard, error) {
	// TODO: dashboards and folders are managed together now... should be its own resource
	if cmd.IsFolder {
		if provisioning != nil {
			return s.DashboardSQLStore.SaveProvisionedDashboard(ctx, cmd, provisioning)
		}
		return s.DashboardSQLStore.SaveDashboard(ctx, cmd)
	}

	dashboardResource, err := s.clientset.GetClientset().GetResourceClient(CRD)
	if err != nil {
		return nil, fmt.Errorf("ProvideServiceWrapper failed to get dashboard resource client: %w", err)
	}

	if cmd.Dashboard == nil {
		return nil, fmt.Errorf("dashboard data is nil")
	}

	anno := entityAnnotations{
		OrgID:     cmd.OrgID,
		Message:   cmd.Message,
		FolderID:  cmd.FolderID,
		FolderUID: cmd.FolderUID,
		PluginID:  cmd.PluginID,
		UpdatedAt: cmd.UpdatedAt.UnixMilli(),
		UpdatedBy: cmd.UserID,
	}

	// Save provisioning info
	if provisioning != nil {
		anno.OriginName = provisioning.Name
		anno.OriginPath = provisioning.ExternalID
		anno.OriginKey = provisioning.CheckSum
		anno.OriginTime = provisioning.Updated
	}

	meta := metav1.ObjectMeta{
		Namespace: s.namespace,
	}

	// FIXME this is not reliable and is spaghetti
	dto := cmd.GetDashboardModel()
	uid := dto.UID
	if uid == "" {
		uid = util.GenerateShortUID()
		meta.Name = GrafanaUIDToK8sName(uid)
	} else {
		// Get the previous version
		meta.Name = GrafanaUIDToK8sName(uid)
		r, err := dashboardResource.Get(ctx, meta.Name, metav1.GetOptions{})
		if err != nil || r == nil {
			//return nil, fmt.Errorf("unable to find k8s dashboard: " + uid)
			fmt.Printf("UNABLE TO FIND: " + uid)
		} else {
			if !cmd.Overwrite {
				fmt.Printf("TODO... verify SQL version: %s\n", r.GetResourceVersion())
			}

			// Keep old metadata
			meta.ResourceVersion = r.GetResourceVersion()
			anno.Merge(r.GetAnnotations())
			if anno.CreatedAt < 100 {
				anno.CreatedAt = r.GetCreationTimestamp().UnixMilli()
			}
		}
	}

	// HACK, remove empty ID!!
	// dto.Data.Del("id") <<<<< MUST KEEP UID since this is the real key
	dto.Data.Set("uid", uid)
	dto.UID = uid
	// strip nulls...
	stripNulls(dto.Data)

	dashbytes, err := dto.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}

	d, _, err := coreReg.Dashboard().JSONValueMux(dashbytes)
	if err != nil {
		fmt.Printf("-------- FAILED TO PARSE ---------")
		fmt.Printf("%s", string(dashbytes))
		return nil, fmt.Errorf("dashboard JSONValueMux failed: %w", err)
	}

	if d.Uid == nil {
		d.Uid = &uid
	}

	if d.Title == nil {
		d.Title = &dto.Title
	}

	if anno.CreatedAt < 100 {
		anno.CreatedAt = time.Now().UnixMilli()
	}
	if anno.CreatedBy < 1 {
		anno.CreatedBy = anno.UpdatedBy
	}
	meta.Annotations = anno.ToMap()
	uObj, err := toUnstructured(d, meta)
	if err != nil {
		return nil, err
	}

	js, _ := json.MarshalIndent(uObj, "", "  ")
	fmt.Printf("-------- WRAPPER BEFORE SAVE ---------")
	fmt.Printf("%s", string(js))

	if meta.ResourceVersion == "" {
		s.log.Debug("k8s action: create")
		uObj, err = dashboardResource.Create(ctx, uObj, metav1.CreateOptions{})
	} else {
		s.log.Debug("k8s action: update")
		uObj, err = dashboardResource.Update(ctx, uObj, metav1.UpdateOptions{})
	}

	// create or update error
	if err != nil {
		return nil, err
	}

	js, _ = json.MarshalIndent(uObj, "", "  ")
	fmt.Printf("-------- WRAPPER AFTER SAVE ---------")
	fmt.Printf("%s", string(js))

	rv := uObj.GetResourceVersion()
	s.log.Debug("wait for revision", "revision", rv)

	// TODO: rather than polling the dashboard service,
	// we could write a status field and listen for changes on that status from k8s directly
	// however, this is likely better since it is checking the SQL instance that needs to be valid
	for i := 0; i < 9; i++ {
		time.Sleep(175 * time.Millisecond)
		out, err := s.DashboardSQLStore.GetDashboard(ctx, &dashboards.GetDashboardQuery{UID: uid, OrgID: dto.OrgID})
		if err != nil {
			fmt.Printf("ERROR: %v", err)
			continue
		}
		if out != nil && out.Data != nil {
			savedRV := out.Data.Get("resourceVersion").MustString()
			if savedRV == rv {
				return out, nil
			} else {
				fmt.Printf("NO MATCH: %v\n", out)
			}
		}
	}

	// too many loops?
	return nil, fmt.Errorf("controller never ran?")
}
