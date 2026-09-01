package ui

import (
	"image"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func (ui *UI) drawAuth(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(50)}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			if ui.isAuth {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(
						drawHeaders("Выполнен вход под логином: "+ui.Fields[1].Text(), th),
					),
					layout.Rigid(
						drawHeaders("Выполнен вход под ИНН: "+ui.Fields[3].Text(), th),
					),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return ui.pairOfButtons(25, ui.exitLogin, ui.deleteData, "выйти", "удалить данные", th, gtx)
					}),
				)
			} else {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.List(th, &ui.list).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
							return layout.Flex{
								Axis: layout.Vertical,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									c := component.Surface(th)
									c.CornerRadius = unit.Dp(20)
									c.Elevation = 4
									c.UmbraColor = c.AmbientColor

									return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.UniformInset(12).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{
													Axis:      layout.Vertical,
													Alignment: layout.Middle,
												}.Layout(gtx,
													layout.Rigid(
														drawHeaders("Данные KitVending", th),
													),
													layout.Rigid(
														ui.drawFields("ID компании", "1,2,3,4,5,6,7,8,9,0", key.HintNumeric, &ui.Fields[0], th)),
													layout.Rigid(
														ui.drawFields("Логин", "", key.HintAny, &ui.Fields[1], th)),
													layout.Rigid(
														ui.drawFields("Пароль", "", key.HintAny, &ui.Fields[2], th)),
												)
											})
										})
									})
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Spacer{Height: unit.Dp(50)}.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									c := component.Surface(th)
									c.CornerRadius = unit.Dp(20)
									c.Elevation = 4
									c.UmbraColor = c.AmbientColor

									return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.UniformInset(12).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Flex{
													Axis:      layout.Vertical,
													Alignment: layout.Middle,
												}.Layout(gtx,
													layout.Rigid(
														drawHeaders("Данные Мой Налог", th),
													),
													layout.Rigid(
														ui.drawFields("ИНН", "1,2,3,4,5,6,7,8,9,0", key.HintNumeric, &ui.Fields[3], th)),
													layout.Rigid(
														ui.drawFields("Пароль", "", key.HintAny, &ui.Fields[4], th)),
												)
											})
										})
									})
								}),
								/*layout.Rigid(
									drawHeaders("Данные KitVending:", th),
								),
								layout.Rigid(
									ui.drawFields("ID компании", "1,2,3,4,5,6,7,8,9,0", key.HintNumeric, &ui.Fields[0], th)),
								layout.Rigid(
									ui.drawFields("Логин", "", key.HintAny, &ui.Fields[1], th)),
								layout.Rigid(
									ui.drawFields("Пароль", "", key.HintAny, &ui.Fields[2], th)),
								layout.Rigid(
									drawHeaders("Данные Мой Налог:", th),
								),
								layout.Rigid(
									ui.drawFields("ИНН", "1,2,3,4,5,6,7,8,9,0", key.HintNumeric, &ui.Fields[3], th)),
								layout.Rigid(
									ui.drawFields("Пароль", "", key.HintAny, &ui.Fields[4], th)),*/
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return ui.pairOfButtons(25, ui.okModalBtn, ui.faqBtn, "авторизация", "помощь", th, gtx)
									},
								),
								layout.Rigid(
									func(gtx layout.Context) layout.Dimensions {
										return layout.Inset{Top: 15}.Layout(gtx,
											func(gtx layout.Context) layout.Dimensions {
												return layout.Center.Layout(gtx,
													func(gtx layout.Context) layout.Dimensions {
														var str string
														if ui.emptyFields {
															str = "заполнены не все поля"
														} else {
															str = ""
														}
														return material.Label(th, 16, str).Layout(gtx)
													})
											})
									}),
							)
						})
					}))
			}
		},
	)
}

func (ui *UI) modalDescription(gtx layout.Context) {
	ui.ModalFaq.Widget = func(gtx layout.Context, th *material.Theme, anim *component.VisibilityAnimation) layout.Dimensions {
		sz, off := popSize(9, 5, gtx)
		defer op.Offset(off).Push(gtx.Ops).Pop()
		gtx.Constraints = layout.Exact(sz)
		rect := clip.UniformRRect(image.Rectangle{Max: sz}, 12)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, rect.Op(gtx.Ops))
		return layout.Center.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							return material.Body1(th, authDescribe).Layout(gtx)
						},
					),
				)
			},
		)
	}
	ui.ModalFaq.FinalAlpha = 245
	ui.ModalFaq.Appear(gtx.Now)
}

func drawHeaders(str string, th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: 15, Bottom: 7}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, 16, str)
						label.Font.Weight = 300
						return label.Layout(gtx)
					})
			})

	}
}

func (ui *UI) drawFields(s, n string, typ key.InputHint, f *widget.Editor, th *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: 6, Bottom: 6}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lblDims := material.Label(th, 14, s).Layout(gtx)
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Flexed(0.4,
							func(gtx layout.Context) layout.Dimensions {
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return lblDims
								})
							}),
						layout.Flexed(0.6,
							func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max.X = gtx.Constraints.Max.X * 84 / 100
								//gtx.Constraints.Min.Y = 80
								gtx.Constraints.Min.Y = lblDims.Size.Y + 10
								fl := material.Editor(th, f, "")
								fl.Editor.Alignment = text.Middle
								fl.Editor.Filter = n
								fl.Editor.InputHint = typ
								return widget.Border{
									Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
									CornerRadius: unit.Dp(4),
									Width:        unit.Dp(gtx.Dp(0.5)),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Left: 5, Right: 5, Top: 5}.Layout(gtx,
										func(gtx layout.Context) layout.Dimensions {
											return fl.Layout(gtx)
										})
								})
							}),
					)
				})
			},
		)
	}
}
