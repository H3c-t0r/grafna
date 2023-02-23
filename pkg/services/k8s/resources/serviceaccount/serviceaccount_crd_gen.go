// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     CRDTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package serviceaccount

import (
	_ "embed"
	"fmt"

	"github.com/grafana/grafana/pkg/kinds/serviceaccount"
	"github.com/grafana/grafana/pkg/registry/corekind"
	"github.com/grafana/grafana/pkg/services/k8s/crd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var coreReg = corekind.NewBase(nil)

var CRD = crd.Kind{
	GrafanaKind: coreReg.ServiceAccount(),
	Object:      &ServiceAccount{},
	ObjectList:  &ServiceAccountList{},
}

// The CRD YAML representation of the ServiceAccount kind.
//
//go:embed serviceaccount.crd.yml
var CRDYaml []byte

// ServiceAccount is the Go CRD representation of a single ServiceAccount object.
// It implements [runtime.Object], and is used in k8s scheme construction.
type ServiceAccount struct {
	crd.Base[serviceaccount.ServiceAccount]
}

// ServiceAccountList is the Go CRD representation of a list ServiceAccount objects.
// It implements [runtime.Object], and is used in k8s scheme construction.
type ServiceAccountList struct {
	crd.ListBase[serviceaccount.ServiceAccount]
}

// fromUnstructured converts an *unstructured.Unstructured object to a *ServiceAccount.
func fromUnstructured(obj any) (*ServiceAccount, error) {
	uObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed to convert to *unstructured.Unstructured")
	}

	var serviceaccount crd.Base[serviceaccount.ServiceAccount]
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), &serviceaccount)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to ServiceAccount: %w", err)
	}

	return &ServiceAccount{serviceaccount}, nil
}

// toUnstructured converts a ServiceAccount to an *unstructured.Unstructured.
func toUnstructured(obj *serviceaccount.ServiceAccount, metadata metav1.ObjectMeta) (*unstructured.Unstructured, error) {
	serviceaccountObj := crd.Base[serviceaccount.ServiceAccount]{
		TypeMeta: metav1.TypeMeta{
			Kind:       CRD.GVK().Kind,
			APIVersion: CRD.GVK().Group + "/" + CRD.GVK().Version,
		},
		ObjectMeta: metadata,
		Spec:       *obj,
	}

	out, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&serviceaccountObj)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: out,
	}, nil
}
