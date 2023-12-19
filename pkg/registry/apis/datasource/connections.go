package datasource

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"

	"github.com/grafana/grafana/pkg/apis"
	"github.com/grafana/grafana/pkg/apis/datasource/v0alpha1"
	"github.com/grafana/grafana/pkg/kinds"
	"github.com/grafana/grafana/pkg/services/datasources"
)

var (
	_ rest.Scoper               = (*connectionAccess)(nil)
	_ rest.SingularNameProvider = (*connectionAccess)(nil)
	_ rest.Getter               = (*connectionAccess)(nil)
	_ rest.Lister               = (*connectionAccess)(nil)
	_ rest.Storage              = (*connectionAccess)(nil)
)

type connectionAccess struct {
	resourceInfo   apis.ResourceInfo
	tableConverter rest.TableConvertor
	builder        *DataSourceAPIBuilder
}

func (s *connectionAccess) New() runtime.Object {
	return s.resourceInfo.NewFunc()
}

func (s *connectionAccess) Destroy() {}

func (s *connectionAccess) NamespaceScoped() bool {
	return true
}

func (s *connectionAccess) GetSingularName() string {
	return s.resourceInfo.GetSingularName()
}

func (s *connectionAccess) ShortNames() []string {
	return s.resourceInfo.GetShortNames()
}

func (s *connectionAccess) NewList() runtime.Object {
	return s.resourceInfo.NewListFunc()
}

func (s *connectionAccess) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return s.tableConverter.ConvertToTable(ctx, object, tableOptions)
}

func (s *connectionAccess) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	ds, err := s.builder.getDataSource(ctx, name)
	if err != nil {
		return nil, err
	}
	return s.asConnection(ds), nil
}

func (s *connectionAccess) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	result := &v0alpha1.DataSourceConnectionList{
		Items: []v0alpha1.DataSourceConnection{},
	}
	vals, err := s.builder.getDataSources(ctx)
	if err == nil {
		for _, ds := range vals {
			result.Items = append(result.Items, *s.asConnection(ds))
		}
	}
	return result, err
}

func (s *connectionAccess) asConnection(ds *datasources.DataSource) *v0alpha1.DataSourceConnection {
	meta := kinds.GrafanaResourceMetadata{
		Name:              ds.UID,
		Namespace:         s.builder.namespacer(ds.OrgID),
		CreationTimestamp: metav1.NewTime(ds.Created),
		ResourceVersion:   fmt.Sprintf("%d", ds.Updated.UnixMilli()),
	}
	meta.SetUpdatedTimestamp(&ds.Updated)
	return &v0alpha1.DataSourceConnection{
		TypeMeta:   s.resourceInfo.TypeMeta(),
		ObjectMeta: metav1.ObjectMeta(meta),
		Title:      ds.Name,
	}
}
