// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: parca/profilestore/v1alpha1/profilestore.proto

package profilestorev1alpha1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1alpha1 "github.com/parca-dev/parca/gen/proto/go/parca/profilestore/v1alpha1"
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
	// ProfileStoreServiceName is the fully-qualified name of the ProfileStoreService service.
	ProfileStoreServiceName = "parca.profilestore.v1alpha1.ProfileStoreService"
)

// ProfileStoreServiceClient is a client for the parca.profilestore.v1alpha1.ProfileStoreService
// service.
type ProfileStoreServiceClient interface {
	// WriteRaw accepts a raw set of bytes of a pprof file
	WriteRaw(context.Context, *connect_go.Request[v1alpha1.WriteRawRequest]) (*connect_go.Response[v1alpha1.WriteRawResponse], error)
}

// NewProfileStoreServiceClient constructs a client for the
// parca.profilestore.v1alpha1.ProfileStoreService service. By default, it uses the Connect protocol
// with the binary Protobuf Codec, asks for gzipped responses, and sends uncompressed requests. To
// use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or connect.WithGRPCWeb()
// options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProfileStoreServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ProfileStoreServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &profileStoreServiceClient{
		writeRaw: connect_go.NewClient[v1alpha1.WriteRawRequest, v1alpha1.WriteRawResponse](
			httpClient,
			baseURL+"/parca.profilestore.v1alpha1.ProfileStoreService/WriteRaw",
			opts...,
		),
	}
}

// profileStoreServiceClient implements ProfileStoreServiceClient.
type profileStoreServiceClient struct {
	writeRaw *connect_go.Client[v1alpha1.WriteRawRequest, v1alpha1.WriteRawResponse]
}

// WriteRaw calls parca.profilestore.v1alpha1.ProfileStoreService.WriteRaw.
func (c *profileStoreServiceClient) WriteRaw(ctx context.Context, req *connect_go.Request[v1alpha1.WriteRawRequest]) (*connect_go.Response[v1alpha1.WriteRawResponse], error) {
	return c.writeRaw.CallUnary(ctx, req)
}

// ProfileStoreServiceHandler is an implementation of the
// parca.profilestore.v1alpha1.ProfileStoreService service.
type ProfileStoreServiceHandler interface {
	// WriteRaw accepts a raw set of bytes of a pprof file
	WriteRaw(context.Context, *connect_go.Request[v1alpha1.WriteRawRequest]) (*connect_go.Response[v1alpha1.WriteRawResponse], error)
}

// NewProfileStoreServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProfileStoreServiceHandler(svc ProfileStoreServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/parca.profilestore.v1alpha1.ProfileStoreService/WriteRaw", connect_go.NewUnaryHandler(
		"/parca.profilestore.v1alpha1.ProfileStoreService/WriteRaw",
		svc.WriteRaw,
		opts...,
	))
	return "/parca.profilestore.v1alpha1.ProfileStoreService/", mux
}

// UnimplementedProfileStoreServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProfileStoreServiceHandler struct{}

func (UnimplementedProfileStoreServiceHandler) WriteRaw(context.Context, *connect_go.Request[v1alpha1.WriteRawRequest]) (*connect_go.Response[v1alpha1.WriteRawResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("parca.profilestore.v1alpha1.ProfileStoreService.WriteRaw is not implemented"))
}
