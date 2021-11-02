package manager

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/grafana/grafana/pkg/plugins"
)

func (m *PluginManager) Plugin(pluginID string) (plugins.PluginDTO, bool) {
	m.pluginsMu.RLock()
	p, ok := m.plugins[pluginID]
	m.pluginsMu.RUnlock()

	if ok && (p.IsDecommissioned()) {
		return plugins.PluginDTO{}, false
	}

	return p.ToDTO(), true
}

func (m *PluginManager) Plugins(pluginTypes ...plugins.Type) []plugins.PluginDTO {
	// if no types passed, assume all
	if len(pluginTypes) == 0 {
		pluginTypes = plugins.PluginTypes
	}

	var requestedTypes = make(map[plugins.Type]struct{})
	for _, pt := range pluginTypes {
		requestedTypes[pt] = struct{}{}
	}

	m.pluginsMu.RLock()
	pluginsList := []plugins.PluginDTO{}
	for _, p := range m.plugins {
		if _, exists := requestedTypes[p.Type]; exists {
			pluginsList = append(pluginsList, p.ToDTO())
		}
	}
	m.pluginsMu.RUnlock()
	return pluginsList
}

func (m *PluginManager) Add(ctx context.Context, pluginID, version string, opts plugins.AddOpts) error {
	var pluginZipURL string

	if opts.PluginRepoURL == "" {
		opts.PluginRepoURL = grafanaComURL
	}

	if plugin, exists := m.plugins[pluginID]; exists {
		if !plugin.IsExternalPlugin() {
			return plugins.ErrInstallCorePlugin
		}

		if plugin.Info.Version == version {
			return plugins.DuplicateError{
				PluginID:          plugin.ID,
				ExistingPluginDir: plugin.PluginDir,
			}
		}

		// get plugin update information to confirm if upgrading is possible
		updateInfo, err := m.pluginInstaller.GetUpdateInfo(ctx, pluginID, version, opts.PluginRepoURL)
		if err != nil {
			return err
		}

		pluginZipURL = updateInfo.PluginZipURL

		// remove existing installation of plugin
		err = m.Remove(ctx, plugin.ID)
		if err != nil {
			return err
		}
	}

	if opts.PluginInstallDir == "" {
		opts.PluginInstallDir = m.cfg.PluginsPath
	}

	if opts.PluginZipURL == "" {
		opts.PluginZipURL = pluginZipURL
	}

	err := m.pluginInstaller.Install(ctx, pluginID, version, opts.PluginInstallDir, opts.PluginZipURL, opts.PluginRepoURL)
	if err != nil {
		return err
	}

	err = m.loadPlugins(opts.PluginInstallDir)
	if err != nil {
		return err
	}

	return nil
}

func (m *PluginManager) Remove(ctx context.Context, pluginID string) error {
	plugin, exists := m.plugins[pluginID]
	if !exists {
		return plugins.ErrPluginNotInstalled
	}

	if !plugin.IsExternalPlugin() {
		return plugins.ErrUninstallCorePlugin
	}

	// extra security check to ensure we only remove plugins that are located in the configured plugins directory
	path, err := filepath.Rel(m.cfg.PluginsPath, plugin.PluginDir)
	if err != nil || strings.HasPrefix(path, ".."+string(filepath.Separator)) {
		return plugins.ErrUninstallOutsideOfPluginDir
	}

	if m.isRegistered(pluginID) {
		err := m.unregisterAndStop(ctx, plugin)
		if err != nil {
			return err
		}
	}

	return m.pluginInstaller.Uninstall(ctx, plugin.PluginDir)
}
