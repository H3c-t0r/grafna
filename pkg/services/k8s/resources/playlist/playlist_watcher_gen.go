// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     CRDWatcherJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package playlist

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/log"
)

type Watcher interface {
	Add(context.Context, *Playlist)
	Update(context.Context, *Playlist, *Playlist)
	Delete(context.Context, *Playlist)
}

type WatcherWrapper struct {
	log     log.Logger
	watcher Watcher
}

func NewWatcherWrapper(watcher Watcher) *WatcherWrapper {
	return &WatcherWrapper{
		log:     log.New("k8s.playlist.watcher"),
		watcher: watcher,
	}
}

func (w *WatcherWrapper) Add(ctx context.Context, obj any) {
	conv, err := fromUnstructured(obj)
	if err != nil {
		w.log.Error("Failed to convert object", "err", err)
		return
	}
	w.watcher.Add(ctx, conv)
}

func (w *WatcherWrapper) Update(ctx context.Context, oldObj, newObj any) {
	convOld, err := fromUnstructured(oldObj)
	if err != nil {
		w.log.Error("Failed to convert oldObj", "err", err)
		return
	}
	convNew, err := fromUnstructured(newObj)
	if err != nil {
		w.log.Error("Failed to convert newObj", "err", err)
		return
	}
	w.watcher.Update(ctx, convOld, convNew)
}

func (w *WatcherWrapper) Delete(ctx context.Context, obj any) {
	conv, err := fromUnstructured(obj)
	if err != nil {
		w.log.Error("Failed to convert object", "err", err)
		return
	}
	w.watcher.Delete(ctx, conv)
}
