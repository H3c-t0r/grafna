package k8sbridge

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/grafana/grafana/pkg/schema"
)

var (
	// ErrCRDAlreadyRegistered is returned when trying to register a duplicate CRD.
	ErrCRDAlreadyRegistered = errors.New("error registering duplicate CRD")
)

// Clientset is the clientset for Kubernetes APIs.
// It provides functionality to talk to the APIs as well as register new API clients for CRDs.
type Clientset struct {
	// TODO: this needs to be exposed, but only specific types (e.g. no pods / deployments / etc.).
	*kubernetes.Clientset
	extset *apiextensionsclient.Clientset
	config *rest.Config
	crds   map[k8schema.GroupVersion]apiextensionsv1.CustomResourceDefinition
	lock   sync.RWMutex
}

// NewClientset returns a new Clientset configured with cfg.
func NewClientset(cfg *rest.Config) (*Clientset, error) {
	k8sset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	extset, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Clientset{
		Clientset: k8sset,
		extset:    extset,
		config:    cfg,
		crds:      make(map[k8schema.GroupVersion]apiextensionsv1.CustomResourceDefinition),
	}, nil
}

// RegisterSchema registers a new client and CRD for schema s.
func (c *Clientset) RegisterSchema(ctx context.Context, s schema.CRD) error {
	ver := k8schema.GroupVersion{
		Group:   s.GroupName(),
		Version: s.GroupVersion(),
	}

	c.lock.RLock()
	_, ok := c.crds[ver]
	c.lock.RUnlock()
	if ok {
		return ErrCRDAlreadyRegistered
	}

	crdObj := newCRD(s.Name(), s.GroupName(), s.GroupVersion(), s.OpenAPISchema())
	crd, err := c.extset.
		ApiextensionsV1().
		CustomResourceDefinitions().
		Create(ctx, &crdObj, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	c.lock.Lock()
	c.crds[ver] = *crd
	c.lock.Unlock()

	return nil
}

func newCRD(
	objectKind, groupName, groupVersion string, schema apiextensionsv1.JSONSchemaProps,
) apiextensionsv1.CustomResourceDefinition {
	return apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%ss.%s", objectKind, groupName),
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: groupName,
			Scope: apiextensionsv1.NamespaceScoped, // TODO: make configurable?
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Plural:   objectKind + "s", // TODO: figure out better approach?
				Singular: objectKind,
				Kind:     capitalize(objectKind),
			},
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    groupVersion,
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &schema,
					},
				},
			},
		},
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}

	u := strings.ToUpper(string(s[0]))

	if len(s) == 1 {
		return u
	}

	return u + s[1:]
}
