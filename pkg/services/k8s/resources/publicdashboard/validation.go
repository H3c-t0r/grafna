package publicdashboard

import (
	"context"
	"fmt"

	"github.com/grafana/grafana/pkg/services/k8s/admission"
	"github.com/grafana/grafana/pkg/services/publicdashboards"
	publicdashboardModels "github.com/grafana/grafana/pkg/services/publicdashboards/models"
	"github.com/grafana/grafana/pkg/services/publicdashboards/validation"
)

var _ admission.ValidatingAdmissionController = (*pdValidation)(nil)

type pdValidation struct {
	publicdashboardsService publicdashboards.Service
	publicdashboardsStore   publicdashboards.Store
}

func ProvideValidation(
	publicdashboardsService publicdashboards.Service,
	publicdashboardsStore publicdashboards.Store,
) *pdValidation {
	return &pdValidation{
		publicdashboardsService: publicdashboardsService,
		publicdashboardsStore:   publicdashboardsStore,
	}
}

func (v *pdValidation) Validate(ctx context.Context, request *admission.AdmissionRequest) error {
	k8Dashboard := request.Object.(*PublicDashboard)
	pdModel, err := k8sObjectToModel(k8Dashboard)
	if err != nil {
		return err
	}

	// API VALIDATIONS
	if !validation.IsValidShortUID(pdModel.DashboardUid) {
		return fmt.Errorf("invalid dashboard ID: %v", pdModel.DashboardUid)
	}

	// validate fields
	// TODO: make validatePublicDashboard take a PublicDashboard Model
	dto := &publicdashboardModels.SavePublicDashboardDTO{
		PublicDashboard: pdModel,
	}
	err = validation.ValidatePublicDashboard(dto)
	if err != nil {
		return err
	}

	// verify public dashboard does not exist and that we didn't get one from the
	// request
	existingPubdash, err := v.publicdashboardsStore.Find(ctx, pdModel.Uid)
	if err != nil {
		return fmt.Errorf("Create: failed to find the public dashboard: %w", err)
	} else if existingPubdash != nil {
		return fmt.Errorf("Create: public dashboard already exists: %s", dto.PublicDashboard.Uid)
	}
	return nil
}
