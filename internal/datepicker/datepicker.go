package datepicker

import (
	"easyVending/internal/slider"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var weekdays = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

var months = []string{
	"Январь", "Февраль", "Март", "Апрель",
	"Май", "Июнь", "Июль", "Август",
	"Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
}

type DatePicker struct {
	Visible  bool
	Selected time.Time
	viewing  time.Time

	OpenBtn   widget.Clickable
	prevMonth widget.Clickable
	nextMonth widget.Clickable
	prevYear  widget.Clickable
	nextYear  widget.Clickable
	todayBtn  widget.Clickable
	dayBtns   [42]widget.Clickable
	list      widget.List
	slider    slider.Slider
}

func New() *DatePicker {
	now := time.Now()
	return &DatePicker{
		Selected: now,
		viewing:  now,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}

}

func (dp *DatePicker) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Обработка навигации
	if dp.prevMonth.Clicked(gtx) {
		dp.viewing = dp.viewing.AddDate(0, -1, 0)
		dp.slider.PushRight()
	}
	if dp.nextMonth.Clicked(gtx) {
		dp.viewing = dp.viewing.AddDate(0, 1, 0)
		dp.slider.PushLeft()
	}
	if dp.prevYear.Clicked(gtx) {
		dp.viewing = dp.viewing.AddDate(-1, 0, 0)
	}
	if dp.nextYear.Clicked(gtx) {
		dp.viewing = dp.viewing.AddDate(1, 0, 0)
	}
	if dp.todayBtn.Clicked(gtx) {
		now := time.Now()
		dp.Selected = now
		dp.viewing = now
		dp.Visible = false
	}
	if dp.OpenBtn.Clicked(gtx) {
		dp.Visible = !dp.Visible
		dp.viewing = dp.Selected
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Кнопка с датой
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !dp.Visible {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

					c := component.Surface(th)
					c.CornerRadius = unit.Dp(16)
					c.Elevation = 4
					c.UmbraColor = c.AmbientColor

					return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(30), Bottom: unit.Dp(20), Left: unit.Dp(40), Right: unit.Dp(40)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = 100
							return widget.Border{
								Color:        color.NRGBA{R: 200, G: 200, B: 200, A: 250},
								CornerRadius: unit.Dp(18),
								Width:        unit.Dp(1),
							}.Layout(gtx,
								func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{
										Axis:      layout.Horizontal,
										Alignment: layout.Middle,
									}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											btn := material.Button(th, &dp.OpenBtn, dp.Selected.Format("02.01.2006"))
											btn.Background = color.NRGBA{R: 74, G: 169, B: 255, A: 10}
											btn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
											btn.CornerRadius = 8
											btn.TextSize = unit.Sp(22)

											return btn.Layout(gtx)
										}),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											icon, _ := widget.NewIcon(icons.ActionDateRange)
											return layout.Inset{Right: 5}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return icon.Layout(gtx, color.NRGBA{R: 1, G: 1, B: 1, A: 255})
											})
										}),
									)

								})
						})
					})
				})
			}
			return layout.Dimensions{}
		}),
		// Календарь
		layout.Flexed(1,
			func(gtx layout.Context) layout.Dimensions {
				if !dp.Visible {
					return layout.Dimensions{}
				}
				return dp.layoutCalendar(gtx, th)
			}),
	)
}

func (dp *DatePicker) layoutCalendar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Фон календаря
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				sz := image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
				rr := gtx.Dp(24)
				rect := clip.RRect{Rect: image.Rectangle{Max: sz}, SE: rr, SW: rr, NE: rr, NW: rr}
				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, rect.Op(gtx.Ops))
				// Тень/бордер
				return layout.Dimensions{Size: sz}
			},
			func(gtx layout.Context) layout.Dimensions {
				return material.List(th, &dp.list).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							// Заголовок с годом
							/*layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return dp.layoutYearNav(gtx, th)
							}),*/
							layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
							// Заголовок с месяцем
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return dp.layoutMonthNav(gtx, th)
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.Divider(th).Layout(gtx)
							}),
							// Дни недели
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return dp.layoutWeekdays(gtx, th)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
							// Сетка дней
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return dp.slider.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return dp.layoutDays(gtx, th)
								})
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
							// Кнопка "Сегодня"
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &dp.todayBtn, "Сегодня")
									return btn.Layout(gtx)
								})
							}),
						)
					})
				})
			},
		)
	})
}

func drawButton(gtx layout.Context, button *widget.Clickable, image string, th *material.Theme, top, bottom, left, right int) layout.Dimensions {
	btn := material.Button(th, button, image)
	btn.Background = color.NRGBA{R: 0, G: 0, B: 0, A: 20}
	btn.Inset = layout.Inset{Top: unit.Dp(top), Bottom: unit.Dp(bottom), Left: unit.Dp(left), Right: unit.Dp(right)}
	btn.CornerRadius = unit.Dp(10)
	btn.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	//return widget.Border{
	//Color:        color.NRGBA{R: 200, G: 200, B: 200, A: 250},
	//CornerRadius: unit.Dp(10),
	//Width:        unit.Dp(1),
	//}.Layout(gtx,
	//	func(gtx layout.Context) layout.Dimensions {
	return btn.Layout(gtx)
	//})
}

func (dp *DatePicker) layoutYearNav(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawButton(gtx, &dp.prevYear, "◀◀", th, 8, 8, 10, 12)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.H6(th, fmt.Sprintf("%d", dp.viewing.Year()))
			//lbl.Alignment = 1 // Center
			lbl.Font.Weight = font.Bold
			return layout.Center.Layout(gtx, lbl.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawButton(gtx, &dp.nextYear, "▶▶", th, 8, 8, 12, 10)
		}),
	)
}

func (dp *DatePicker) layoutMonthNav(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawButton(gtx, &dp.prevMonth, "◀", th, 6, 6, 10, 12)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			monthName := months[dp.viewing.Month()-1]
			lbl := material.H6(th, monthName+" "+strconv.Itoa(dp.viewing.Year()))
			lbl.Font.Weight = font.Bold
			return layout.Inset{Left: 14, Right: 14}.Layout(gtx, lbl.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawButton(gtx, &dp.nextMonth, "▶", th, 6, 6, 12, 10)
		}),
	)
}

func (dp *DatePicker) layoutWeekdays(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var children []layout.FlexChild
	for _, day := range weekdays {
		day := day
		children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(th, day)
			lbl.Font.Weight = font.Bold
			lbl.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}

			if day != "Пн" {
				w := gtx.Dp(2)
				h := gtx.Dp(20)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 100, G: 100, B: 100, A: 255},
					clip.Rect{Max: image.Pt(w, h)}.Op(),
				)
			}

			return layout.Center.Layout(gtx, lbl.Layout)
		}))
	}
	return layout.Inset{Top: 20, Bottom: 20}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx, children...)
		})
}

func (dp *DatePicker) layoutDays(gtx layout.Context, th *material.Theme) layout.Dimensions {
	first := time.Date(dp.viewing.Year(), dp.viewing.Month(), 1, 0, 0, 0, 0, dp.viewing.Location())
	startOffset := int(first.Weekday())
	if startOffset == 0 {
		startOffset = 7
	}
	startOffset--

	daysInMonth := time.Date(dp.viewing.Year(), dp.viewing.Month()+1, 0, 0, 0, 0, 0, dp.viewing.Location()).Day()
	today := time.Now()

	var rows []layout.FlexChild
	for week := 0; week < 6; week++ {
		week := week
		rows = append(rows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var cols []layout.FlexChild
			for day := 0; day < 7; day++ {
				day := day
				idx := week*7 + day
				dayNum := idx - startOffset + 1

				cols = append(cols, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(3).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						cellSize := gtx.Constraints.Max.X

						if dayNum < 1 || dayNum > daysInMonth {
							return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 46)}
							//return layout.Dimensions{Size: image.Pt(cellSize, cellSize)}
						}

						if dp.dayBtns[idx].Clicked(gtx) {
							dp.Selected = time.Date(dp.viewing.Year(), dp.viewing.Month(), dayNum, 0, 0, 0, 0, dp.viewing.Location())
							dp.Visible = false
							gtx.Execute(op.InvalidateCmd{})
						}

						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							isSelected := dp.Selected.Day() == dayNum &&
								dp.Selected.Month() == dp.viewing.Month() &&
								dp.Selected.Year() == dp.viewing.Year()

							isToday := today.Day() == dayNum &&
								today.Month() == dp.viewing.Month() &&
								today.Year() == dp.viewing.Year()

							return widget.Border{
								Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 70},
								CornerRadius: 18,
								Width:        1,
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return dp.dayBtns[idx].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									sz := image.Pt(cellSize, cellSize)
									gtx.Constraints = layout.Exact(sz)

									// Фон ячейки
									//rr := cellSize / 2
									rect := clip.RRect{
										Rect: image.Rectangle{Max: sz},
										SE:   gtx.Dp(22),
										SW:   gtx.Dp(22),
										NE:   gtx.Dp(22),
										NW:   gtx.Dp(22),
									}

									var bgColor color.NRGBA
									if isSelected {
										bgColor = color.NRGBA{R: 63, G: 81, B: 181, A: 255}
									} else if isToday {
										bgColor = color.NRGBA{R: 255, G: 193, B: 7, A: 255}
									} else {
										bgColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
									}
									paint.FillShape(gtx.Ops, bgColor, rect.Op(gtx.Ops))

									// Текст
									lbl := material.Body1(th, fmt.Sprintf("%d", dayNum))
									if isSelected {
										lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
									} else if isToday {
										lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
										lbl.Font.Weight = font.Bold
									}
									return layout.Center.Layout(gtx, lbl.Layout)
								})
							})
						})
					})
				}))
			}
			return layout.Flex{}.Layout(gtx, cols...)
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, rows...)
}
