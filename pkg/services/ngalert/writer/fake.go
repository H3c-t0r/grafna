package writer

import (
	"context"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type FakeWriter struct {
	WriteFunc func(ctx context.Context, name string, t time.Time, frames data.Frames) error
}

func (w FakeWriter) Write(ctx context.Context, name string, t time.Time, frames data.Frames) error {
	if w.WriteFunc == nil {
		return nil
	}

	return w.WriteFunc(ctx, name, t, frames)
}
