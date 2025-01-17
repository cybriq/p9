// SPDX-License-Identifier: Unlicense OR MIT

// Package profiles provides access to rendering
// profiles.
package profile

import (
	"github.com/cybriq/p9/pkg/gel/gio/internal/opconst"
	"github.com/cybriq/p9/pkg/gel/gio/io/event"
	"github.com/cybriq/p9/pkg/gel/gio/op"
)

// Op registers a handler for receiving
// Events.
type Op struct {
	Tag event.Tag
}

// Event contains profile data from a single
// rendered frame.
type Event struct {
	// Timings. Very likely to change.
	Timings string
}

func (p Op) Add(o *op.Ops) {
	data := o.Write1(opconst.TypeProfileLen, p.Tag)
	data[0] = byte(opconst.TypeProfile)
}

func (p Event) ImplementsEvent() {}
