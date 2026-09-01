package ui

import (
	"image"
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func (ui *UI) drawSales(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if ui.sales == nil {
		img := widget.Image{
			Src:   ui.emptySalesIcon,
			Scale: 0.6,
			Fit:   widget.ScaleDown,
		}
		return img.Layout(gtx)
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(7).Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return material.Label(th, 18, "Статистика продаж").Layout(gtx)
							},
						)
					},
				)
			},
		),
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(10).Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(4),
							Width:        unit.Dp(2),
						}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(gtx,
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{
												Axis: layout.Horizontal,
											}.Layout(gtx,
												layout.Flexed(0.5,
													func(gtx layout.Context) layout.Dimensions {
														return layout.Center.Layout(gtx,
															func(gtx layout.Context) layout.Dimensions {
																return material.Label(th, 16, "Дата").Layout(gtx)
															})
													}),
												layout.Flexed(0.5,
													func(gtx layout.Context) layout.Dimensions {
														return layout.Center.Layout(gtx,
															func(gtx layout.Context) layout.Dimensions {
																return material.Label(th, 16, "Количество").Layout(gtx)
															})
													}))
										}),
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											d := component.Divider(th)
											d.Thickness = unit.Dp(2)
											return d.Layout(gtx)
										}),
									layout.Rigid(
										func(gtx layout.Context) layout.Dimensions {
											var dims layout.Dimensions
											var cs2 layout.Constraints
											return material.List(th, &ui.list).Layout(gtx, 1,
												func(gtx layout.Context, i int) layout.Dimensions {
													return layout.Flex{
														Axis: layout.Horizontal,
													}.Layout(gtx,
														layout.Flexed(0.5,
															func(gtx layout.Context) layout.Dimensions {
																list := &layout.List{Axis: layout.Vertical}
																return layout.Center.Layout(gtx,
																	func(gtx layout.Context) layout.Dimensions {
																		return list.Layout(gtx, len(ui.sales), func(gtx layout.Context, i int) layout.Dimensions {
																			return layout.Inset{Bottom: 7}.Layout(gtx,
																				func(gtx layout.Context) layout.Dimensions {
																					dims = material.Label(th, 20, ui.sales[i].Date).Layout(gtx)
																					size := dims.Size
																					cs2 = layout.Constraints{Min: size, Max: size}
																					return dims
																				})
																		})
																	},
																)
															},
														),
														layout.Flexed(0.5,
															func(gtx layout.Context) layout.Dimensions {
																list := &layout.List{Axis: layout.Vertical}
																return layout.Center.Layout(gtx,
																	func(gtx layout.Context) layout.Dimensions {
																		return list.Layout(gtx, len(ui.sales),
																			func(gtx layout.Context, i int) layout.Dimensions {
																				return layout.Inset{Bottom: 7}.Layout(gtx,
																					func(gtx layout.Context) layout.Dimensions {
																						gtx.Constraints = cs2
																						return layout.Center.Layout(gtx,
																							func(gtx layout.Context) layout.Dimensions {
																								return material.Clickable(gtx, &ui.clkLabel[i],
																									func(gtx layout.Context) layout.Dimensions {
																										s := material.Label(th, 20, strconv.Itoa(len(ui.sales[i].Details.Drinks))).Layout(gtx)
																										h := 2
																										r := image.Rect(0, s.Size.Y-h, s.Size.X, s.Size.Y)
																										paint.FillShape(gtx.Ops, color.NRGBA{R: 1, G: 1, B: 1, A: 255}, clip.Rect(r).Op())
																										return s
																									})
																							})
																					})
																			})
																	},
																)
															},
														),
													)

												})
										},
									),
								)
							})
					},
				)
			}),
		/*layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: 20, Right: 20, Bottom: 10}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, ui.moreBtn, "загрузить еще").Layout(gtx)
				})
		}),*/
	)
}

func (ui *UI) modalDetail(gtx layout.Context, m []string, date string) {
	ui.ModalDetails.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {

		gtx.Constraints.Min = image.Point{}

		recordGtx := gtx
		mOps := new(op.Ops)
		recordGtx.Ops = mOps
		macro := op.Record(mOps)

		dims := layout.Inset{Left: unit.Dp(26), Right: unit.Dp(26), Top: unit.Dp(6), Bottom: unit.Dp(26)}.Layout(recordGtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(recordGtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Bottom: 20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.Body1(th, date).Layout(gtx)
							})
						},
					),
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							list := &layout.List{Axis: layout.Vertical}
							return list.Layout(gtx, len(m), func(gtx layout.Context, i int) layout.Dimensions {
								return material.Body1(th, m[i]).Layout(gtx)
							})
						}),
				)
			})

		stop := macro.Stop()

		defer op.Offset(image.Pt((gtx.Constraints.Max.X-dims.Size.X)/2, (gtx.Constraints.Max.Y-dims.Size.Y)/2)).Push(gtx.Ops).Pop()
		rect := clip.UniformRRect(image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)}, 12)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, rect.Op(gtx.Ops))

		stop.Add(gtx.Ops)

		return dims

	}
	ui.ModalDetails.FinalAlpha = 245
	ui.ModalDetails.Appear(gtx.Now)
}
