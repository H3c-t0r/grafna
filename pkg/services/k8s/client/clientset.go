package client

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/grafana/grafana/pkg/kindsys/k8ssys"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
)

var (
	// ErrCRDAlreadyRegistered is returned when trying to register a duplicate CRD.
	ErrCRDAlreadyRegistered = errors.New("error registering duplicate CRD")
)

type Resource interface {
	dynamic.ResourceInterface
	cache.SharedIndexInformer
}

// Clientset is the clientset for Kubernetes APIs.
// It provides functionality to talk to the APIs as well as register new API clients for CRDs.
type Clientset struct {
	config *rest.Config

	clientset kubernetes.Interface
	extset    apiextensionsclient.Interface
	dynamic   dynamic.Interface
	mapper    meta.RESTMapper
	factory   dynamicinformer.DynamicSharedInformerFactory

	crds map[k8schema.GroupVersion]apiextensionsv1.CustomResourceDefinition
	lock sync.RWMutex
}

func ProvideClientset() (*Clientset, error) {
	cfg, err := GetRESTConfig()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg)
}

// NewFromConfig returns a new Clientset configured with cfg.
func NewFromConfig(cfg *rest.Config) (*Clientset, error) {
	k8sset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	extset, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(k8sset))
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dyn, time.Minute)

	return NewClientset(cfg, k8sset, extset, dyn, mapper, factory)
}

// NewClientset returns a new Clientset.
func NewClientset(
	cfg *rest.Config,
	k8sset kubernetes.Interface,
	extset apiextensionsclient.Interface,
	dyn dynamic.Interface,
	mapper meta.RESTMapper,
	factory dynamicinformer.DynamicSharedInformerFactory,
) (*Clientset, error) {
	return &Clientset{
		config: cfg,

		clientset: k8sset,
		extset:    extset,
		dynamic:   dyn,
		mapper:    mapper,
		factory:   factory,

		crds: make(map[k8schema.GroupVersion]apiextensionsv1.CustomResourceDefinition),
		lock: sync.RWMutex{},
	}, nil
}

// RegisterSchema registers a k8ssys.Kind with the Kubernetes API.
func (c *Clientset) RegisterKind(ctx context.Context, gcrd k8ssys.Kind) error {
	gvk := gcrd.GVK()
	gv := gvk.GroupVersion()

	c.lock.RLock()
	_, ok := c.crds[gv]
	c.lock.RUnlock()
	if ok {
		return ErrCRDAlreadyRegistered
	}

	crd, err := c.extset.
		ApiextensionsV1().
		CustomResourceDefinitions().
		Create(ctx, &gcrd.Schema, metav1.CreateOptions{})

	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	c.lock.Lock()
	c.crds[gv] = *crd
	c.lock.Unlock()

	return nil
}

// GetResourceClient returns a dynamic client for the given Kind and optional namespace.
func (c *Clientset) GetResourceClient(gcrd k8ssys.Kind, namespace ...string) (dynamic.ResourceInterface, error) {
	gvk := gcrd.GVK()
	gk := gvk.GroupKind()

	mapping, err := c.mapper.RESTMapping(gk, gvk.Version)
	if err != nil {
		return nil, err
	}

	var resourceClient dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if len(namespace) == 0 {
			return nil, errors.New("namespace is required for namespaced resources")
		}
		resourceClient = c.dynamic.Resource(mapping.Resource).Namespace(namespace[0])
	} else {
		resourceClient = c.dynamic.Resource(mapping.Resource)
	}

	return resourceClient, nil
}

// GetResourceInformer returns a SharedIndexInformer for the given Kind.
func (c *Clientset) GetResourceInformer(gcrd k8ssys.Kind) cache.SharedIndexInformer {
	gvk := gcrd.GVK()
	gvr := k8schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: gcrd.Schema.Spec.Names.Plural,
	}
	return c.factory.ForResource(gvr).Informer()
}
