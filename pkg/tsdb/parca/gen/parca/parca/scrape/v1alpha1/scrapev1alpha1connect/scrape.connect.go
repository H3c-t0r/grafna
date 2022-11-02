// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: parca/scrape/v1alpha1/scrape.proto

package scrapev1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/parca-dev/parca/gen/proto/go/parca/scrape/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ScrapeServiceName is the fully-qualified name of the ScrapeService service.
	ScrapeServiceName = "parca.scrape.v1alpha1.ScrapeService"
)

// ScrapeServiceClient is a client for the parca.scrape.v1alpha1.ScrapeService service.
type ScrapeServiceClient interface {
	// Targets returns the set of scrape targets that are configured
	Targets(context.Context, *connect_go.Request[v1alpha1.TargetsRequest]) (*connect_go.Response[v1alpha1.TargetsResponse], error)
}

// NewScrapeServiceClient constructs a client for the parca.scrape.v1alpha1.ScrapeService service.
// By default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped
// responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewScrapeServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ScrapeServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &scrapeServiceClient{
		targets: connect_go.NewClient[v1alpha1.TargetsRequest, v1alpha1.TargetsResponse](
			httpClient,
			baseURL+"/parca.scrape.v1alpha1.ScrapeService/Targets",
			opts...,
		),
	}
}

// scrapeServiceClient implements ScrapeServiceClient.
type scrapeServiceClient struct {
	targets *connect_go.Client[v1alpha1.TargetsRequest, v1alpha1.TargetsResponse]
}

// Targets calls parca.scrape.v1alpha1.ScrapeService.Targets.
func (c *scrapeServiceClient) Targets(ctx context.Context, req *connect_go.Request[v1alpha1.TargetsRequest]) (*connect_go.Response[v1alpha1.TargetsResponse], error) {
	return c.targets.CallUnary(ctx, req)
}

// ScrapeServiceHandler is an implementation of the parca.scrape.v1alpha1.ScrapeService service.
type ScrapeServiceHandler interface {
	// Targets returns the set of scrape targets that are configured
	Targets(context.Context, *connect_go.Request[v1alpha1.TargetsRequest]) (*connect_go.Response[v1alpha1.TargetsResponse], error)
}

// NewScrapeServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewScrapeServiceHandler(svc ScrapeServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/parca.scrape.v1alpha1.ScrapeService/Targets", connect_go.NewUnaryHandler(
		"/parca.scrape.v1alpha1.ScrapeService/Targets",
		svc.Targets,
		opts...,
	))
	return "/parca.scrape.v1alpha1.ScrapeService/", mux
}

// UnimplementedScrapeServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedScrapeServiceHandler struct{}

func (UnimplementedScrapeServiceHandler) Targets(context.Context, *connect_go.Request[v1alpha1.TargetsRequest]) (*connect_go.Response[v1alpha1.TargetsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("parca.scrape.v1alpha1.ScrapeService.Targets is not implemented"))
}
