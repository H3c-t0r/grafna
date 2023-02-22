// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     CRDTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package librarypanel

import (
	_ "embed"
	"fmt"

	"github.com/grafana/grafana/pkg/kinds/librarypanel"
	"github.com/grafana/grafana/pkg/registry/corekind"
	"github.com/grafana/grafana/pkg/services/k8s/crd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var coreReg = corekind.NewBase(nil)

var CRD = crd.Kind{
	GrafanaKind: coreReg.LibraryPanel(),
	Object:      &LibraryPanel{},
	ObjectList:  &LibraryPanelList{},
}

// The CRD YAML representation of the LibraryPanel kind.
//
//go:embed librarypanel.crd.yml
var CRDYaml []byte

// LibraryPanel is the Go CRD representation of a single LibraryPanel object.
// It implements [runtime.Object], and is used in k8s scheme construction.
type LibraryPanel struct {
	crd.Base[librarypanel.LibraryPanel]
}

// LibraryPanelList is the Go CRD representation of a list LibraryPanel objects.
// It implements [runtime.Object], and is used in k8s scheme construction.
type LibraryPanelList struct {
	crd.ListBase[librarypanel.LibraryPanel]
}

// fromUnstructured converts an *unstructured.Unstructured object to a *LibraryPanel.
func fromUnstructured(obj any) (*LibraryPanel, error) {
	uObj, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("failed to convert to *unstructured.Unstructured")
	}

	var librarypanel *LibraryPanel
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(uObj.UnstructuredContent(), &librarypanel)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to LibraryPanel: %w", err)
	}

	return librarypanel, nil
}

// toUnstructured converts a LibraryPanel to an *unstructured.Unstructured.
func toUnstructured(obj *librarypanel.LibraryPanel, metadata metav1.ObjectMeta) (*unstructured.Unstructured, error) {
	librarypanelObj := crd.Base[librarypanel.LibraryPanel]{
		TypeMeta: metav1.TypeMeta{
			Kind:       CRD.GVK().Kind,
			APIVersion: CRD.GVK().Group + "/" + CRD.GVK().Version,
		},
		ObjectMeta: metadata,
		Spec:       *obj,
	}

	out, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&librarypanelObj)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{
		Object: out,
	}, nil
}
