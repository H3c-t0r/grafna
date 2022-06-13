package azuremonitor

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-azure-sdk-go/azcredentials"
	"github.com/grafana/grafana-azure-sdk-go/azsettings"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"

	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/tsdb/azuremonitor/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstanceSettings(t *testing.T) {
	tests := []struct {
		name          string
		settings      backend.DataSourceInstanceSettings
		expectedModel types.DatasourceInfo
		Err           require.ErrorAssertionFunc
	}{
		{
			name: "creates an instance",
			settings: backend.DataSourceInstanceSettings{
				JSONData:                []byte(`{"azureAuthType":"msi"}`),
				DecryptedSecureJSONData: map[string]string{"key": "value"},
				ID:                      40,
			},
			expectedModel: types.DatasourceInfo{
				Cloud:                   azsettings.AzurePublic,
				Credentials:             &azcredentials.AzureManagedIdentityCredentials{},
				Settings:                types.AzureMonitorSettings{},
				Routes:                  routes[azsettings.AzurePublic],
				JSONData:                map[string]interface{}{"azureAuthType": "msi"},
				DatasourceID:            40,
				DecryptedSecureJSONData: map[string]string{"key": "value"},
				Services:                map[string]types.DatasourceService{},
			},
			Err: require.NoError,
		},
	}

	cfg := &setting.Cfg{
		Azure: &azsettings.AzureSettings{
			Cloud: azsettings.AzurePublic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewInstanceSettings(cfg, &httpclient.Provider{}, map[string]azDatasourceExecutor{})
			instance, err := factory(tt.settings)
			tt.Err(t, err)
			if !cmp.Equal(instance, tt.expectedModel) {
				t.Errorf("Unexpected instance: %v", cmp.Diff(instance, tt.expectedModel))
			}
		})
	}
}

type fakeInstance struct {
	cloud    string
	routes   map[string]types.AzRoute
	services map[string]types.DatasourceService
	settings types.AzureMonitorSettings
}

func (f *fakeInstance) Get(pluginContext backend.PluginContext) (instancemgmt.Instance, error) {
	return types.DatasourceInfo{
		Cloud:    f.cloud,
		Routes:   f.routes,
		Services: f.services,
		Settings: f.settings,
	}, nil
}

func (f *fakeInstance) Do(pluginContext backend.PluginContext, fn instancemgmt.InstanceCallbackFunc) error {
	return nil
}

type fakeExecutor struct {
	t           *testing.T
	queryType   string
	expectedURL string
}

func (f *fakeExecutor) ResourceRequest(rw http.ResponseWriter, req *http.Request, cli *http.Client) {
}

func (f *fakeExecutor) ExecuteTimeSeriesQuery(ctx context.Context, originalQueries []backend.DataQuery, dsInfo types.DatasourceInfo, client *http.Client,
	url string, tracer tracing.Tracer) (*backend.QueryDataResponse, error) {
	if client == nil {
		f.t.Errorf("The HTTP client for %s is missing", f.queryType)
	} else {
		if url != f.expectedURL {
			f.t.Errorf("Unexpected URL %s wanted %s", url, f.expectedURL)
		}
	}
	return &backend.QueryDataResponse{}, nil
}

func Test_newMux(t *testing.T) {
	tests := []struct {
		name        string
		queryType   string
		expectedURL string
		Err         require.ErrorAssertionFunc
	}{
		{
			name:        "creates an Azure Monitor executor",
			queryType:   azureMonitor,
			expectedURL: routes[azureMonitorPublic][azureMonitor].URL,
			Err:         require.NoError,
		},
		{
			name:        "creates an Azure Log Analytics executor",
			queryType:   azureLogAnalytics,
			expectedURL: routes[azureMonitorPublic][azureLogAnalytics].URL,
			Err:         require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				im: &fakeInstance{
					routes: routes[azureMonitorPublic],
					services: map[string]types.DatasourceService{
						tt.queryType: {
							URL:        routes[azureMonitorPublic][tt.queryType].URL,
							HTTPClient: &http.Client{},
						},
					},
				},
				executors: map[string]azDatasourceExecutor{
					tt.queryType: &fakeExecutor{
						t:           t,
						queryType:   tt.queryType,
						expectedURL: tt.expectedURL,
					},
				},
			}
			mux := s.newQueryMux()
			res, err := mux.QueryData(context.Background(), &backend.QueryDataRequest{
				PluginContext: backend.PluginContext{},
				Queries: []backend.DataQuery{
					{QueryType: tt.queryType},
				},
			})
			tt.Err(t, err)
			// Dummy response from the fake implementation
			if res == nil {
				t.Errorf("Expecting a response")
			}
		})
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestCheckHealth(t *testing.T) {
	okClient := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
			Header:     make(http.Header),
		}
	})
	failClient := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 404,
			Body:       ioutil.NopCloser(bytes.NewBufferString("not found")),
			Header:     make(http.Header),
		}
	})

	cloud := "AzureCloud"
	tests := []struct {
		name                   string
		errorExpected          bool
		expectedResult         *backend.CheckHealthResult
		customServices         map[string]types.DatasourceService
		customInstanceSettings types.AzureMonitorSettings
	}{
		{
			name:          "Successfully queries all endpoints",
			errorExpected: false,
			expectedResult: &backend.CheckHealthResult{
				Status:  backend.HealthStatusOk,
				Message: "1. Successfully connected to Azure Monitor endpoint.\n2. Successfully connected to Azure Log Analytics endpoint.\n3. Successfully connected to Azure Resource Graph endpoint.",
			},
			customServices: map[string]types.DatasourceService{
				azureMonitor: {
					URL:        routes[cloud]["Azure Monitor"].URL,
					HTTPClient: okClient,
				},
				azureLogAnalytics: {
					URL:        routes[cloud]["Azure Log Analytics"].URL,
					HTTPClient: okClient,
				},
				azureResourceGraph: {
					URL:        routes[cloud]["Azure Resource Graph"].URL,
					HTTPClient: okClient,
				}}, customInstanceSettings: types.AzureMonitorSettings{
				LogAnalyticsDefaultWorkspace: "workspace-id",
			},
		},
		{
			name:          "Successfully queries all endpoints except metrics",
			errorExpected: false,
			expectedResult: &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "1. Error connecting to Azure Monitor endpoint.\n2. Successfully connected to Azure Log Analytics endpoint.\n3. Successfully connected to Azure Resource Graph endpoint.",
			},
			customServices: map[string]types.DatasourceService{
				azureMonitor: {
					URL:        routes[cloud]["Azure Monitor"].URL,
					HTTPClient: failClient,
				},
				azureLogAnalytics: {
					URL:        routes[cloud]["Azure Log Analytics"].URL,
					HTTPClient: okClient,
				},
				azureResourceGraph: {
					URL:        routes[cloud]["Azure Resource Graph"].URL,
					HTTPClient: okClient,
				}}, customInstanceSettings: types.AzureMonitorSettings{
				LogAnalyticsDefaultWorkspace: "workspace-id",
			},
		},
		{
			name:          "Successfully queries all endpoints except log analytics",
			errorExpected: false,
			expectedResult: &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "1. Successfully connected to Azure Monitor endpoint.\n2. Error connecting to Azure Log Analytics endpoint.\n3. Successfully connected to Azure Resource Graph endpoint.",
			},
			customServices: map[string]types.DatasourceService{
				azureMonitor: {
					URL:        routes[cloud]["Azure Monitor"].URL,
					HTTPClient: okClient,
				},
				azureLogAnalytics: {
					URL:        routes[cloud]["Azure Log Analytics"].URL,
					HTTPClient: failClient,
				},
				azureResourceGraph: {
					URL:        routes[cloud]["Azure Resource Graph"].URL,
					HTTPClient: okClient,
				}}, customInstanceSettings: types.AzureMonitorSettings{
				LogAnalyticsDefaultWorkspace: "workspace-id",
			},
		},
		{
			name:          "Successfully queries all endpoints except resource graph",
			errorExpected: false,
			expectedResult: &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "1. Successfully connected to Azure Monitor endpoint.\n2. Successfully connected to Azure Log Analytics endpoint.\n3. Error connecting to Azure Resource Graph endpoint.",
			},
			customServices: map[string]types.DatasourceService{
				azureMonitor: {
					URL:        routes[cloud]["Azure Monitor"].URL,
					HTTPClient: okClient,
				},
				azureLogAnalytics: {
					URL:        routes[cloud]["Azure Log Analytics"].URL,
					HTTPClient: okClient,
				},
				azureResourceGraph: {
					URL:        routes[cloud]["Azure Resource Graph"].URL,
					HTTPClient: failClient,
				}},
			customInstanceSettings: types.AzureMonitorSettings{
				LogAnalyticsDefaultWorkspace: "workspace-id",
			},
		},
		{
			name:          "Successfully queries all endpoints even when default workspace ID is undefined",
			errorExpected: false,
			expectedResult: &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: "1. Successfully connected to Azure Monitor endpoint.\n2. Successfully connected to Azure Log Analytics endpoint.\n3. Error connecting to Azure Resource Graph endpoint.",
			},
			customServices: map[string]types.DatasourceService{
				azureMonitor: {
					URL: routes[cloud]["Azure Monitor"].URL,
					HTTPClient: NewTestClient(func(req *http.Request) *http.Response {
						if strings.Contains(req.URL.String(), "workspaces") {
							body := struct {
								Value []types.LogAnalyticsWorkspaceResponse
							}{Value: []types.LogAnalyticsWorkspaceResponse{{
								Id:       "abcd-1234",
								Location: "location",
								Name:     "test-workspace",
								Properties: types.LogAnalyticsWorkspaceProperties{
									CreatedDate: "",
									CustomerId:  "abcd-1234",
									Features:    types.LogAnalyticsWorkspaceFeatures{},
								},
								ProvisioningState:               "provisioned",
								PublicNetworkAccessForIngestion: "enabled",
								PublicNetworkAccessForQuery:     "disabled",
								RetentionInDays:                 0},
							}}
							bodyMarshal, err := json.Marshal(body)
							if err != nil {
								t.Error(err)
							}
							return &http.Response{
								StatusCode: 200,
								Body:       ioutil.NopCloser(bytes.NewBuffer(bodyMarshal)),
								Header:     make(http.Header),
							}
						} else {
							return &http.Response{
								StatusCode: 200,
								Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
								Header:     make(http.Header),
							}
						}
					}),
				},
				azureLogAnalytics: {
					URL:        routes[cloud]["Azure Log Analytics"].URL,
					HTTPClient: okClient,
				},
				azureResourceGraph: {
					URL:        routes[cloud]["Azure Resource Graph"].URL,
					HTTPClient: failClient,
				}},
			customInstanceSettings: types.AzureMonitorSettings{
				LogAnalyticsDefaultWorkspace: "",
			},
		},
	}

	instance := &fakeInstance{
		cloud:    cloud,
		routes:   routes[cloud],
		services: map[string]types.DatasourceService{},
		settings: types.AzureMonitorSettings{
			LogAnalyticsDefaultWorkspace: "workspace-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance.services = tt.customServices
			instance.settings = tt.customInstanceSettings
			s := &Service{
				im: instance,
			}
			res, err := s.CheckHealth(context.Background(), &backend.CheckHealthRequest{
				PluginContext: backend.PluginContext{
					DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{},
				}})
			if tt.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
