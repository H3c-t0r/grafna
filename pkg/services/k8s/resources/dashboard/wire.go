package dashboard

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideWatcher,
	wire.Bind(new(Watcher), new(*watcher)),
	ProvideStoreWrapper, // Replace the original store with a wrapper
)
