package ctrl

import (
	"github.com/cybriq/p9/pkg/log"
	"github.com/cybriq/p9/version"
)

var subsystem = log.AddLoggerSubsystem(version.PathBase)
var F, E, W, I, D, T log.LevelPrinter = log.GetLogPrinterSet(subsystem)
