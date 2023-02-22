// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     CRDTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package publicdashboard

import (
	_ "embed"
	"fmt"

	"github.com/grafana/grafana/pkg/kinds/publicdashboard"
	"github.com/grafana/grafana/pkg/registry/corekind"
	"github.com/grafana/grafana/pkg/services/k8s/crd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var coreReg = corekind.NewBase(nil)

var CRD = crd.Kind{
	GrafanaKind: coreReg.PublicDashboard(),
	Object:      &PublicDashboard{},
	ObjectList:  &PublicDashboardList{},
}

// The CRD YAML representation of the PublicDashboard kind.
//
//go:embed publicdashboard.crd.yml
var CRDYaml []byte

// PublicDashboard is the Go CRD representation of a single PublicDashboard object.
// It implements [runtime.Object], and is used in k8s scheme construction.
type PublicDashboard struct {
	crd.Base[publicdashboard.PublicDashboard]
}

// PublicDashboardList is the Go CRD representation of a list PublicDashboard objects.
// It implements [runtime.Object], and is used in k8s scheme construction.
type PublicDashboardList struct {
	crd.ListBase[publicdashboard.PublicDashboard]
}

// fromUnstructured converts an *unstructured.Unstructured object to a *PublicDashboard.
func fromUnstructured(obj any) (*PublicDashboard, error) {
	uObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed to convert to *unstructured.Unstructured")
	}

	var publicdashboard *PublicDashboard
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), &publicdashboard)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to PublicDashboard: %w", err)
	}

	return publicdashboard, nil
}

// toUnstructured converts a PublicDashboard to an *unstructured.Unstructured.
func toUnstructured(obj *publicdashboard.PublicDashboard, metadata metav1.ObjectMeta) (*unstructured.Unstructured, error) {
	publicdashboardObj := &PublicDashboard{
		crd.Base[publicdashboard.PublicDashboard]{
			TypeMeta: metav1.TypeMeta{
				Kind:       CRD.GVK().Kind,
				APIVersion: CRD.GVK().Group + "/" + CRD.GVK().Version,
			},
			ObjectMeta: metadata,
			Spec:       *obj,
		},
	}

	out, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&publicdashboardObj)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: out,
	}, nil
}
