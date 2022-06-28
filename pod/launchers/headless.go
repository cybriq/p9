//go:build headless
// +build headless

package launchers

import (
	"os"
)

func GUIHandle(ifc interface{}) (e error) {
	W.Ln("GUI was disabled for this build (server only version)")
	os.Exit(1)
	return nil
}
