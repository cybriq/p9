package main

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/cybriq/p9/pkg/qu"

	_ "github.com/cybriq/p9/pkg/gel/gio/app/permission/networkstate" // todo: integrate this into routeable package
	_ "github.com/cybriq/p9/pkg/gel/gio/app/permission/storage"      // this enables the home folder appdata directory to work on android (and ios)

	"github.com/cybriq/p9/pkg/log"
	"github.com/cybriq/p9/pod/config"
	"github.com/cybriq/p9/pod/podcfgs"
	"github.com/cybriq/p9/pod/podhelp"
	"github.com/cybriq/p9/pod/state"
	"github.com/cybriq/p9/version"

	// This ensures the database drivers get registered
	_ "github.com/cybriq/p9/pkg/database/ffldb"
	// _ "github.com/cybriq/p9/pkg/gel/gio/app/permission/bluetooth"
	// _ "github.com/cybriq/p9/pkg/gel/gio/app/permission/camera"
)

func main() {
	<-Main()
}

func Main() (quit qu.C) {
	quit = qu.T()
	go func() {
		log.SetLogLevel("trace")
		T.Ln(os.Args)
		T.Ln(version.Get())
		var cx *state.State
		var e error
		if cx, e = state.GetNew(podcfgs.GetDefaultConfig(),
			podhelp.HelpFunction, quit,
		); E.Chk(e) {
			fail()
		}

		// fail()
		// if e = debugConfig(cx.Config); E.Chk(e) {
		// }

		D.Ln("running command:", cx.Config.RunningCommand.Name)
		if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
			fail()
		}
		quit.Q()
	}()
	return quit
}

func fail() {
	os.Exit(1)
}

func debugConfig(c *config.Config) (e error) {
	c.ShowAll = true
	defer func() { c.ShowAll = false }()
	var j []byte
	if j, e = c.MarshalJSON(); E.Chk(e) {
		return
	}
	var b []byte
	jj := bytes.NewBuffer(b)
	if e = json.Indent(jj, j, "", "\t"); E.Chk(e) {
		return
	}
	T.Ln("\n" + jj.String())
	return
}
