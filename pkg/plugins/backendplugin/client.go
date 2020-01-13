package backendplugin

import (
	"os/exec"

	"github.com/grafana/grafana/pkg/infra/log"

	datasourceV1 "github.com/grafana/grafana-plugin-model/go/datasource"
	rendererV1 "github.com/grafana/grafana-plugin-model/go/renderer"
	backend "github.com/grafana/grafana-plugin-sdk-go/backend"
	sdk "github.com/grafana/grafana-plugin-sdk-go/common"
	"github.com/hashicorp/go-plugin"
)

const (
	// DefaultProtocolVersion is the protocol version assumed for legacy clients that don't specify
	// a particular version or version 1 during their handshake. This is currently the version used
	// since Grafana launched support for backend plugins.
	DefaultProtocolVersion = 1
)

// Handshake is the HandshakeConfig used to configure clients and servers.
var handshake = plugin.HandshakeConfig{
	// The ProtocolVersion is the version that must match between Grafana core
	// and Grafana plugins. This should be bumped whenever a (breaking) change
	// happens in one or the other that makes it so that they can't safely communicate.
	ProtocolVersion: DefaultProtocolVersion,

	// The magic cookie values should NEVER be changed.
	MagicCookieKey:   sdk.MagicCookieKey,
	MagicCookieValue: sdk.MagicCookieValue,
}

func newClientConfig(executablePath string, logger log.Logger, versionedPlugins map[int]plugin.PluginSet) *plugin.ClientConfig {
	return &plugin.ClientConfig{
		Cmd:              exec.Command(executablePath),
		HandshakeConfig:  handshake,
		VersionedPlugins: versionedPlugins,
		Logger:           logWrapper{Logger: logger},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}
}

// LegacyStartFunc callback function called when a plugin with old plugin protocol is started.
type LegacyStartFunc func(pluginID string, client *LegacyClient, logger log.Logger) error

// StartFunc callback function called when a plugin with current plugin protocol version is started.
type StartFunc func(pluginID string, client *Client, logger log.Logger) error

// PluginCallbacks callback functions called for plugin throughout its lifetime.
type PluginCallbacks struct {
	OnLegacyStart LegacyStartFunc
	OnStart       StartFunc
}

// PluginDescriptor descriptor used for registering backend plugins.
type PluginDescriptor struct {
	pluginID         string
	executablePath   string
	managed          bool
	versionedPlugins map[int]plugin.PluginSet
	callbacks        PluginCallbacks
}

// NewBackendPluginDescriptor creates a new backend plugin descriptor
// used for registering a backend datasource plugin.
func NewBackendPluginDescriptor(pluginID, executablePath string, callbacks PluginCallbacks) PluginDescriptor {
	return PluginDescriptor{
		pluginID:       pluginID,
		executablePath: executablePath,
		managed:        true,
		versionedPlugins: map[int]plugin.PluginSet{
			DefaultProtocolVersion: {
				pluginID: &datasourceV1.DatasourcePluginImpl{},
			},
			sdk.ProtocolVersion: {
				"backend":   &backend.CoreGRPCPlugin{},
				"transform": &backend.TransformGRPCPlugin{},
			},
		},
		callbacks: callbacks,
	}
}

// NewRendererPluginDescriptor creates a new renderer plugin descriptor
// used for registering a backend renderer plugin.
func NewRendererPluginDescriptor(pluginID, executablePath string, callbacks PluginCallbacks) PluginDescriptor {
	return PluginDescriptor{
		pluginID:       pluginID,
		executablePath: executablePath,
		managed:        false,
		versionedPlugins: map[int]plugin.PluginSet{
			DefaultProtocolVersion: {
				pluginID: &rendererV1.RendererPluginImpl{},
			},
		},
		callbacks: callbacks,
	}
}

// LegacyClient client for communicating with a plugin using the old plugin protocol.
type LegacyClient struct {
	DatasourcePlugin datasourceV1.DatasourcePlugin
	RendererPlugin   rendererV1.RendererPlugin
}

// Client client for communicating with a plugin using the current plugin protocol.
type Client struct {
	BackendPlugin   backend.BackendPlugin
	TransformPlugin backend.TransformPlugin
}
