package ui

import (
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (ui *UI) drawAbout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	q := strconv.Itoa(len(ui.workingAuto))
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 20}.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Label(th, 18, "В автоматическом режиме работает "+q+" процессов").Layout(gtx)
						})
					})
			}),
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Flexed(0.7,
						func(gtx layout.Context) layout.Dimensions {
							list := &layout.List{Axis: layout.Vertical}
							return layout.Inset{Top: 10}.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Center.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return list.Layout(gtx, len(ui.workingAuto), func(gtx layout.Context, i int) layout.Dimensions {
												k := ui.workingAuto[i]
												return layout.Inset{Bottom: 7}.Layout(gtx,
													func(gtx layout.Context) layout.Dimensions {
														return material.Label(th, 16, k).Layout(gtx)
													})
											})
										},
									)
								})
						}))
			}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			b := material.Body1(th, "deviceId: "+ui.deviceId)
			b.State = &ui.deviceSelectable
			return b.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			now := time.Now()
			zone, _ := now.Zone()

			return material.Body1(th, zone).Layout(gtx)
		}),
	)
}
