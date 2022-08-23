package store_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/ngalert/store"
	"github.com/grafana/grafana/pkg/services/ngalert/tests"
)

const baseIntervalSeconds = 10

// Every time this is called, time advances by 1 second.
func mockTimeNow() {
	var timeSeed int64
	store.TimeNow = func() time.Time {
		fakeNow := time.Unix(timeSeed, 0).UTC()
		timeSeed++
		return fakeNow
	}
}

func BenchmarkAlertInstanceOperations(b *testing.B) {
	b.StopTimer()
	ctx := context.Background()
	_, dbstore := tests.SetupTestEnv(b, baseIntervalSeconds)

	const mainOrgID int64 = 1

	alertRule := tests.CreateTestAlertRule(b, ctx, dbstore, 60, mainOrgID)

	// Create some instances to write down and then delete.
	count := 10_000
	instances := make([]models.AlertInstance, 0, count)
	keys := make([]models.AlertInstanceKey, 0, count)
	for i := 0; i < count; i++ {
		labels := models.InstanceLabels{"test": fmt.Sprint(i)}
		_, labelsHash, _ := labels.StringAndHash()
		instance := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule.OrgID,
				RuleUID:    alertRule.UID,
				LabelsHash: labelsHash,
			},
			CurrentState:  models.InstanceStateFiring,
			CurrentReason: string(models.InstanceStateError),
			Labels:        labels,
		}
		instances = append(instances, instance)
		keys = append(keys, instance.AlertInstanceKey)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = dbstore.SaveAlertInstances(ctx, instances...)
		_ = dbstore.DeleteAlertInstances(ctx, keys...)
	}
}

func TestAlertInstanceBulkWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	_, dbstore := tests.SetupTestEnv(t, baseIntervalSeconds)

	const mainOrgID int64 = 1

	alertRule := tests.CreateTestAlertRule(t, ctx, dbstore, 60, mainOrgID)

	// Create some instances to write down and then delete.
	count := 1_000_000
	instances := make([]models.AlertInstance, 0, count)
	keys := make([]models.AlertInstanceKey, 0, count)
	for i := 0; i < count; i++ {
		labels := models.InstanceLabels{"test": fmt.Sprint(i)}
		_, labelsHash, _ := labels.StringAndHash()
		instance := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule.OrgID,
				RuleUID:    alertRule.UID,
				LabelsHash: labelsHash,
			},
			CurrentState:  models.InstanceStateFiring,
			CurrentReason: string(models.InstanceStateError),
			Labels:        labels,
		}
		instances = append(instances, instance)
		keys = append(keys, instance.AlertInstanceKey)
	}

	err := dbstore.SaveAlertInstances(ctx, instances...)
	require.NoError(t, err)

	//// List our instances. Make sure we have 100k.
	q := &models.ListAlertInstancesQuery{
		RuleOrgID: alertRule.OrgID,
	}
	err = dbstore.ListAlertInstances(ctx, q)
	require.NoError(t, err)
	require.Equal(t, count, len(q.Result), "Expected %v instances but got %v", count, len(q.Result))

	err = dbstore.DeleteAlertInstances(ctx, keys...)
	require.NoError(t, err)

	err = dbstore.ListAlertInstances(ctx, q)
	require.NoError(t, err)
	require.Zero(t, len(q.Result), "Deleted instances but still had %v", len(q.Result))
}

func TestIntegrationAlertInstanceOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	_, dbstore := tests.SetupTestEnv(t, baseIntervalSeconds)

	const mainOrgID int64 = 1

	alertRule1 := tests.CreateTestAlertRule(t, ctx, dbstore, 60, mainOrgID)
	orgID := alertRule1.OrgID

	alertRule2 := tests.CreateTestAlertRule(t, ctx, dbstore, 60, mainOrgID)
	require.Equal(t, orgID, alertRule2.OrgID)

	alertRule3 := tests.CreateTestAlertRule(t, ctx, dbstore, 60, mainOrgID)
	require.Equal(t, orgID, alertRule3.OrgID)

	alertRule4 := tests.CreateTestAlertRule(t, ctx, dbstore, 60, mainOrgID)
	require.Equal(t, orgID, alertRule4.OrgID)

	t.Run("can save and read new alert instance", func(t *testing.T) {
		labels := models.InstanceLabels{"test": "testValue"}
		_, hash, _ := labels.StringAndHash()
		instance := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule1.OrgID,
				RuleUID:    alertRule1.UID,
				LabelsHash: hash,
			},
			CurrentState:  models.InstanceStateFiring,
			CurrentReason: string(models.InstanceStateError),
			Labels:        labels,
		}
		err := dbstore.SaveAlertInstances(ctx, instance)
		require.NoError(t, err)

		getCmd := &models.GetAlertInstanceQuery{
			RuleOrgID: instance.RuleOrgID,
			RuleUID:   instance.RuleUID,
			Labels:    models.InstanceLabels{"test": "testValue"},
		}

		err = dbstore.GetAlertInstance(ctx, getCmd)
		require.NoError(t, err)

		require.Equal(t, instance.Labels, getCmd.Result.Labels)
		require.Equal(t, alertRule1.OrgID, getCmd.Result.RuleOrgID)
		require.Equal(t, alertRule1.UID, getCmd.Result.RuleUID)
		require.Equal(t, instance.CurrentReason, getCmd.Result.CurrentReason)
	})

	t.Run("can save and read new alert instance with no labels", func(t *testing.T) {
		labels := models.InstanceLabels{}
		_, hash, _ := labels.StringAndHash()
		instance := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule2.OrgID,
				RuleUID:    alertRule2.UID,
				LabelsHash: hash,
			},
			CurrentState: models.InstanceStateNormal,
			Labels:       labels,
		}
		err := dbstore.SaveAlertInstances(ctx, instance)
		require.NoError(t, err)

		getCmd := &models.GetAlertInstanceQuery{
			RuleOrgID: instance.RuleOrgID,
			RuleUID:   instance.RuleUID,
		}

		err = dbstore.GetAlertInstance(ctx, getCmd)
		require.NoError(t, err)

		require.Equal(t, alertRule2.OrgID, getCmd.Result.RuleOrgID)
		require.Equal(t, alertRule2.UID, getCmd.Result.RuleUID)
		require.Equal(t, instance.Labels, getCmd.Result.Labels)
	})

	t.Run("can save two instances with same org_id, uid and different labels", func(t *testing.T) {
		labels := models.InstanceLabels{"test": "testValue"}
		_, hash, _ := labels.StringAndHash()
		instance1 := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule3.OrgID,
				RuleUID:    alertRule3.UID,
				LabelsHash: hash,
			},
			CurrentState: models.InstanceStateFiring,
			Labels:       labels,
		}

		err := dbstore.SaveAlertInstances(ctx, instance1)
		require.NoError(t, err)

		labels = models.InstanceLabels{"test": "testValue2"}
		_, hash, _ = labels.StringAndHash()
		instance2 := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  instance1.RuleOrgID,
				RuleUID:    instance1.RuleUID,
				LabelsHash: hash,
			},
			CurrentState: models.InstanceStateFiring,
			Labels:       labels,
		}
		err = dbstore.SaveAlertInstances(ctx, instance2)
		require.NoError(t, err)

		listQuery := &models.ListAlertInstancesQuery{
			RuleOrgID: instance1.RuleOrgID,
			RuleUID:   instance1.RuleUID,
		}

		err = dbstore.ListAlertInstances(ctx, listQuery)
		require.NoError(t, err)

		require.Len(t, listQuery.Result, 2)
	})

	t.Run("can list all added instances in org", func(t *testing.T) {
		listQuery := &models.ListAlertInstancesQuery{
			RuleOrgID: orgID,
		}

		err := dbstore.ListAlertInstances(ctx, listQuery)
		require.NoError(t, err)

		require.Len(t, listQuery.Result, 4)
	})

	t.Run("can list all added instances in org filtered by current state", func(t *testing.T) {
		listQuery := &models.ListAlertInstancesQuery{
			RuleOrgID: orgID,
			State:     models.InstanceStateNormal,
		}

		err := dbstore.ListAlertInstances(ctx, listQuery)
		require.NoError(t, err)

		require.Len(t, listQuery.Result, 1)
	})

	t.Run("update instance with same org_id, uid and different state", func(t *testing.T) {
		labels := models.InstanceLabels{"test": "testValue"}
		_, hash, _ := labels.StringAndHash()
		instance1 := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule4.OrgID,
				RuleUID:    alertRule4.UID,
				LabelsHash: hash,
			},
			CurrentState: models.InstanceStateFiring,
			Labels:       labels,
		}

		err := dbstore.SaveAlertInstances(ctx, instance1)
		require.NoError(t, err)

		instance2 := models.AlertInstance{
			AlertInstanceKey: models.AlertInstanceKey{
				RuleOrgID:  alertRule4.OrgID,
				RuleUID:    instance1.RuleUID,
				LabelsHash: instance1.LabelsHash,
			},
			CurrentState: models.InstanceStateNormal,
			Labels:       instance1.Labels,
		}
		err = dbstore.SaveAlertInstances(ctx, instance2)
		require.NoError(t, err)

		listQuery := &models.ListAlertInstancesQuery{
			RuleOrgID: alertRule4.OrgID,
			RuleUID:   alertRule4.UID,
		}

		err = dbstore.ListAlertInstances(ctx, listQuery)
		require.NoError(t, err)

		require.Len(t, listQuery.Result, 1)

		require.Equal(t, instance2.RuleOrgID, listQuery.Result[0].RuleOrgID)
		require.Equal(t, instance2.RuleUID, listQuery.Result[0].RuleUID)
		require.Equal(t, instance2.Labels, listQuery.Result[0].Labels)
		require.Equal(t, instance2.CurrentState, listQuery.Result[0].CurrentState)
	})
}
