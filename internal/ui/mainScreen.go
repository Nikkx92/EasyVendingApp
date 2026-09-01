package ui

import (
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func (ui *UI) drawMain(gtx layout.Context, th *material.Theme) layout.Dimensions {

	return layout.Inset{Left: 7, Right: 7, Top: 20}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, 18, "Выберите дату:")
						return layout.Inset{Top: 15, Bottom: 10}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return lbl.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 15, Left: 20, Right: 20, Bottom: 50}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								return ui.datePicker.Layout(gtx, th)
							})
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						if ui.datePicker.Visible {
							return layout.Dimensions{}
						} else {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis: layout.Horizontal,
										}.Layout(gtx,
											layout.Rigid(
												func(gtx layout.Context) layout.Dimensions {
													return material.Label(th, 16, "Автоматическая отправка продаж").Layout(gtx)
												},
											),
											layout.Rigid(
												func(gtx layout.Context) layout.Dimensions {
													return layout.Inset{Top: unit.Dp(2), Left: 8}.Layout(gtx,
														func(gtx layout.Context) layout.Dimensions {
															return material.Switch(th, ui.AutoSwt, "autoSwitch").Layout(gtx)
														})
												},
											),
										)
									},
								),
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return layout.Center.Layout(gtx,
											func(gtx layout.Context) layout.Dimensions {
												return layout.Inset{Top: 15}.Layout(gtx,
													func(gtx layout.Context) layout.Dimensions {
														var str string
														if ui.emptyFields {
															str = "вы не авторизованы"
														} else if ui.dateInvalid {
															str = "неверный формат даты"
														} else {
															str = ""
														}
														lb := material.Label(th, 18, str)
														return lb.Layout(gtx)
													},
												)
											},
										)
									}),
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 30}.Layout(gtx,
											func(gtx layout.Context) layout.Dimensions {
												return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													gtx.Constraints.Min.X = gtx.Dp(230)
													btn := material.Button(th, ui.requestKitBtn, "получить продажи")
													btn.CornerRadius = unit.Dp(10)

													return btn.Layout(gtx)

												})
											},
										)
									},
								))
						}
					}),
			)
		})
}

func (ui *UI) modalResponse(gtx layout.Context, m []string, s string) {
	ui.ModalResp.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		if ui.isLoad {
			return layout.Center.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X / 4
					return material.Loader(th).Layout(gtx)
				},
			)
		} else if m == nil && s == "" {
			ui.ModalResp.VisibilityAnimation.Disappear(time.Now())
			gtx.Execute(op.InvalidateCmd{})
			return layout.Dimensions{}
		} else {
			gtx.Constraints.Min = image.Point{}

			recordGtx := gtx
			mOps := new(op.Ops)
			recordGtx.Ops = mOps
			macro := op.Record(mOps)

			dims := layout.UniformInset(15).Layout(recordGtx, func(gtx layout.Context) layout.Dimensions {
				return ui.click.Layout(recordGtx,
					func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									if s != "" {
										return material.Label(th, 16, s).Layout(gtx)
									} else if m != nil {
										return layout.Flex{
											Axis: layout.Horizontal,
										}.Layout(gtx,
											layout.Rigid(
												ui.drawModalMessage("Количество новых продаж:", th),
											),
											layout.Rigid(
												ui.drawModalMessage(strconv.Itoa(len(m)), th),
											),
										)
									} else {
										return material.Label(th, 16, "нет новых продаж").Layout(gtx)
									}
								},
							),
							layout.Rigid(
								func(gtx layout.Context) layout.Dimensions {
									if m == nil {
										return layout.Dimensions{}
									} else {
										return ui.pairOfButtons(30, ui.sendToFnsBtn, ui.cancelBtn, "отправить в ФНС", "не отправлять", th, gtx)
									}
								},
							),
						)
					})
			})

			stop := macro.Stop()

			defer op.Offset(image.Pt((gtx.Constraints.Max.X-dims.Size.X)/2, (gtx.Constraints.Max.Y-dims.Size.Y)/2)).Push(gtx.Ops).Pop()
			rect := clip.UniformRRect(image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)}, 12)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, rect.Op(gtx.Ops))

			stop.Add(gtx.Ops)

			return dims

		}
	}

	ui.ModalResp.FinalAlpha = 245
	ui.ModalResp.Appear(gtx.Now)

}

func (ui *UI) drawModalMessage(s string, th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Left: 10, Right: 10}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						if s == "" {
							return layout.Dimensions{}
						} else {
							return material.Label(th, 16, s).Layout(gtx)
						}
					},
				)
			},
		)
	}
}
