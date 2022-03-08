package main

import (
	l "github.com/cybriq/p9/pkg/gel/gio/layout"

	"github.com/cybriq/p9/pkg/qu"

	"github.com/cybriq/p9/cmd/misc/glom/pkg/pathtree"
	"github.com/cybriq/p9/pkg/gel"
	"github.com/cybriq/p9/pkg/interrupt"
)

type State struct {
	*gel.Window
}

func NewState(quit qu.C) *State {
	return &State{
		Window: gel.NewWindowP9(quit),
	}
}

func main() {
	quit := qu.T()
	state := NewState(quit)
	var e error
	folderView := pathtree.New(state.Window)
	state.Window.SetDarkTheme(folderView.Dark.True())
	if e = state.Window.
		Size(48, 32).
		Title("glom, the visual code editor").
		Open().
		Run(func(gtx l.Context) l.Dimensions { return folderView.Fn(gtx) },
			func() {
				interrupt.Request()
				quit.Q()
			}, quit,
		); E.Chk(e) {

	}
}
