package proc

import (
	"github.com/cybriq/p9/version"
)

var F, E, W, I, D, T = GetLogPrinterSet(AddLoggerSubsystem(version.PathBase))
