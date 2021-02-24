package api

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana/pkg/models"
	amv2 "github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/config"
)

// swagger:route POST /api/v1/config alertmanager RoutePostAlertingConfig
//
// sets an Alerting config
//
//     Responses:
//       201: Ack
//       400: ValidationError

// swagger:route GET /api/v1/config alertmanager RouteGetAlertingConfig
//
// gets an Alerting config
//
//     Responses:
//       200: AlertingConfigResponse
//       400: ValidationError

// swagger:route DELETE /api/v1/config alertmanager RouteDeleteAlertingConfig
//
// deletes the Alerting config for a tenant
//
//     Responses:
//       200: Ack
//       400: ValidationError

// swagger:route GET /api/v1/alerts alertmanager RouteGetAmAlerts
//
// get alertmanager alerts
//
//     Responses:
//       200: GettableAlerts
//       400: ValidationError

// swagger:route POST /api/v1/alerts alertmanager RoutePostAmAlerts
//
// create alertmanager alerts
//
//     Responses:
//       200: Ack
//       400: ValidationError

// swagger:route GET /api/v1/alerts/groups alertmanager RouteGetAmAlertGroups
//
// get alertmanager alerts
//
//     Responses:
//       200: AlertGroups
//       400: ValidationError

// swagger:route GET /api/v1/silences alertmanager RouteGetSilences
//
// get silences
//
//     Responses:
//       200: GettableSilences
//       400: ValidationError

// swagger:route POST /api/v1/silences alertmanager RouteCreateSilence
//
// create silence
//
//     Responses:
//       201: GettableSilence
//       400: ValidationError

// swagger:route GET /api/v1/silence/{SilenceId} alertmanager RouteGetSilence
//
// get silence
//
//     Responses:
//       200: GettableSilence
//       400: ValidationError

// swagger:route DELETE /api/v1/silence/{SilenceId} alertmanager RouteDeleteSilence
//
// delete silence
//
//     Responses:
//       200: Ack
//       400: ValidationError

// swagger:parameters RouteCreateSilence
type CreateSilenceParams struct {
	// in:body
	Body SilenceBody
}

//swagger:parameters RouteGetSilence RouteDeleteSilence
type GetDeleteSilenceParams struct {
	// in:path
	SilenceId string
}

// swagger:model
type SilenceBody struct {
	Idd string `json:"id"`
	amv2.Silence
}

// swagger:model
type GettableSilences []amv2.GettableSilences

// swagger:model
type GetableSilence amv2.Silence

// swagger:model
type GettableAlerts amv2.GettableAlerts

// swagger:model
type AlertGroups amv2.AlertGroups

// swagger:parameters RouteGetAmAlerts RouteGetAmAlertGroups
type AlertsParams struct {

	// Show active alerts
	// in: query
	// required: false
	Active bool `json:"active"`

	// Show silenced alerts
	// in: query
	// required: false
	Silenced bool `json:"silenced"`

	// Show inhibited alerts
	// in: query
	// required: false
	Inhibited bool `json:"inhibited"`

	// A list of matchers to filter alerts by
	// in: query
	// required: false
	Matchers []string `json:"matchers"`

	// A list of receivers to filter alerts by
	// in: query
	// required: false
	Receivers []string `json:"receivers"`
}

// swagger:parameters RoutePostAmAlerts
type PostableAlerts struct {
	// in:body
	PostableAlerts []amv2.PostableAlert `yaml:"" json:""`
}

// swagger:parameters RoutePostAlertingConfig
type BodyAlertingConfig struct {
	// in:body
	Body UserConfig
}

// swagger:model
type UserConfig struct {
	TemplateFiles      map[string]string `yaml:"template_files" json:"template_files"`
	AlertmanagerConfig ApiAlertingConfig `yaml:"alertmanager_config" json:"alertmanager_config"`
}

// swagger:model
type AlertingConfigResponse struct {
	BodyAlertingConfig
}

type ApiAlertingConfig struct {
	config.Config

	// Override with our superset receiver type
	Receivers []*ApiReceiver `yaml:"receivers,omitempty" json:"receivers,omitempty"`
}

type GrafanaReceiver models.CreateAlertNotificationCommand

type ReceiverType int

const (
	GrafanaManagedReceiver ReceiverType = iota
	AMReceiver
)

type ApiReceiver struct {
	config.Receiver
	GrafanaReceivers
}

func (r *ApiReceiver) UnmarshalJSON(b []byte) error {
	type plain ApiReceiver
	if err := json.Unmarshal(b, (*plain)(r)); err != nil {
		return err
	}

	hasGrafanaReceivers := len(r.GrafanaReceivers.GrafanaManagedReceivers) > 0

	if hasGrafanaReceivers {
		if len(r.EmailConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager EmailConfigs & Grafana receivers together")
		}
		if len(r.PagerdutyConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager PagerdutyConfigs & Grafana receivers together")
		}
		if len(r.SlackConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager SlackConfigs & Grafana receivers together")
		}
		if len(r.WebhookConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager WebhookConfigs & Grafana receivers together")
		}
		if len(r.OpsGenieConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager OpsGenieConfigs & Grafana receivers together")
		}
		if len(r.WechatConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager WechatConfigs & Grafana receivers together")
		}
		if len(r.PushoverConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager PushoverConfigs & Grafana receivers together")
		}
		if len(r.VictorOpsConfigs) > 0 {
			return fmt.Errorf("cannot have both Alertmanager VictorOpsConfigs & Grafana receivers together")
		}

	}

	return nil

}

type GrafanaReceivers struct {
	GrafanaManagedReceivers []*GrafanaReceiver `yaml:"grafana_managed_receiver_configs,omitempty" json:"grafana_managed_receiver_configs,omitempty"`
}
