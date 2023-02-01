package models

const AlertConfigurationVersion = 1

// AlertConfiguration represents a single version of the Alerting Engine Configuration.
type AlertConfiguration struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	AlertmanagerConfiguration string
	ConfigurationHash         string
	ConfigurationVersion      string
	CreatedAt                 int64 `xorm:"created"`
	Default                   bool
	OrgID                     int64 `xorm:"org_id"`
}

// HistoricAlertConfiguration represents a previously used alerting configuration.
type HistoricAlertConfiguration struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	AlertmanagerConfiguration string
	ConfigurationHash         string
	ConfigurationVersion      string
	CreatedAt                 int64 `xorm:"created"`
	Default                   bool
	OrgID                     int64 `xorm:"org_id"`
	LastApplied               int64 `xorm:"last_applied"`
}

// GetLatestAlertmanagerConfigurationQuery is the query to get the latest alertmanager configuration.
type GetLatestAlertmanagerConfigurationQuery struct {
	OrgID  int64
	Result *AlertConfiguration
}

// SaveAlertmanagerConfigurationCmd is the command to save an alertmanager configuration.
type SaveAlertmanagerConfigurationCmd struct {
	AlertmanagerConfiguration string
	FetchedConfigurationHash  string
	ConfigurationVersion      string
	Default                   bool
	OrgID                     int64
	LastApplied               int64
}

// MarkConfigurationAsAppliedCmd is the command for marking a previously saved configuration as successfully applied.
type MarkConfigurationAsAppliedCmd struct {
	ConfigID          int64
	OrgID             int64
	ConfigurationHash string
}

// GetAppliedConfigurationsQuery is the query for getting configurations that have been previously applied with no errors.
type GetAppliedConfigurationsQuery struct {
	OrgID  int64
	Result []*HistoricAlertConfiguration
}

func HistoricConfigFromAlertConfig(config AlertConfiguration) HistoricAlertConfiguration {
	return HistoricAlertConfiguration{
		AlertmanagerConfiguration: config.AlertmanagerConfiguration,
		ConfigurationHash:         config.ConfigurationHash,
		ConfigurationVersion:      config.ConfigurationVersion,
		Default:                   config.Default,
		OrgID:                     config.OrgID,
	}
}
