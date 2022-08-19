package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cybriq/p9/pkg/helpers"
	"github.com/davecgh/go-spew/spew"
	"github.com/gookit/color"
	uberatomic "go.uber.org/atomic"
)

const (
	_Off = iota
	_Fatal
	_Error
	_Chek
	_Warn
	_Info
	_Debug
	_Trace
)

type (
	// LevelPrinter defines a set of terminal printing primitives that output with
	// extra data, time, log logLevelList, and code location
	LevelPrinter struct {
		// Ln prints lists of interfaces with spaces in between
		Ln func(a ...interface{})
		// F prints like fmt.Println surrounded by log details
		F func(format string, a ...interface{})
		// S prints a spew.Sdump for an interface slice
		S func(a ...interface{})
		// C accepts a function so that the extra computation can be avoided if it is
		// not being viewed
		C func(closure func() string)
		// Chk is a shortcut for printing if there is an error, or returning true
		Chk func(e error) bool
	}
	logLevelList struct {
		Off, Fatal, Error, Check, Warn, Info, Debug, Trace int32
	}
	LevelSpec struct {
		ID        int32
		Name      string
		Colorizer func(format string, a ...interface{}) string
	}

	// Entry is a log entry to be printed as json to the log file
	Entry struct {
		Time         time.Time
		Level        string
		Package      string
		CodeLocation string
		Text         string
	}
)

var (
	App          = "   pod"
	AppColorizer = color.White.Sprint
	CurrentLevel = uberatomic.NewInt32(logLevels.Info)
	// writer can be swapped out for any io.*writer* that you want to use instead of
	// stdout.
	writer io.Writer = os.Stderr
	// allSubsystems stores all of the package subsystem names found in the current
	// application
	allSubsystems []string
	// highlighted is a text that helps visually distinguish a log entry by category
	highlighted = make(map[string]struct{})
	// logFilter specifies a set of packages that will not pr logs
	logFilter = make(map[string]struct{})
	// mutexes to prevent concurrent map accesses
	highlightMx, _logFilterMx sync.Mutex
	// logLevels is a shorthand access that minimises possible Name collisions in the
	// dot import
	logLevels = logLevelList{
		Off:   _Off,
		Fatal: _Fatal,
		Error: _Error,
		Check: _Chek,
		Warn:  _Warn,
		Info:  _Info,
		Debug: _Debug,
		Trace: _Trace,
	}
	// LevelSpecs specifies the id, string name and color-printing function
	LevelSpecs = []LevelSpec{
		{logLevels.Off, "off  ", color.Bit24(0, 0, 0, false).Sprintf},
		{logLevels.Fatal, "fatal",
			color.Bit24(128, 0, 0, false).Sprintf},
		{logLevels.Error, "error",
			color.Bit24(255, 0, 0, false).Sprintf},
		{logLevels.Check, "check",
			color.Bit24(255, 255, 0, false).Sprintf},
		{logLevels.Warn, "warn ",
			color.Bit24(0, 255, 0, false).Sprintf},
		{logLevels.Info, "info ",
			color.Bit24(255, 255, 0, false).Sprintf},
		{logLevels.Debug, "debug",
			color.Bit24(0, 128, 255, false).Sprintf},
		{logLevels.Trace, "trace",
			color.Bit24(128, 0, 255, false).Sprintf},
	}
	Levels = []string{
		Off,
		Fatal,
		Error,
		Check,
		Warn,
		Info,
		Debug,
		Trace,
	}
	LogChanDisabled = uberatomic.NewBool(true)
	LogChan         chan Entry
)

const (
	Off   = "off"
	Fatal = "fatal"
	Error = "error"
	Warn  = "warn"
	Info  = "info"
	Check = "check"
	Debug = "debug"
	Trace = "trace"
)

func ListAllSubsystems() []string { return allSubsystems }
func ListAllFilteredSubsystems() (out []string) {
	var counter int
	counter = len(logFilter)
	out = make([]string, counter)
	for i := range logFilter {
		out[counter] = i
	}
	sort.Strings(out)
	return
}
func ListAllHighlightedSubsystems() (out []string) {
	var counter int
	counter = len(highlighted)
	out = make([]string, counter)
	for i := range highlighted {
		out[counter] = i
	}
	sort.Strings(out)
	return
}

// AddLogChan adds a channel that log entries are sent to
func AddLogChan() (ch chan Entry) {
	LogChanDisabled.Store(false)
	if LogChan != nil {
		panic("warning warning")
	}
	// L.Writer.Write.Store( false
	LogChan = make(chan Entry)
	return LogChan
}

// LogPrinters is a struct that bundles a set of log printers for a subsystem
type LogPrinters struct {
	F, E, W, I, D, T LevelPrinter
}

// GetLogPrinters returns a set of log printers wrapped in a struct
func GetLogPrinters(subsystem string) (log LogPrinters) {
	log.F, log.E, log.W, log.I, log.D, log.T = GetLogPrinterSet(subsystem)
	return
}

// GetLogPrinterSet returns a set of LevelPrinter with their subsystem preloaded
func GetLogPrinterSet(subsystem string) (Fatal, Error, Warn, Info, Debug, Trace LevelPrinter) {
	return GetOnePrinter(_Fatal, subsystem),
		GetOnePrinter(_Error, subsystem),
		GetOnePrinter(_Warn, subsystem),
		GetOnePrinter(_Info, subsystem),
		GetOnePrinter(_Debug, subsystem),
		GetOnePrinter(_Trace, subsystem)
}

func GetOnePrinter(level int32, subsystem string) LevelPrinter {
	return LevelPrinter{
		Ln:  GetPrintln(level, subsystem),
		F:   GetPrintf(level, subsystem),
		S:   GetPrints(level, subsystem),
		C:   GetPrintc(level, subsystem),
		Chk: GetChk(level, subsystem),
	}
}

// SetLogLevel sets the log level via a string, which can be truncated down to
// one character, similar to nmcli's argument processor, as the first letter is
// unique. This could be used with a linter to make larger command sets.
func SetLogLevel(l string) {
	if l == "" {
		l = "info"
	}
	// fmt.Fprintln(os.Stderr, "setting log level", l)
	lvl := logLevels.Info
	for i := range LevelSpecs {
		if LevelSpecs[i].Name[:1] == l[:1] {
			lvl = LevelSpecs[i].ID
		}
	}
	CurrentLevel.Store(lvl)
}

// SetLogWriter atomically changes the log io.Writer interface
func SetLogWriter(wr io.Writer) {
	// w := unsafe.Pointer(writer)
	// c := unsafe.Pointer(wr)
	// atomic.SwapPointer(&w, c)
	writer = wr
}

func SetLogWriteToFile(path, appName string) (e error) {
	// copy existing log file to dated log file as we will truncate it per
	// session
	path = filepath.Join(path, "log"+appName)
	if _, e = os.Stat(path); e == nil {
		var b []byte
		b, e = ioutil.ReadFile(path)
		if e == nil {
			ioutil.WriteFile(path+fmt.Sprint(time.Now().Unix()), b,
				0600)
		}
	}
	var fileWriter *os.File
	if fileWriter, e = os.OpenFile(
		path, os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0600,
	); e != nil {
		// fmt.Fprintln(os.Stderr, "unable to write log to", path, "error:", e)
		return
	}
	mw := io.MultiWriter(os.Stderr, fileWriter)
	// fileWriter.Write([]byte("logging to file '" + path + "'\n"))
	// mw.Write([]byte("logging to file '" + path + "'\n"))
	SetLogWriter(mw)
	return
}

// SortSubsystemsList sorts the list of subsystems, to keep the data read-only,
// call this function right at the top of the main, which runs after
// declarations and main/init. Really this is just here to alert the reader.
func SortSubsystemsList() {
	sort.Strings(allSubsystems)
	// fmt.Fprintln(
	// 	os.Stderr,
	// 	spew.Sdump(allSubsystems),
	// 	spew.Sdump(highlighted),
	// 	spew.Sdump(logFilter),
	// )
}

// AddLoggerSubsystem adds a subsystem to the list of known subsystems and returns the
// string so it is nice and neat in the package logg.go file
func AddLoggerSubsystem(pathBase string) (subsystem string) {
	// var split []string
	var ok bool
	var file string
	_, file, _, ok = runtime.Caller(1)
	if ok {
		r := strings.Split(file, pathBase)
		// fmt.Fprintln(os.Stderr, version.PathBase, r)
		fromRoot := filepath.Base(file)
		if len(r) > 1 {
			fromRoot = r[1]
		}
		split := strings.Split(fromRoot, "/")
		// fmt.Fprintln(os.Stderr, version.PathBase, "file", file, r, fromRoot, split)
		subsystem = strings.Join(split[:len(split)-1], "/")
		// fmt.Fprintln(os.Stderr, "adding subsystem", subsystem)
		allSubsystems = append(allSubsystems, subsystem)
	}
	return
}

// StoreHighlightedSubsystems sets the list of subsystems to highlight
func StoreHighlightedSubsystems(highlights []string) (found bool) {
	highlightMx.Lock()
	highlighted = make(map[string]struct{}, len(highlights))
	for i := range highlights {
		highlighted[highlights[i]] = struct{}{}
	}
	highlightMx.Unlock()
	return
}

// LoadHighlightedSubsystems returns a copy of the map of highlighted subsystems
func LoadHighlightedSubsystems() (o []string) {
	highlightMx.Lock()
	o = make([]string, len(logFilter))
	var counter int
	for i := range logFilter {
		o[counter] = i
		counter++
	}
	highlightMx.Unlock()
	sort.Strings(o)
	return
}

// StoreSubsystemFilter sets the list of subsystems to filter
func StoreSubsystemFilter(filter []string) {
	_logFilterMx.Lock()
	logFilter = make(map[string]struct{}, len(filter))
	for i := range filter {
		logFilter[filter[i]] = struct{}{}
	}
	_logFilterMx.Unlock()
}

// LoadSubsystemFilter returns a copy of the map of filtered subsystems
func LoadSubsystemFilter() (o []string) {
	_logFilterMx.Lock()
	o = make([]string, len(logFilter))
	var counter int
	for i := range logFilter {
		o[counter] = i
		counter++
	}
	_logFilterMx.Unlock()
	sort.Strings(o)
	return
}

// IsHighlighted returns true if the subsystem is in the list to have attention
// getters added to them
func IsHighlighted(subsystem string) (found bool) {
	highlightMx.Lock()
	_, found = highlighted[subsystem]
	highlightMx.Unlock()
	return
}

// AddHighlightedSubsystem adds a new subsystem Name to the highlighted list
func AddHighlightedSubsystem(hl string) struct{} {
	highlightMx.Lock()
	highlighted[hl] = struct{}{}
	highlightMx.Unlock()
	return struct{}{}
}

// IsSubsystemFiltered returns true if the subsystem should not pr logs
func IsSubsystemFiltered(subsystem string) (found bool) {
	_logFilterMx.Lock()
	_, found = logFilter[subsystem]
	_logFilterMx.Unlock()
	return
}

// AddFilteredSubsystem adds a new subsystem Name to the highlighted list
func AddFilteredSubsystem(hl string) struct{} {
	_logFilterMx.Lock()
	logFilter[hl] = struct{}{}
	_logFilterMx.Unlock()
	return struct{}{}
}

func getTimeText(level int32) string {
	// since := time.Now().Sub(logger_started).Round(time.Millisecond).String()
	// diff := 12 - len(since)
	// if diff > 0 {
	// 	since = strings.Repeat(" ", diff) + since + " "
	// }
	return color.Bit24(99, 99, 99, false).Sprint(
		time.Now().
			Format(time.StampMilli),
	)
}

func GetPrintln(level int32, subsystem string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if level <= CurrentLevel.Load() && !IsSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if IsHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%-6v %s\n",
					helpers.GetLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(helpers.JoinStrings(" ",
						a...)),
				),
			)
		}
	}
}

func GetPrintf(level int32, subsystem string) func(
	format string,
	a ...interface{},
) {
	return func(format string, a ...interface{}) {
		if level <= CurrentLevel.Load() && !IsSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if IsHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%-6v %s\n",
					helpers.GetLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(fmt.Sprintf(format, a...)),
				),
			)
		}
	}
}

func GetPrints(level int32, subsystem string) func(a ...interface{}) {
	return func(a ...interface{}) {
		if level <= CurrentLevel.Load() && !IsSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if IsHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%s%s%s\n",
					helpers.GetLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(
						" spew:",
					),
					fmt.Sprint(
						color.Bit24(
							20,
							20,
							20,
							true,
						).Sprint("\n\n"+spew.Sdump(a)),
						"\n",
					),
				),
			)
		}
	}
}

func GetPrintc(level int32, subsystem string) func(closure func() string) {
	return func(closure func() string) {
		if level <= CurrentLevel.Load() && !IsSubsystemFiltered(subsystem) {
			printer := fmt.Sprintf
			if IsHighlighted(subsystem) {
				printer = color.Bold.Sprintf
			}
			fmt.Fprintf(
				writer,
				printer(
					"%-58v%s%s%-6v %s\n",
					helpers.GetLoc(2, level, subsystem),
					getTimeText(level),
					color.Bit24(20, 20, 20, true).
						Sprint(AppColorizer(" "+App)),
					LevelSpecs[level].Colorizer(
						color.Bit24(20, 20, 20, true).
							Sprint(" "+LevelSpecs[level].Name+" "),
					),
					AppColorizer(closure()),
				),
			)
		}
	}
}

func GetChk(level int32, subsystem string) func(e error) bool {
	return func(e error) bool {
		if level <= CurrentLevel.Load() && !IsSubsystemFiltered(subsystem) {
			if e != nil {
				printer := fmt.Sprintf
				if IsHighlighted(subsystem) {
					printer = color.Bold.Sprintf
				}
				fmt.Fprintf(
					writer,
					printer(
						"%-58v%s%s%-6v %s\n",
						helpers.GetLoc(2, level,
							subsystem),
						getTimeText(level),
						color.Bit24(20, 20, 20, true).
							Sprint(AppColorizer(" "+App)),
						LevelSpecs[level].Colorizer(
							color.Bit24(20, 20, 20,
								true).
								Sprint(" "+LevelSpecs[level].Name+" "),
						),
						LevelSpecs[level].Colorizer(
							helpers.JoinStrings(
								" ",
								e.Error(),
							),
						),
					),
				)
				return true
			}
		}
		return false
	}
}
