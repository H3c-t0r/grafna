package dashboard

import "context"

var _ Watcher = (*watcher)(nil)

type watcher struct{}

func ProvideWatcher() *watcher {
	return &watcher{}
}

func (w *watcher) Add(ctx context.Context, obj *Dashboard) {
	// TODO
}

func (w *watcher) Update(ctx context.Context, oldObj, newObj *Dashboard) {
	// TODO
}

func (w *watcher) Delete(ctx context.Context, obj *Dashboard) {
	// TODO
}
