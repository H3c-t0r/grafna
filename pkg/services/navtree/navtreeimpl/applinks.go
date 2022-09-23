package navtreeimpl

import (
	"fmt"
	"path"
	"sort"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/navtree"
	"github.com/grafana/grafana/pkg/services/pluginsettings"
)

func (s *ServiceImpl) addAppLinks(treeRoot *navtree.NavTreeRoot, c *models.ReqContext) error {
	topNavEnabled := s.features.IsEnabled(featuremgmt.FlagTopnav)
	hasAccess := ac.HasAccess(s.accessControl, c)
	appLinks := []*navtree.NavLink{}

	pss, err := s.pluginSettings.GetPluginSettings(c.Req.Context(), &pluginsettings.GetArgs{OrgID: c.OrgID})
	if err != nil {
		return err
	}

	isPluginEnabled := func(plugin plugins.PluginDTO) bool {
		if plugin.AutoEnabled {
			return true
		}
		for _, ps := range pss {
			if ps.PluginID == plugin.ID {
				return ps.Enabled
			}
		}
		return false
	}

	for _, plugin := range s.pluginStore.Plugins(c.Req.Context(), plugins.App) {
		if !isPluginEnabled(plugin) {
			fmt.Printf("not enabled")
			continue
		}

		if !hasAccess(ac.ReqSignedIn,
			ac.EvalPermission(plugins.ActionAppAccess, plugins.ScopeProvider.GetResourceScope(plugin.ID))) {
			fmt.Printf("not permission")
			continue
		}

		if appNode := s.processAppPlugin(plugin, c, topNavEnabled, treeRoot); appNode != nil {
			appLinks = append(appLinks, appNode)
		}
	}

	if len(appLinks) > 0 {
		sort.SliceStable(appLinks, func(i, j int) bool {
			return appLinks[i].Text < appLinks[j].Text
		})
	}

	if topNavEnabled {
		treeRoot.AddSection(&navtree.NavLink{
			Text:        "Apps",
			Icon:        "apps",
			Description: "App plugins",
			Id:          "apps",
			Children:    appLinks,
			Section:     navtree.NavSectionCore,
			Url:         s.cfg.AppSubURL + "/apps",
		})
	} else {
		for _, appLink := range appLinks {
			treeRoot.AddSection(appLink)
		}
	}

	return nil
}

func (s *ServiceImpl) processAppPlugin(plugin plugins.PluginDTO, c *models.ReqContext, topNavEnabled bool, treeRoot *navtree.NavTreeRoot) *navtree.NavLink {
	appLink := &navtree.NavLink{
		Text:       plugin.Name,
		Id:         "plugin-page-" + plugin.ID,
		Img:        plugin.Info.Logos.Small,
		Section:    navtree.NavSectionPlugin,
		SortWeight: navtree.WeightPlugin,
	}

	if s.features.IsEnabled(featuremgmt.FlagTopnav) {
		appLink.Url = s.cfg.AppSubURL + "/a/" + plugin.ID
	} else {
		appLink.Url = path.Join(s.cfg.AppSubURL, plugin.DefaultNavURL)
	}

	for _, include := range plugin.Includes {
		if !c.HasUserRole(include.Role) {
			continue
		}

		if include.Type == "page" && include.AddToNav {
			link := &navtree.NavLink{
				Text: include.Name,
				Icon: include.Icon,
			}

			if len(include.Path) > 0 {
				link.Url = s.cfg.AppSubURL + include.Path
				if include.DefaultNav && !topNavEnabled {
					appLink.Url = link.Url
				}
			} else {
				link.Url = s.cfg.AppSubURL + "/plugins/" + plugin.ID + "/page/" + include.Slug
			}

			if sectionForPageID, ok := s.cfg.NavigationNavIdOverrides[include.Path]; ok {
				if sectionForPage := treeRoot.FindById(sectionForPageID); sectionForPage != nil {
					link.Id = "standalone-plugin-page-" + include.Path
					sectionForPage.Children = append(sectionForPage.Children, link)
				}
			} else {
				appLink.Children = append(appLink.Children, link)
			}
		}

		if include.Type == "dashboard" && include.AddToNav {
			dboardURL := include.DashboardURLPath()
			if dboardURL != "" {
				link := &navtree.NavLink{
					Url:  path.Join(s.cfg.AppSubURL, dboardURL),
					Text: include.Name,
				}
				appLink.Children = append(appLink.Children, link)
			}
		}
	}

	if len(appLink.Children) > 0 {
		// If we only have one child and it's the app default nav then remove it from children
		if len(appLink.Children) == 1 && appLink.Children[0].Url == appLink.Url {
			appLink.Children = []*navtree.NavLink{}
		}

		alertingNode := treeRoot.FindById("alerting")

		if navId, hasNavId := s.cfg.NavigationAppNavIds[plugin.ID]; hasNavId && topNavEnabled {
			if navNode := treeRoot.FindById(navId); navNode != nil {
				navNode.Children = append(navNode.Children, appLink)
			} else {
				if navId == navtree.NavIDMonitoring {
					treeRoot.AddSection(&navtree.NavLink{
						Text:        "Monitoring",
						Id:          navtree.NavIDMonitoring,
						Description: "Monitoring and infrastructure apps",
						Icon:        "heart-rate",
						Section:     navtree.NavSectionCore,
						Children:    []*navtree.NavLink{appLink},
						Url:         s.cfg.AppSubURL + "/monitoring",
					})
				}

				if navId == navtree.NavIDAlertsAndIncidents && alertingNode != nil {
					treeRoot.AddSection(&navtree.NavLink{
						Text:        "Alerts & incidents",
						Id:          navtree.NavIDAlertsAndIncidents,
						Description: "Alerting and incident management apps",
						Icon:        "bell",
						Section:     navtree.NavSectionCore,
						Children:    []*navtree.NavLink{alertingNode, appLink},
						Url:         s.cfg.AppSubURL + "/alerts-and-incidents",
					})
					treeRoot.RemoveSection(alertingNode)
				}
				s.log.Error("Plugin app nav id not found", "pluginId", plugin.ID, "navId", navId)
			}
		} else {
			return appLink
		}
	}

	return nil
}
