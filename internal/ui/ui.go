package ui

import (
	"easyVending/internal/app"
	"easyVending/internal/datepicker"
	"easyVending/internal/domain"
	"easyVending/internal/slider"
	"easyVending/internal/storage"
	"fmt"
	"image"
	"image/color"
	"os"
	"sync"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Screen int

const (
	MainScreen Screen = iota
	AuthScreen
	HistoryScreen
	AboutScreen
	LogScreen

	Title        = "EasyVending"
	authDescribe = `  CompanyId: идентификатор компании из личного кабинета.

  UserLogin: логин пользователя с правами "API".

  Password: пароль пользователя с правами "API".

  INN: ИНН индивидуального предпринимателя на НПД.

  Password: пароль от личного кабинета налогоплательщика НПД.`
)

type UI struct {
	mu               sync.Mutex
	hasStatusData    bool
	statusData       domain.Result
	a                *app.App
	datePicker       *datepicker.DatePicker
	currentScreen    Screen
	authBtn          *widget.Clickable
	mainBtn          *widget.Clickable
	salesBtn         *widget.Clickable
	AutoSwt          *widget.Bool
	requestKitBtn    *widget.Clickable
	okModalBtn       *widget.Clickable
	faqBtn           *widget.Clickable
	click            *widget.Clickable
	sendToFnsBtn     *widget.Clickable
	cancelBtn        *widget.Clickable
	moreBtn          *widget.Clickable
	exitLogin        *widget.Clickable
	deleteData       *widget.Clickable
	prevGTM          *widget.Clickable
	nextGTM          *widget.Clickable
	clkLabel         []widget.Clickable
	DateField        *component.TextField
	ModalResp        *component.ModalLayer
	ModalFaq         *component.ModalLayer
	ModalForBar      *component.ModalLayer
	ModalDetails     *component.ModalLayer
	Fields           *[5]widget.Editor
	dateMass         []string
	temporaryDataKit []string
	workingAuto      []string
	jwt              string
	deviceId         string
	emptyFields      bool
	dateInvalid      bool
	isLoad           bool
	isAuth           bool
	isPaid           bool
	authIcon         widget.Icon
	mainIcon         widget.Icon
	salesIcon        widget.Icon
	bar              *component.AppBar
	co               []component.OverflowAction
	mapDrinks        map[string]map[string]int32
	list             widget.List
	deviceSelectable widget.Selectable
	tempDrinks       chan SalesData
	sales            []storage.DrinkSession
	alreadyGetSales  bool
	slider           slider.Slider
	emptySalesIcon   paint.ImageOp
}

func NewUI(a *app.App) *UI {
	m := component.NewModal()
	b := component.NewAppBar(m)
	b.Title = Title
	auth, _ := widget.NewIcon(icons.ActionAccountBox)
	sales, _ := widget.NewIcon(icons.ActionAssessment)
	main, _ := widget.NewIcon(icons.NavigationArrowBack)
	img, _ := loadImage()
	return &UI{
		a:             a,
		datePicker:    datepicker.New(),
		authBtn:       new(widget.Clickable),
		mainBtn:       new(widget.Clickable),
		salesBtn:      new(widget.Clickable),
		AutoSwt:       new(widget.Bool),
		requestKitBtn: new(widget.Clickable),
		okModalBtn:    new(widget.Clickable),
		faqBtn:        new(widget.Clickable),
		click:         new(widget.Clickable),
		sendToFnsBtn:  new(widget.Clickable),
		cancelBtn:     new(widget.Clickable),
		moreBtn:       new(widget.Clickable),
		exitLogin:     new(widget.Clickable),
		deleteData:    new(widget.Clickable),
		prevGTM:       new(widget.Clickable),
		ModalResp:     component.NewModal(),
		ModalFaq:      component.NewModal(),
		ModalDetails:  component.NewModal(),
		DateField:     new(component.TextField),
		Fields:        new([5]widget.Editor),
		emptyFields:   false,
		dateInvalid:   false,
		isLoad:        false,
		ModalForBar:   m,
		bar:           b,
		authIcon:      *auth,
		mainIcon:      *main,
		salesIcon:     *sales,
		mapDrinks:     make(map[string]map[string]int32),
		workingAuto:   make([]string, 0),
		dateMass:      make([]string, 0),
		co: []component.OverflowAction{
			{Name: "История", Tag: "history"},
			{Name: "Лог", Tag: "log"},
			{Name: "О программе", Tag: "about"},
		},
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
		currentScreen:  MainScreen,
		tempDrinks:     make(chan SalesData, 1),
		emptySalesIcon: img,
	}
}

func loadImage() (paint.ImageOp, error) {
	f, err := os.Open("C:/Users/user/GolandProjects/easyVending/internal/chart-histogram.png")
	if err != nil {
		return paint.ImageOp{}, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return paint.ImageOp{}, err
	}
	return paint.NewImageOp(img), nil
}

func (ui *UI) SendStateToUI(s storage.StateApp) {
	ui.AutoSwt.Value = s.AutoMode
	ui.isPaid = s.IsPaid
	ui.isAuth = s.IsAuth
	ui.Fields[0].SetText(s.CompanyId)
	ui.Fields[1].SetText(s.LoginKit)
	ui.Fields[2].SetText(s.PassKit)
	ui.Fields[3].SetText(s.Inn)
	ui.Fields[4].SetText(s.PassFns)
	ui.deviceId = s.DeviceId
}

func (ui *UI) AppBar() {
	actions := []component.AppBarAction{
		component.SimpleIconAction(ui.mainBtn, &ui.mainIcon,
			component.OverflowAction{
				Name: "sent screen",
				Tag:  "main",
			}),
		component.SimpleIconAction(ui.authBtn, &ui.authIcon,
			component.OverflowAction{
				Name: "auth",
				Tag:  "auth",
			}),
		component.SimpleIconAction(ui.salesBtn, &ui.salesIcon,
			component.OverflowAction{
				Name: "sales",
				Tag:  "sales",
			}),
	}
	ui.bar.SetActions(actions, nil)

}

func (ui *UI) InvokeModals(gtx layout.Context, th *material.Theme) {
	ui.ModalForBar.Layout(gtx, th)
	ui.ModalResp.Layout(gtx, th)
	ui.ModalFaq.Layout(gtx, th)
	ui.ModalDetails.Layout(gtx, th)
}

func (ui *UI) Draw(gtx layout.Context, th *material.Theme) layout.Dimensions {

	paint.LinearGradientOp{
		Stop1:  f32.Pt(float32(gtx.Constraints.Min.X), 0),
		Stop2:  f32.Pt(float32(unit.Dp(gtx.Constraints.Max.X)), float32(unit.Dp(gtx.Dp(500)))),
		Color1: color.NRGBA{R: 74, G: 169, B: 255, A: 150},
		Color2: color.NRGBA{R: 74, G: 169, B: 255, A: 0},
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ui.bar.Layout(gtx, th, "nav", "over")
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			switch ui.currentScreen {
			case MainScreen:
				ui.slider.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return ui.drawMain(gtx, th)
				})
			case AuthScreen:
				ui.slider.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return ui.drawAuth(gtx, th)
				})
			case HistoryScreen:
				ui.slider.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return ui.drawSales(gtx, th)
				})
			case LogScreen:
				return ui.drawLog(gtx, th)
			case AboutScreen:
				return ui.drawAbout(gtx, th)
			default:
				fmt.Println("default screen")
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: 15}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					if ui.a.StateApp.AutoMode {
						if ui.a.GRPC.GetConnStatus() {
							return layout.Dimensions{}
						}
						b := material.Body1(th, "отсутствует соединение с сервером")
						b.Color = color.NRGBA{R: 255, G: 1, B: 1, A: 235}
						return b.Layout(gtx)
					}
					return layout.Dimensions{}
				})
			})

		}),
	)
}

func popSize(x, y int, gtx layout.Context) (sz, off image.Point) {
	parent := gtx.Constraints.Max
	sz = image.Pt(parent.X*x/10, parent.Y*y/10)
	off = image.Pt((parent.X-sz.X)/2, (parent.Y-sz.Y)/2)
	return sz, off
}

func (ui *UI) pairOfButtons(i int, b1, b2 *widget.Clickable, s1, s2 string, th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(i)}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			c := gtx.Constraints.Max.X
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceSides,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = c / 3
						return material.Button(th, b1, s1).Layout(gtx)
					}),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						return layout.Spacer{Width: unit.Dp(50)}.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = c / 3
						return material.Button(th, b2, s2).Layout(gtx)
					}),
			)
		})

}
