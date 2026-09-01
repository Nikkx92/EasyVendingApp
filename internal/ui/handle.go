package ui

import (
	"context"
	"easyVending/internal/domain"
	"easyVending/internal/errs"
	"easyVending/internal/storage"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
)

type SalesData struct {
	Old []domain.Sale
	New []domain.Sale
}

func (ui *UI) isDataEmpty() bool {
	for i := range ui.Fields {
		if ui.Fields[i].Text() == "" {
			ui.emptyFields = true
			break
		} else {
			ui.emptyFields = false
		}
	}
	return ui.emptyFields
}

func getTimeZone() *domain.Time {
	now := time.Now()
	zone, offset := now.Zone()

	t := &domain.Time{
		Zone:   zone,
		Offset: int32(offset),
	}

	return t
}

func (ui *UI) UpdateStateApp() {
	ui.a.StateApp.IsAuth = ui.isAuth
	ui.a.StateApp.AutoMode = ui.AutoSwt.Value
	ui.a.StateApp.CompanyId = ui.Fields[0].Text()
	ui.a.StateApp.LoginKit = ui.Fields[1].Text()
	ui.a.StateApp.PassKit = ui.Fields[2].Text()
	ui.a.StateApp.Inn = ui.Fields[3].Text()
	ui.a.StateApp.PassFns = ui.Fields[4].Text()
}

func (ui *UI) initAutoModeData() *domain.AutoModeData {
	return &domain.AutoModeData{
		CompanyID:    ui.a.StateApp.CompanyId,
		UserLogin:    ui.a.StateApp.LoginKit,
		PasswordKit:  ui.a.StateApp.PassKit,
		INN:          ui.a.StateApp.Inn,
		PasswordFns:  ui.a.StateApp.PassFns,
		DeviceInfo:   ui.a.Device,
		RefreshToken: ui.a.StateApp.RefreshTokenFns,
		Token:        ui.a.StateApp.TokenFns,
		Time:         getTimeZone(),
	}
}

func convertDateTime(input string) string {
	if input == "" {
		log.Println("empty string")
		return ""
	}

	// Парсим исходную дату
	t, err := time.Parse("02.01.2006", input)
	if err != nil {
		log.Println(err)
		return ""
	}

	// Форматируем в целевой формат
	output := t.Format("2006-01-02")
	return output
}

func (ui *UI) HandleEvents(gtx layout.Context) {
	if ui.requestKitBtn.Clicked(gtx) {
		if !ui.isAuth {
			ui.modalResponse(gtx, nil, "вы не авторизованы")
		} else if !ui.a.GRPC.GetConnStatus() && ui.a.StateApp.AutoMode {
			ui.modalResponse(gtx, nil, errs.NoConnect.Error())
		} else {
			if ui.AutoSwt.Value {
				ui.modalResponse(gtx, nil, "Отправка продаж невозможна т.к. включена Автоматическая отправка")
			} else {
				if ui.isLoad {
					return
				}
				ui.isLoad = true

				go func() {
					defer func() { ui.isLoad = false }()

					selectedDate := ui.datePicker.Selected.Format("2006-01-02")
					dateDay := storage.GetSessionByDate(ui.a.DB, selectedDate)

					var dateForKit string
					endTime := "23:59:59"

					var oldDrinks []domain.Sale

					if dateDay != nil {
						timeEndDB := dateDay.TimeEnd

						if timeEndDB == endTime {
							ui.modalResponse(gtx, nil, "нет новых продаж")
							return
						}

						parsedTime, err := time.Parse("15:04:05", timeEndDB)
						if err != nil {
							fmt.Println("Ошибка парсинга:", err)
							return
						}
						incTimeEndDB := parsedTime.Add(1 * time.Second).Format("15:04:05")

						if incTimeEndDB == endTime {
							ui.modalResponse(gtx, nil, "нет новых продаж")
							return
						}

						dateForKit = selectedDate + " " + incTimeEndDB + "--" + selectedDate + " " + endTime

						for i := range dateDay.Drinks {
							oldDrinks = append(oldDrinks, domain.Sale{
								DateTime:  selectedDate,
								GoodsName: dateDay.Drinks[i],
							})
						}

					} else {
						dateForKit = selectedDate + " 00:00:00--" + selectedDate + " " + endTime
					}

					sales, err := ui.a.HTTP.GetDataKitVending(
						ui.a.StateApp.CompanyId,
						ui.a.StateApp.LoginKit,
						ui.a.StateApp.PassKit,
						dateForKit,
					)

					if err != nil {
						ui.modalResponse(gtx, nil, err.Error())
						return
					}

					if sales == nil {
						ui.modalResponse(gtx, nil, "нет новых продаж")
						return
					}

					for i := range sales {
						sales[i].DateTime = convertDateTime(sales[i].DateTime)
					}

					ui.tempDrinks <- SalesData{
						Old: oldDrinks,
						New: sales,
					}

					stub := make([]string, len(sales))
					ui.modalResponse(gtx, stub, "")

					go func() {
						for {
							if !ui.ModalResp.Visible() {
								for len(ui.tempDrinks) > 0 {
									<-ui.tempDrinks
								}
								return
							}
						}
					}()
				}()
			}
		}

	}

	if ui.exitLogin.Clicked(gtx) {
		ui.isAuth = false
		ui.AutoSwt.Value = false
		for i := range ui.Fields {
			ui.Fields[i].SetText("")
		}
		ui.UpdateStateApp()
		storage.SaveState(ui.a.StateApp)
	}

	if ui.deleteData.Clicked(gtx) {
		//date := strings.Split(dateValidator(ui.datePicker.Selected.Format("02.01.2006")), " ")
		/*session := &storage.DrinkSession{
			Date: ui.datePicker.Selected.Format("02.01.2006"),
			Details: storage.Details{
				TimeEnd: "02:00:00",
				Drinks: []string{
					"coffee",
				},
			},
		}

		storage.SaveSession(ui.a.DB, session)*/
		/*fullDate := strings.Split(dateValidator(ui.datePicker.Selected.Format("02.01.2006")), " ")
		dateDay := storage.GetSessionByDate(ui.a.DB, fullDate[0])
		fmt.Println(dateDay)*/
	}

	if ui.isLoad {
		ui.modalResponse(gtx, nil, "")
	}

	if ui.okModalBtn.Clicked(gtx) {
		if ui.isDataEmpty() {
			ui.emptyFields = true
		} else {
			go func() {
				ui.isLoad = true
				defer func() { ui.isLoad = false }()

				/*err := ui.a.HTTP.LoginKitVending(
					ui.Fields[0].Text(),
					ui.Fields[1].Text(),
					ui.Fields[2].Text(),
				)
				if err != nil {
					ui.modalResponse(gtx, nil, err.Error())
					return
				}

				ui.a.StateApp.RefreshTokenFns, ui.a.StateApp.TokenFns, err = ui.a.HTTP.LoginFns(
					ui.Fields[3].Text(),
					ui.Fields[4].Text(),
					ui.a.Device,
				)
				if err != nil {
					ui.modalResponse(gtx, nil, err.Error())
					return
				}*/

				ui.isAuth = true

				ui.UpdateStateApp()
				storage.SaveState(ui.a.StateApp)
			}()
		}
	}

	if ui.sendToFnsBtn.Clicked(gtx) {
		if ui.isLoad {
			return
		}
		ui.isLoad = true

		go func() {
			defer func() { ui.isLoad = false }()

			drinks := <-ui.tempDrinks

			const maxRetries = 2
		outer:
			for attempt := 0; attempt <= maxRetries; attempt++ {

				err := ui.a.HTTP.SendSalesToFns(
					ui.a.StateApp,
					drinks.New,
				)

				if err != nil {

					if errors.Is(err, errs.FnsTokenErr) {
						token, err := ui.a.HTTP.GetToken(ui.a.StateApp.RefreshTokenFns, ui.a.Device)

						if err != nil {
							if errors.Is(err, errs.FnsRefreshTokenErr) {
								ui.a.StateApp.RefreshTokenFns, ui.a.StateApp.TokenFns, err = ui.a.HTTP.LoginFns(
									ui.Fields[3].Text(),
									ui.Fields[4].Text(),
									ui.a.Device,
								)
								if err != nil {
									ui.modalResponse(gtx, nil, err.Error())
									return
								}
								continue outer
							}

							ui.modalResponse(gtx, nil, err.Error())
							return
						}

						ui.a.StateApp.TokenFns = token
						continue outer
					}

					ui.modalResponse(gtx, nil, err.Error())
					return
				}

				if attempt > 0 {
					storage.SaveState(ui.a.StateApp)
				}

				allDrinks := append(drinks.Old, drinks.New...)
				var session storage.DrinkSession

				for i := range allDrinks {
					session.Details.Drinks = append(session.Details.Drinks, allDrinks[i].GoodsName)
				}

				session.Date = allDrinks[0].DateTime

				if time.Now().Format("2006-01-02") != allDrinks[0].DateTime {
					session.Details.TimeEnd = "23:59:59"
				} else {
					session.Details.TimeEnd = time.Now().Format("15:04:05")
				}

				storage.SaveSession(ui.a.DB, &session)

				ui.modalResponse(gtx, nil, "успешная отправка")
				break outer
			}
		}()
	}

	if ui.cancelBtn.Clicked(gtx) {
		<-ui.tempDrinks
		ui.ModalResp.Disappear(time.Now())
	}

	for i := range ui.clkLabel {
		if ui.clkLabel[i].Clicked(gtx) {
			counts := make(map[string]int)
			for _, drink := range ui.sales[i].Details.Drinks {
				counts[drink]++
			}

			items := make([]string, 0, len(counts))
			for drink, count := range counts {
				items = append(items, fmt.Sprintf("%s : %d", drink, count))
			}

			ui.modalDetail(gtx, items, ui.sales[i].Date)
		}
	}

	if ui.faqBtn.Clicked(gtx) {
		ui.modalDescription(gtx)
	}

	if ui.AutoSwt.Update(gtx) {
		if !ui.isAuth {
			ui.AutoSwt.Value = false
			ui.modalResponse(gtx, nil, "вы не авторизованы")
		} else {
			ui.isLoad = true

			go func() {
				defer func() { ui.isLoad = false }()

				var result string
				var err error

				if ui.AutoSwt.Value {
					result, err = ui.a.GRPC.Start(context.Background(), ui.initAutoModeData())

					if err != nil {
						ui.AutoSwt.Value = false
						result = errs.NoConnect.Error()
					} else {
						ui.AutoSwt.Value = true
					}

				} else {
					result, err = ui.a.GRPC.Stop(context.Background(), ui.a.StateApp.Inn, ui.a.StateApp.LoginKit, ui.a.StateApp.CompanyId)
					if err != nil {
						ui.AutoSwt.Value = true
						result = errs.NoConnect.Error()
					}
				}

				ui.modalResponse(gtx, nil, result)
				ui.UpdateStateApp()
				storage.SaveState(ui.a.StateApp)
			}()
		}
	}

	if ui.authBtn.Clicked(gtx) {
		if ui.currentScreen != AuthScreen {

			if ui.currentScreen > AuthScreen {
				ui.slider.PushRight()
			} else {
				ui.slider.PushLeft()
			}
			ui.currentScreen = AuthScreen

		}
	}

	if ui.mainBtn.Clicked(gtx) {
		if ui.currentScreen != MainScreen {
			ui.slider.PushRight()
			ui.currentScreen = MainScreen
		}
	}

	if ui.salesBtn.Clicked(gtx) {
		if ui.currentScreen != HistoryScreen {
			ui.slider.PushLeft()
			ui.currentScreen = HistoryScreen
		}

		if ui.isAuth {
			ui.isLoad = true

			go func() {
				defer func() { ui.isLoad = false }()

				ui.getSales()

			}()
		}

	}

	/*for _, ev := range ui.bar.Events(gtx) {
		switch ev := ev.(type) {
		case component.AppBarOverflowActionClicked:
			switch ev.Tag {
			case "main":
				ui.currentScreen = MainScreen
			case "auth":
				ui.currentScreen = AuthScreen
			case "history":
				ui.isLoad = true

				go func() {
					defer func() { ui.isLoad = false }()

					ui.getSales()

				}()

				ui.currentScreen = HistoryScreen
			case "log":
				ui.currentScreen = LogScreen
			case "about":
				ui.isLoad = true
				go func() {
					ui.workingAuto = nil

					m, message := ui.a.GRPC.WorkingAutoMode(context.Background(), ui.singleRequestInit())
					if message == "" {
						for i := range m {
							ui.workingAuto = append(ui.workingAuto, i)
						}
					} else {
						ui.modalResponse(gtx, nil, message)
					}
					ui.isLoad = false
				}()
				ui.currentScreen = AboutScreen
			}
		}
	}*/

}

func (ui *UI) getSales() {
	salesDB, err := storage.GetSales(ui.a.DB, 30)
	if err != nil {
		fmt.Println(err)
	}

	allSales := make(map[string]*storage.Details, len(salesDB))
	stubDate := "2026-01-01"

	if salesDB == nil {
		salesDB = []storage.DrinkSession{
			{
				Date: stubDate,
				Details: storage.Details{
					TimeEnd: "00:00:00",
				},
			},
		}
	} else {
		for i := range salesDB {
			details, ok := allSales[salesDB[i].Date]
			if !ok {
				details = &storage.Details{}
				allSales[salesDB[i].Date] = details
			}

			details.TimeEnd = salesDB[i].Details.TimeEnd
			details.Drinks = append(details.Drinks, salesDB[i].Details.Drinks...)
		}
	}

	if ui.AutoSwt.Value {

		salesServer, err := ui.a.GRPC.Sales(context.Background(),
			ui.a.StateApp.Inn,
			ui.a.StateApp.LoginKit,
			ui.a.StateApp.CompanyId,
			salesDB[0].Date,
			salesDB[0].Details.TimeEnd,
		)
		if err != nil {
			log.Println(err)
		}

		for i := 1; i <= len(salesServer); i++ {
			sep := strings.Split(salesServer[len(salesServer)-i].DateTime, " ")

			details, ok := allSales[sep[0]]
			if !ok {
				details = &storage.Details{}
				allSales[sep[0]] = details
			}

			details.TimeEnd = sep[1]
			details.Drinks = append(details.Drinks, salesServer[len(salesServer)-i].GoodsName)
		}

	}

	var sales []storage.DrinkSession

	for i := range allSales {
		session := storage.DrinkSession{
			Date: i,
			Details: storage.Details{
				TimeEnd: allSales[i].TimeEnd,
				Drinks:  allSales[i].Drinks,
			},
		}
		sales = append(sales, session)
		storage.SaveSession(ui.a.DB, &session)
	}

	sort.Slice(sales, func(i, j int) bool {
		return sales[j].Date < sales[i].Date
	})

	ui.clkLabel = make([]widget.Clickable, len(sales))
	ui.sales = make([]storage.DrinkSession, len(sales))
	ui.sales = sales
}
