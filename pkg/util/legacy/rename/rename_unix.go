//go:build !windows && !plan9
// +build !windows,!plan9

package rename

import (
	"os"
)

// Atomic provides an atomic file rename. newpath is replaced if it already exists.
func Atomic(oldpath, newpath string) (e error) {
	return os.Rename(oldpath, newpath)
}
