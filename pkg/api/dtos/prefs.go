package dtos

type Prefs struct {
	Theme           string `json:"theme"`
	HomeDashboardID int64  `json:"homeDashboardId"`
	Timezone        string `json:"timezone"`
	WeekStart       string `json:"weekStart"`
}

type UpdatePrefsCmd struct {
	Theme           string `json:"theme"`
	HomeDashboardID int64  `json:"homeDashboardId"`
	Timezone        string `json:"timezone"`
	WeekStart       string `json:"weekStart"`
}
