package node

import (
	"github.com/cybriq/p9/pkg/proc"
	"github.com/cybriq/p9/version"
)

var subsystem = proc.AddLoggerSubsystem(version.PathBase)
var F, E, W, I, D, T proc.LevelPrinter = proc.GetLogPrinterSet(subsystem)
