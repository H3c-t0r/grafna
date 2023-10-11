package v0alpha1

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/rest"

	playlistkind "github.com/grafana/grafana/pkg/kinds/playlist"
	grafanarequest "github.com/grafana/grafana/pkg/services/grafana-apiserver/endpoints/request"
	"github.com/grafana/grafana/pkg/services/playlist"
)

var (
	_ rest.Scoper               = (*legacyStorage)(nil)
	_ rest.SingularNameProvider = (*legacyStorage)(nil)
	_ rest.Getter               = (*legacyStorage)(nil)
	_ rest.Lister               = (*legacyStorage)(nil)
	_ rest.Storage              = (*legacyStorage)(nil)
)

type legacyStorage struct {
	service playlist.Service
}

func newLegacyStorage(s playlist.Service) *legacyStorage {
	return &legacyStorage{
		service: s,
	}
}

func (s *legacyStorage) New() runtime.Object {
	return &Playlist{}
}

func (s *legacyStorage) Destroy() {}

func (s *legacyStorage) NamespaceScoped() bool {
	return true // namespace == org
}

func (s *legacyStorage) GetSingularName() string {
	return "playlist"
}

func (s *legacyStorage) NewList() runtime.Object {
	return &PlaylistList{}
}

func (s *legacyStorage) ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rest.NewDefaultTableConvertor(Resource("playlists")).ConvertToTable(ctx, object, tableOptions)
}

func (s *legacyStorage) List(ctx context.Context, options *internalversion.ListOptions) (runtime.Object, error) {
	// TODO: handle fetching all available orgs when no namespace is specified
	// To test: kubectl get playlists --all-namespaces
	orgId, ok := grafanarequest.OrgIDFrom(ctx)
	if !ok {
		orgId = 1 // TODO: default org ID 1 for now
	}

	limit := 100
	if options.Limit > 0 {
		limit = int(options.Limit)
	}
	res, err := s.service.Search(ctx, &playlist.GetPlaylistsQuery{
		OrgId: orgId,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	list := &PlaylistList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PlaylistList",
			APIVersion: APIVersion,
		},
	}
	for _, v := range res {
		p := Playlist{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Playlist",
				APIVersion: APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:              v.UID,
				CreationTimestamp: metav1.NewTime(time.UnixMilli(v.CreatedAt)),
				ResourceVersion:   fmt.Sprintf("%d", v.UpdatedAt),
				UID:               types.UID(v.UID),
			},
			Spec: playlistkind.Spec{
				Title:    v.Name,
				Interval: v.Interval,
			},
		}
		list.Items = append(list.Items, p)
	}
	if len(list.Items) == limit {
		list.Continue = "<more>" // TODO?
	}
	return list, nil
}

func (s *legacyStorage) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	orgId, ok := grafanarequest.OrgIDFrom(ctx)
	if !ok {
		orgId = 1 // TODO: default org ID 1 for now
	}

	dto, err := s.service.Get(ctx, &playlist.GetPlaylistByUidQuery{
		UID:   name,
		OrgId: orgId,
	})
	if err != nil {
		return nil, err
	}
	if dto == nil {
		return nil, fmt.Errorf("not found?")
	}

	p := &Playlist{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Playlist",
			APIVersion: APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              dto.Uid,
			CreationTimestamp: metav1.NewTime(time.UnixMilli(dto.CreatedAt)),
			ResourceVersion:   fmt.Sprintf("%d", dto.UpdatedAt),
			UID:               types.UID(dto.Uid),
		},
		Spec: playlistkind.Spec{
			Title:    dto.Name,
			Interval: dto.Interval,
		},
	}
	for _, item := range dto.Items {
		p.Spec.Items = append(p.Spec.Items, playlistkind.Item{
			Type:  playlistkind.ItemType(item.Type),
			Value: item.Value,
		})
	}

	return p, nil
}
