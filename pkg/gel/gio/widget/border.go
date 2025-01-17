// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image/color"

	"github.com/cybriq/p9/pkg/gel/gio/f32"
	"github.com/cybriq/p9/pkg/gel/gio/layout"
	"github.com/cybriq/p9/pkg/gel/gio/op/clip"
	"github.com/cybriq/p9/pkg/gel/gio/op/paint"
	"github.com/cybriq/p9/pkg/gel/gio/unit"
)

// Border lays out a widget and draws a border inside it.
type Border struct {
	Color        color.NRGBA
	CornerRadius unit.Value
	Width        unit.Value
}

func (b Border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	sz := layout.FPt(dims.Size)

	rr := float32(gtx.Px(b.CornerRadius))
	width := float32(gtx.Px(b.Width))
	sz.X -= width
	sz.Y -= width

	r := f32.Rectangle{Max: sz}
	r = r.Add(f32.Point{X: width * 0.5, Y: width * 0.5})

	paint.FillShape(
		gtx.Ops,
		b.Color,
		clip.Stroke{
			Path:  clip.UniformRRect(r, rr).Path(gtx.Ops),
			Style: clip.StrokeStyle{Width: width},
		}.Op(),
	)

	return dims
}
