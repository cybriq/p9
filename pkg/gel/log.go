package gel

import (
	"github.com/cybriq/p9/pkg/proc"

	"github.com/cybriq/p9/version"
)

var F, E, W, I, D, T = proc.GetLogPrinterSet(proc.AddLoggerSubsystem(version.PathBase))
