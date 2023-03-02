package initializer

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/stretchr/testify/assert"

	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/plugins/backendplugin"
	"github.com/grafana/grafana/pkg/plugins/config"
	"github.com/grafana/grafana/pkg/plugins/log"
	"github.com/grafana/grafana/pkg/plugins/manager/fakes"
)

func TestInitializer_Initialize(t *testing.T) {
	absCurPath, err := filepath.Abs(".")
	assert.NoError(t, err)

	t.Run("core backend datasource", func(t *testing.T) {
		p := &plugins.Plugin{
			JSONData: plugins.JSONData{
				ID:   "test",
				Type: plugins.DataSource,
				Includes: []*plugins.Includes{
					{
						Name: "Example dashboard",
						Type: plugins.TypeDashboard,
					},
				},
				Backend: true,
			},
			PluginDir: absCurPath,
			Class:     plugins.Core,
		}

		i := &Initializer{
			cfg: &config.Cfg{},
			log: log.NewTestLogger(),
			backendProvider: &fakeBackendProvider{
				plugin: p,
			},
		}

		err := i.Initialize(context.Background(), p)
		assert.NoError(t, err)

		c, exists := p.Client()
		assert.True(t, exists)
		assert.NotNil(t, c)
	})

	t.Run("renderer", func(t *testing.T) {
		p := &plugins.Plugin{
			JSONData: plugins.JSONData{
				ID:   "test",
				Type: plugins.Renderer,
				Dependencies: plugins.Dependencies{
					GrafanaVersion: ">=8.x",
				},
				Backend: true,
			},
			PluginDir: absCurPath,
			Class:     plugins.External,
		}

		i := &Initializer{
			cfg: &config.Cfg{},
			log: log.NewTestLogger(),
			backendProvider: &fakeBackendProvider{
				plugin: p,
			},
		}

		err := i.Initialize(context.Background(), p)
		assert.NoError(t, err)

		c, exists := p.Client()
		assert.True(t, exists)
		assert.NotNil(t, c)
	})

	t.Run("secretsmanager", func(t *testing.T) {
		p := &plugins.Plugin{
			JSONData: plugins.JSONData{
				ID:   "test",
				Type: plugins.SecretsManager,
				Dependencies: plugins.Dependencies{
					GrafanaVersion: ">=8.x",
				},
				Backend: true,
			},
			PluginDir: absCurPath,
			Class:     plugins.External,
		}

		i := &Initializer{
			cfg: &config.Cfg{},
			log: log.NewTestLogger(),
			backendProvider: &fakeBackendProvider{
				plugin: p,
			},
		}

		err := i.Initialize(context.Background(), p)
		assert.NoError(t, err)

		c, exists := p.Client()
		assert.True(t, exists)
		assert.NotNil(t, c)
	})

	t.Run("non backend plugin app", func(t *testing.T) {
		p := &plugins.Plugin{
			JSONData: plugins.JSONData{
				Backend: false,
			},
		}

		i := &Initializer{
			cfg: &config.Cfg{},
			log: log.NewTestLogger(),
			backendProvider: &fakeBackendProvider{
				plugin: p,
			},
		}

		err := i.Initialize(context.Background(), p)
		assert.NoError(t, err)

		c, exists := p.Client()
		assert.False(t, exists)
		assert.Nil(t, c)
	})
}

func TestInitializer_envVars(t *testing.T) {
	t.Run("backend datasource with license", func(t *testing.T) {
		p := &plugins.Plugin{
			JSONData: plugins.JSONData{
				ID: "test",
			},
		}

		licensing := &fakes.FakeLicensingService{
			LicenseEdition: "test",
			TokenRaw:       "token",
			LicensePath:    "/path/to/ent/license",
		}

		i := &Initializer{
			cfg: &config.Cfg{
				PluginSettings: map[string]map[string]string{
					"test": {
						"custom_env_var": "customVal",
					},
				},
			},
			license: licensing,
			log:     log.NewTestLogger(),
			backendProvider: &fakeBackendProvider{
				plugin: p,
			},
		}

		envVars := i.envVars(p)
		assert.Len(t, envVars, 5)
		assert.Equal(t, "GF_PLUGIN_CUSTOM_ENV_VAR=customVal", envVars[0])
		assert.Equal(t, "GF_VERSION=", envVars[1])
		assert.Equal(t, "GF_EDITION=test", envVars[2])
		assert.Equal(t, "GF_ENTERPRISE_LICENSE_PATH=/path/to/ent/license", envVars[3])
		assert.Equal(t, "GF_ENTERPRISE_LICENSE_TEXT=token", envVars[4])
	})
}

func TestInitializer_tracingEnvironmentVariables(t *testing.T) {
	p := &plugins.Plugin{}

	for _, tc := range []struct {
		name    string
		otelCfg tracing.OpentelemetryCfg
		exp     func(t *testing.T, envVars []string)
	}{
		{
			name:    "disabled",
			otelCfg: tracing.OpentelemetryCfg{},
			exp: func(t *testing.T, envVars []string) {
				assert.Len(t, envVars, 1)
				assert.Equal(t, "GF_VERSION=", envVars[0])
			},
		},
		{
			name: "otlp no propagation",
			otelCfg: tracing.OpentelemetryCfg{
				Address:     "127.0.0.1:4317",
				Propagation: "",
			},
			exp: func(t *testing.T, envVars []string) {
				assert.Len(t, envVars, 3)
				assert.Equal(t, "GF_VERSION=", envVars[0])
				assert.Equal(t, "GF_TRACING_OPENTELEMETRY_OTLP_ADDRESS=127.0.0.1:4317", envVars[1])
				assert.Equal(t, "GF_TRACING_OPENTELEMETRY_OTLP_PROPAGATION=", envVars[2])
			},
		},
		{
			name: "otlp propagation",
			otelCfg: tracing.OpentelemetryCfg{
				Address:     "127.0.0.1:4317",
				Propagation: "w3c",
			},
			exp: func(t *testing.T, envVars []string) {
				assert.Len(t, envVars, 3)
				assert.Equal(t, "GF_VERSION=", envVars[0])
				assert.Equal(t, "GF_TRACING_OPENTELEMETRY_OTLP_ADDRESS=127.0.0.1:4317", envVars[1])
				assert.Equal(t, "GF_TRACING_OPENTELEMETRY_OTLP_PROPAGATION=w3c", envVars[2])
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := &Initializer{
				cfg: &config.Cfg{
					Opentelemetry: tc.otelCfg,
				},
				log: log.NewNopLogger(),
			}
			envVars := i.envVars(p)
			tc.exp(t, envVars)
		})
	}
}

func TestInitializer_getAWSEnvironmentVariables(t *testing.T) {

}

func TestInitializer_handleModuleDefaults(t *testing.T) {

}

func Test_defaultLogoPath(t *testing.T) {

}

func Test_evalRelativePluginUrlPath(t *testing.T) {

}

func Test_getPluginLogoUrl(t *testing.T) {

}

func Test_getPluginSettings(t *testing.T) {

}

func Test_pluginSettings_ToEnv(t *testing.T) {

}

type fakeBackendProvider struct {
	plugins.BackendFactoryProvider

	plugin *plugins.Plugin
}

func (f *fakeBackendProvider) BackendFactory(_ context.Context, _ *plugins.Plugin) backendplugin.PluginFactoryFunc {
	return func(_ string, _ log.Logger, _ []string) (backendplugin.Plugin, error) {
		return f.plugin, nil
	}
}
