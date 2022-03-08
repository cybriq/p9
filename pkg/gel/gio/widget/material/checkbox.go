// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"github.com/cybriq/p9/pkg/gel/gio/layout"
	"github.com/cybriq/p9/pkg/gel/gio/unit"
	"github.com/cybriq/p9/pkg/gel/gio/widget"
)

type CheckBoxStyle struct {
	checkable
	CheckBox *widget.Bool
}

func CheckBox(th *Theme, checkBox *widget.Bool, label string) CheckBoxStyle {
	return CheckBoxStyle{
		CheckBox: checkBox,
		checkable: checkable{
			Label:              label,
			Color:              th.Palette.Fg,
			IconColor:          th.Palette.ContrastBg,
			TextSize:           th.TextSize.Scale(14.0 / 16.0),
			Size:               unit.Dp(26),
			shaper:             th.Shaper,
			checkedStateIcon:   th.Icon.CheckBoxChecked,
			uncheckedStateIcon: th.Icon.CheckBoxUnchecked,
		},
	}
}

// Layout updates the checkBox and displays it.
func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	dims := c.layout(gtx, c.CheckBox.Value, c.CheckBox.Hovered())
	gtx.Constraints.Min = dims.Size
	c.CheckBox.Layout(gtx)
	return dims
}
