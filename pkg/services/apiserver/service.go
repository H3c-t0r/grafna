package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"

	"github.com/grafana/dskit/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/responsewriter"
	genericapiserver "k8s.io/apiserver/pkg/server"
	clientrest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/appcontext"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/modules"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/services/apiserver/auth/authorizer"
	"github.com/grafana/grafana/pkg/services/apiserver/builder"
	grafanaapiserveroptions "github.com/grafana/grafana/pkg/services/apiserver/options"
	entitystorage "github.com/grafana/grafana/pkg/services/apiserver/storage/entity"
	filestorage "github.com/grafana/grafana/pkg/services/apiserver/storage/file"
	"github.com/grafana/grafana/pkg/services/apiserver/utils"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/store/entity"
	"github.com/grafana/grafana/pkg/services/store/entity/db/dbimpl"
	"github.com/grafana/grafana/pkg/services/store/entity/sqlstash"
	"github.com/grafana/grafana/pkg/setting"
)

var (
	_ Service                    = (*service)(nil)
	_ RestConfigProvider         = (*service)(nil)
	_ registry.BackgroundService = (*service)(nil)
	_ registry.CanBeDisabled     = (*service)(nil)

	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)

	unversionedVersion = schema.GroupVersion{Group: "", Version: "v1"}
	unversionedTypes   = []runtime.Object{
		&metav1.Status{},
		&metav1.WatchEvent{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	}
)

func init() {
	// we need to add the options to empty v1
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Group: "", Version: "v1"})
	Scheme.AddUnversionedTypes(unversionedVersion, unversionedTypes...)
}

type Service interface {
	services.NamedService
	registry.BackgroundService
	registry.CanBeDisabled
}

type RestConfigProvider interface {
	GetRestConfig() *clientrest.Config
}

type DirectRestConfigProvider interface {
	// GetDirectRestConfig returns a k8s client configuration that will use the same
	// logged logged in user as the current request context.  This is useful when
	// creating clients that map legacy API handlers to k8s backed services
	GetDirectRestConfig(c *contextmodel.ReqContext) *clientrest.Config

	// This can be used to rewrite incoming requests to path now supported under /apis
	DirectlyServeHTTP(w http.ResponseWriter, r *http.Request)
}

type service struct {
	*services.BasicService

	options    *grafanaapiserveroptions.Options
	restConfig *clientrest.Config

	cfg      *setting.Cfg
	features featuremgmt.FeatureToggles

	stopCh    chan struct{}
	stoppedCh chan error

	db       db.DB
	rr       routing.RouteRegister
	handler  http.Handler
	builders []builder.APIGroupBuilder

	tracing *tracing.TracingService

	authorizer *authorizer.GrafanaAuthorizer
}

func ProvideService(
	cfg *setting.Cfg,
	features featuremgmt.FeatureToggles,
	rr routing.RouteRegister,
	orgService org.Service,
	tracing *tracing.TracingService,
	db db.DB,
) (*service, error) {
	s := &service{
		cfg:        cfg,
		features:   features,
		rr:         rr,
		stopCh:     make(chan struct{}),
		builders:   []builder.APIGroupBuilder{},
		authorizer: authorizer.NewGrafanaAuthorizer(cfg, orgService),
		tracing:    tracing,
		db:         db, // For Unified storage
	}

	// This will be used when running as a dskit service
	s.BasicService = services.NewBasicService(s.start, s.running, nil).WithName(modules.GrafanaAPIServer)

	// TODO: this is very hacky
	// We need to register the routes in ProvideService to make sure
	// the routes are registered before the Grafana HTTP server starts.
	proxyHandler := func(k8sRoute routing.RouteRegister) {
		handler := func(c *contextmodel.ReqContext) {
			if s.handler == nil {
				c.Resp.WriteHeader(404)
				_, _ = c.Resp.Write([]byte("Not found"))
				return
			}

			req := c.Req
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}

			resp := responsewriter.WrapForHTTP1Or2(c.Resp)
			s.handler.ServeHTTP(resp, req)
		}
		k8sRoute.Any("/", middleware.ReqSignedIn, handler)
		k8sRoute.Any("/*", middleware.ReqSignedIn, handler)
	}

	s.rr.Group("/apis", proxyHandler)
	s.rr.Group("/livez", proxyHandler)
	s.rr.Group("/readyz", proxyHandler)
	s.rr.Group("/healthz", proxyHandler)
	s.rr.Group("/openapi", proxyHandler)

	return s, nil
}

func (s *service) GetRestConfig() *clientrest.Config {
	return s.restConfig
}

func (s *service) IsDisabled() bool {
	return false
}

// Run is an adapter for the BackgroundService interface.
func (s *service) Run(ctx context.Context) error {
	if err := s.start(ctx); err != nil {
		return err
	}
	return s.running(ctx)
}

func (s *service) RegisterAPI(b builder.APIGroupBuilder) {
	s.builders = append(s.builders, b)
}

func (s *service) start(ctx context.Context) error {
	// Get the list of groups the server will support
	builders := s.builders

	groupVersions := make([]schema.GroupVersion, 0, len(builders))
	// Install schemas
	for _, b := range builders {
		groupVersions = append(groupVersions, b.GetGroupVersion())
		if err := b.InstallSchema(Scheme); err != nil {
			return err
		}

		auth := b.GetAuthorizer()
		if auth != nil {
			s.authorizer.Register(b.GetGroupVersion(), auth)
		}
	}

	o := grafanaapiserveroptions.NewOptions(Codecs.LegacyCodec(groupVersions...))
	applyGrafanaConfig(s.cfg, s.features, o)

	if errs := o.Validate(); len(errs) != 0 {
		// TODO: handle multiple errors
		return errs[0]
	}

	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)
	if err := o.ApplyTo(serverConfig); err != nil {
		return err
	}
	serverConfig.Authorization.Authorizer = s.authorizer
	serverConfig.TracerProvider = s.tracing.GetTracerProvider()

	// setup loopback transport
	transport := &roundTripperFunc{ready: make(chan struct{})}
	serverConfig.LoopbackClientConfig.Transport = transport
	serverConfig.LoopbackClientConfig.TLSClientConfig = clientrest.TLSClientConfig{}

	switch o.StorageOptions.StorageType {
	case grafanaapiserveroptions.StorageTypeEtcd:
		if err := o.RecommendedOptions.Etcd.Validate(); len(err) > 0 {
			return err[0]
		}
		if err := o.RecommendedOptions.Etcd.ApplyTo(&serverConfig.Config); err != nil {
			return err
		}

	case grafanaapiserveroptions.StorageTypeUnified:
		if !s.features.IsEnabledGlobally(featuremgmt.FlagUnifiedStorage) {
			return fmt.Errorf("unified storage requires the unifiedStorage feature flag (and app_mode = development)")
		}

		eDB, err := dbimpl.ProvideEntityDB(s.db, s.cfg, s.features)
		if err != nil {
			return err
		}

		store, err := sqlstash.ProvideSQLEntityServer(eDB)
		if err != nil {
			return err
		}

		serverConfig.Config.RESTOptionsGetter = entitystorage.NewRESTOptionsGetter(s.cfg, store, o.RecommendedOptions.Etcd.StorageConfig.Codec)

	case grafanaapiserveroptions.StorageTypeUnifiedGrpc:
		// Create a connection to the gRPC server
		// TODO: support configuring the gRPC server address
		conn, err := grpc.Dial("localhost:10000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}

		// TODO: determine when to close the connection, we cannot defer it here
		// defer conn.Close()

		// Create a client instance
		store := entity.NewEntityStoreClientWrapper(conn)

		serverConfig.Config.RESTOptionsGetter = entitystorage.NewRESTOptionsGetter(s.cfg, store, o.RecommendedOptions.Etcd.StorageConfig.Codec)

	case grafanaapiserveroptions.StorageTypeLegacy:
		fallthrough
	case grafanaapiserveroptions.StorageTypeFile:
		serverConfig.RESTOptionsGetter = filestorage.NewRESTOptionsGetter(o.StorageOptions.DataPath, o.RecommendedOptions.Etcd.StorageConfig)
	}

	// Add OpenAPI specs for each group+version
	err := builder.SetupConfig(Scheme, serverConfig, builders)
	if err != nil {
		return err
	}

	// support folder selection
	err = entitystorage.RegisterFieldSelectorSupport(Scheme)
	if err != nil {
		return err
	}

	// Create the server
	server, err := serverConfig.Complete().New("grafana-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return err
	}

	dualWriteEnabled := o.StorageOptions.StorageType != grafanaapiserveroptions.StorageTypeLegacy

	// Install the API Group+version
	err = builder.InstallAPIs(Scheme, Codecs, server, serverConfig.RESTOptionsGetter, builders, dualWriteEnabled)
	if err != nil {
		return err
	}

	// set the transport function and signal that it's ready
	transport.fn = func(req *http.Request) (*http.Response, error) {
		w := newWrappedResponseWriter()
		resp := responsewriter.WrapForHTTP1Or2(w)
		server.Handler.ServeHTTP(resp, req)
		return w.Result(), nil
	}
	close(transport.ready)

	// only write kubeconfig in dev mode
	if o.ExtraOptions.DevMode {
		if err := ensureKubeConfig(server.LoopbackClientConfig, o.StorageOptions.DataPath); err != nil {
			return err
		}
	}

	// Used by the proxy wrapper registered in ProvideService
	s.handler = server.Handler
	s.restConfig = server.LoopbackClientConfig
	s.options = o

	prepared := server.PrepareRun()

	go func() {
		s.stoppedCh <- prepared.Run(s.stopCh)
	}()
	return nil
}

func (s *service) GetDirectRestConfig(c *contextmodel.ReqContext) *clientrest.Config {
	return &clientrest.Config{
		Transport: &roundTripperFunc{
			fn: func(req *http.Request) (*http.Response, error) {
				ctx := appcontext.WithUser(req.Context(), c.SignedInUser)
				w := httptest.NewRecorder()
				s.handler.ServeHTTP(w, req.WithContext(ctx))
				return w.Result(), nil
			},
		},
	}
}

func (s *service) DirectlyServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *service) running(ctx context.Context) error {
	select {
	case err := <-s.stoppedCh:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		close(s.stopCh)
	}
	return nil
}

func ensureKubeConfig(restConfig *clientrest.Config, dir string) error {
	return clientcmd.WriteToFile(
		utils.FormatKubeConfig(restConfig),
		path.Join(dir, "grafana.kubeconfig"),
	)
}

type roundTripperFunc struct {
	ready chan struct{}
	fn    func(req *http.Request) (*http.Response, error)
}

func (f *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	if f.fn == nil {
		<-f.ready
	}
	return f.fn(req)
}

var _ http.ResponseWriter = (*wrappedResponseWriter)(nil)
var _ responsewriter.UserProvidedDecorator = (*wrappedResponseWriter)(nil)

type wrappedResponseWriter struct {
	*httptest.ResponseRecorder
}

func newWrappedResponseWriter() *wrappedResponseWriter {
	w := httptest.NewRecorder()
	return &wrappedResponseWriter{w}
}

func (w *wrappedResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseRecorder
}

func (w *wrappedResponseWriter) CloseNotify() <-chan bool {
	// TODO: this is probably not the right thing to do here
	return make(<-chan bool)
}
