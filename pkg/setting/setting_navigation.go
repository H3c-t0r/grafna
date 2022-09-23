package setting

type NavigationAppConfig struct {
	SectionID  string
	SortWeight int64
}

func (cfg *Cfg) readNavigationSettings() error {
	//sec := cfg.Raw.Section("navigation")
	cfg.NavigationAppConfig = map[string]NavigationAppConfig{
		"grafana-k8s-app":                  {SectionID: "monitoring", SortWeight: 1},
		"grafana-synthetic-monitoring-app": {SectionID: "monitoring", SortWeight: 2},
		"grafana-ml-app":                   {SectionID: "alerts-and-incidents"},
		"grafana-incident-app":             {SectionID: "alerts-and-incidents"},
	}

	cfg.NavigationAppPathConfig = map[string]NavigationAppConfig{
		"/a/myorgid-simple-app/catalog": {SectionID: "cfg"},
	}

	// for _, key := range sec.Keys() {
	// 	if strings.HasPrefix(key.Name(), "nav_id_") {
	// 		pluginId := strings.Replace(key.Name(), "nav_id_", "", 1)
	// 		cfg.NavigationAppConfig[pluginId] = sec.Key(key.Name()).MustString("")
	// 	}
	// }

	return nil
}
