package ui

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (ui *UI) drawLog(gtx layout.Context, th *material.Theme) layout.Dimensions {
	/*return layout.UniformInset(8).Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, main.bufClient.String()+main.bufServer).Layout(gtx)
		},
	)*/
	return layout.Dimensions{}
}
