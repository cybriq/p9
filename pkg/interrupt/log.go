package interrupt

import (
	"github.com/cybriq/p9/pkg/log"
	"github.com/cybriq/p9/version"
)

var F, E, W, I, D, T = log.GetLogPrinterSet(log.AddLoggerSubsystem(version.PathBase))
