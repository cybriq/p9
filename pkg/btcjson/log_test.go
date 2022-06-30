package btcjson_test

import (
	"github.com/cybriq/p9/pkg/proc"
	"github.com/cybriq/p9/version"
)

var subsystem = proc.AddLoggerSubsystem(version.PathBase)
var F, E, W, I, D, T proc.LevelPrinter = proc.GetLogPrinterSet(subsystem)

func init() {
	// var _ = log.AddFilteredSubsystem(subsystem)
	// var _ = log.AddHighlightedSubsystem(subsystem)
	// F.Ln("F.Ln")
	// E.Ln("E.Ln")
	// W.Ln("W.Ln")
	// I.Ln("I.Ln")
	// D.Ln("D.Ln")
	// F.Ln("T.Ln")
	// F.F("%s", "F.F")
	// E.F("%s", "E.F")
	// W.F("%s", "W.F")
	// I.F("%s", "I.F")
	// D.F("%s", "D.F")
	// T.F("%s", "T.F")
	// F.C(func() string { return "ftl.C" })
	// E.C(func() string { return "err.C" })
	// W.C(func() string { return "W.C" })
	// I.C(func() string { return "inf.C" })
	// D.C(func() string { return "D.C" })
	// T.C(func() string { return "T.C" })
	// F.C(func() string { return "ftl.C" })
	// E.Chk(errors.New("E.Chk"))
	// W.Chk(errors.New("W.Chk"))
	// I.Chk(errors.New("inf.Chk"))
	// D.Chk(errors.New("D.Chk"))
	// T.Chk(errors.New("T.Chk"))
}
