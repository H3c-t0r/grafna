package grafanaapiserver

import (
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/generic"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/spec3"
)

// TODO: this (or something like it) belongs in grafana-app-sdk,
// but lets keep it here while we iterate on a few simple examples
type APIGroupBuilder interface {
	// Add the kinds to the server scheme
	InstallSchema(scheme *runtime.Scheme) error

	// Build the group+version behavior
	GetAPIGroupInfo(
		scheme *runtime.Scheme,
		codecs serializer.CodecFactory, // pointer?
		optsGetter generic.RESTOptionsGetter,
	) (*genericapiserver.APIGroupInfo, error)

	// Get OpenAPI definitions
	GetOpenAPIDefinitions() common.GetOpenAPIDefinitions

	// Get explicit API routes
	GetAPIRoutes() *APIRoutes
}

// This is used to implement dynamic sub-resources like pods/x/logs
type APIRouteHandler struct {
	Path    string           // added to the appropriate level
	Spec    *spec3.PathProps // Exposed in the open api service discovery
	Handler http.HandlerFunc // when Level = resource, the resource will be available in context
}

// APIRoutes define the
type APIRoutes struct {
	// Root handlers are registered directly after the apiVersion identifier
	Root []APIRouteHandler

	// Namespace handlers are mounted under the namespace
	Namespace []APIRouteHandler

	// Resource routes behave the same as pod/logs
	// it looks like a sub-resource, however the response is backed directly by an http handler
	// The current resource can be fetched through context
	Resource map[string]APIRouteHandler
}
