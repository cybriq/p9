package helpers

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

// JoinStrings constructs a string from an slice of interface same as Println but
// without the terminal newline
func JoinStrings(sep string, a ...interface{}) (o string) {
	for i := range a {
		o += fmt.Sprint(a[i])
		if i < len(a)-1 {
			o += sep
		}
	}
	return
}

// GetLoc calls runtime.Caller and formats as expected by source code editors
// for terminal hyperlinks
//
// Regular expressions and the substitution texts to make these clickable in
// Tilix and other RE hyperlink configurable terminal emulators:
//
// This matches the shortened paths generated in this command and printed at
// the very beginning of the line as this logger prints:
//
// ^((([\/a-zA-Z@0-9-_.]+/)+([a-zA-Z@0-9-_.]+)):([0-9]+))
//
// 		goland --line $5 $GOPATH/src/github.com/p9c/matrjoska/$2
//
// I have used a shell variable there but tilix doesn't expand them,
// so put your GOPATH in manually, and obviously change the repo subpath.
//
//
// Change the path to use with another repository's logging output (
// someone with more time on their hands could probably come up with
// something, but frankly the custom links feature of Tilix has the absolute
// worst UX I have encountered since the 90s...
// Maybe in the future this library will be expanded with a tool that more
// intelligently sets the path, ie from CWD or other cleverness.
//
// This matches full paths anywhere on the commandline delimited by spaces:
//
// ([/](([\/a-zA-Z@0-9-_.]+/)+([a-zA-Z@0-9-_.]+)):([0-9]+))
//
// 		goland --line $5 /$2
//
// Adapt the invocation to open your preferred editor if it has the capability,
// the above is for Jetbrains Goland
//
func GetLoc(skip int, level int32, subsystem string) (output string) {
	_, file, line, _ := runtime.Caller(skip)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(
				os.Stderr,
				"getloc panic on subsystem",
				subsystem,
				file,
			)
		}
	}()
	split := strings.Split(file, subsystem)
	if len(split) < 2 {
		output = fmt.Sprint(
			color.White.Sprint(subsystem),
			color.Gray.Sprint(
				file, ":", line,
			),
		)
	} else {
		output = fmt.Sprint(
			color.White.Sprint(subsystem),
			color.Gray.Sprint(
				split[1], ":", line,
			),
		)
	}
	return
}

// DirectionString is a helper function that returns a string that represents the direction of a connection (inbound or outbound).
func DirectionString(inbound bool) string {
	if inbound {
		return "inbound"
	}
	return "outbound"
}

func PickNoun(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

func FileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return e == nil
}

func Caller(comment string, skip int) string {
	_, file, line, _ := runtime.Caller(skip + 1)
	o := fmt.Sprintf("%s: %s:%d", comment, file, line)
	// L.Debug(o)
	return o
}
