// SPDX-License-Identifier: Unlicense OR MIT

package clipboard

import (
	"github.com/cybriq/p9/pkg/gel/gio/internal/opconst"
	"github.com/cybriq/p9/pkg/gel/gio/io/event"
	"github.com/cybriq/p9/pkg/gel/gio/op"
)

// Event is generated when the clipboard content is requested.
type Event struct {
	Text string
}

// ReadOp requests the text of the clipboard, delivered to
// the current handler through an Event.
type ReadOp struct {
	Tag event.Tag
}

// WriteOp copies Text to the clipboard.
type WriteOp struct {
	Text string
}

func (h ReadOp) Add(o *op.Ops) {
	data := o.Write1(opconst.TypeClipboardReadLen, h.Tag)
	data[0] = byte(opconst.TypeClipboardRead)
}

func (h WriteOp) Add(o *op.Ops) {
	data := o.Write1(opconst.TypeClipboardWriteLen, &h.Text)
	data[0] = byte(opconst.TypeClipboardWrite)
}

func (Event) ImplementsEvent() {}
