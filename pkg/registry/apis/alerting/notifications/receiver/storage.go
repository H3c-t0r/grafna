package receiver

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8srequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"

	"github.com/prometheus/client_golang/prometheus"

	model "github.com/grafana/grafana/pkg/apis/alerting_notifications/v0alpha1"
	grafanaregistry "github.com/grafana/grafana/pkg/apiserver/registry/generic"
	grafanarest "github.com/grafana/grafana/pkg/apiserver/rest"
	"github.com/grafana/grafana/pkg/services/apiserver/endpoints/request"
	"github.com/grafana/grafana/pkg/services/apiserver/utils"
)

var _ grafanarest.Storage = (*storage)(nil)

type storage struct {
	*genericregistry.Store
}

func (s storage) Compare(storageObj, legacyObj runtime.Object) bool {
	// TODO implement when supported dual write mode is not Mode0
	return false
}

func NewStorage(
	legacySvc ReceiverService,
	namespacer request.NamespaceMapper,
	scheme *runtime.Scheme,
	desiredMode grafanarest.DualWriterMode,
	optsGetter generic.RESTOptionsGetter,
	reg prometheus.Registerer) (rest.Storage, error) {
	legacyStore := &legacyStorage{
		service:    legacySvc,
		namespacer: namespacer,
		tableConverter: utils.NewTableConverter(
			resourceInfo.GroupResource(),
			[]metav1.TableColumnDefinition{
				{Name: "Name", Type: "string", Format: "name"},
				{Name: "Title", Type: "string", Format: "string", Description: "The receiver name"}, // TODO: Add integration types.
			},
			func(obj any) ([]interface{}, error) {
				r, ok := obj.(*model.Receiver)
				if ok {
					return []interface{}{
						r.Name,
						r.Spec.Title,
						// r.Spec, //TODO implement formatting for Spec, same as UI?
					}, nil
				}
				return nil, fmt.Errorf("expected resource or info")
			}),
	}
	if optsGetter != nil && desiredMode != grafanarest.Mode0 {
		strategy := grafanaregistry.NewStrategy(scheme)
		s := &genericregistry.Store{
			NewFunc:                   resourceInfo.NewFunc,
			NewListFunc:               resourceInfo.NewListFunc,
			PredicateFunc:             grafanaregistry.Matcher,
			DefaultQualifiedResource:  resourceInfo.GroupResource(),
			SingularQualifiedResource: resourceInfo.SingularGroupResource(),
			TableConvertor:            legacyStore.tableConverter,
			CreateStrategy:            strategy,
			UpdateStrategy:            strategy,
			DeleteStrategy:            strategy,
		}
		options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: grafanaregistry.GetAttrs}
		if err := s.CompleteWithOptions(options); err != nil {
			return nil, err
		}

		requestInfo := &k8srequest.RequestInfo{
			APIGroup:  model.GROUP,
			Resource:  "receivers",
			Name:      "",
			Namespace: namespacer(int64(1)),
		}

		return grafanarest.NewDualWriter(desiredMode, legacyStore, storage{Store: s}, reg, requestInfo, nil), nil
	}
	return legacyStore, nil
}
