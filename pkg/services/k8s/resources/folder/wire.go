package folder

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideWatcher,
)
