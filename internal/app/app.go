package app

import (
	"context"
	"easyVending/internal/domain"
	"easyVending/internal/storage"
	"easyVending/internal/transport/grpc"
	"easyVending/internal/transport/httpNet"
	"image/color"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	_ "time/tzdata"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	bolt "go.etcd.io/bbolt"
)

type View interface {
	HandleEvents(gtx layout.Context)
	Draw(gtx layout.Context, th *material.Theme) layout.Dimensions
	AppBar()
	InvokeModals(gtx layout.Context, th *material.Theme)
	SendStateToUI(s storage.StateApp)
	SendStatusToUI(gtx layout.Context) bool
	InitLoadIcon()
}

type App struct {
	Window   *app.Window
	Th       *material.Theme
	StateApp *storage.StateApp
	Device   *domain.DeviceInfo
	GRPC     *grpc.GRPCClient
	HTTP     *httpNet.ClientHTTP
	DB       *bolt.DB
	status   bool
	DataChan chan domain.Result
}

func initDb() (*bolt.DB, error) {
	dir, err := storage.Path()
	if err != nil {
		return nil, err
	}

	if err = os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dir, "sales.db")

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Создаем bucket при первом запуске
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("sales"))
		return err
	})

	return db, nil
}

func NewApp() (*App, error) {

	cli, err := grpc.NewGRPCClient("192.168.1.88:50052")
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	httpClient := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	db, errDB := initDb()
	if errDB != nil {
		return nil, errDB
	}

	a := &App{
		Window:   new(app.Window),
		Th:       InitTheme(),
		Device:   Device(),
		StateApp: storage.LoadState(),
		GRPC:     cli,
		HTTP:     httpNet.NewClient(httpClient),
		DB:       db,
		DataChan: make(chan domain.Result),
	}

	if a.StateApp.DeviceId == "" {
		a.StateApp.DeviceId = GenerateId()
	}

	a.Device.SourceDeviceID = a.StateApp.DeviceId

	InitAndroidTimezoneProperty()

	return a, nil
}

func InitTheme() *material.Theme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(
		text.WithCollection(gofont.Collection()),
	)
	th.Palette.ContrastBg = color.NRGBA{R: 33, G: 92, B: 150, A: 255}
	return th
}

func Device() *domain.DeviceInfo {
	return &domain.DeviceInfo{
		SourceDeviceID: "",
		SourceType:     runtime.GOOS,
		AppVersion:     "4.5.2-prod",
		MetaDetails: domain.MetaDetails{
			UserAgent: strings.ToUpper(runtime.GOOS) + runtime.GOARCH,
		},
	}
}

func GenerateId() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	b := make([]byte, 21)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (a *App) GetStatus() {
	mess, err := a.GRPC.GetStatus(context.Background(), a.StateApp.Inn, a.StateApp.CompanyId, a.StateApp.LoginKit)
	a.DataChan <- domain.Result{Message: mess, Err: err}
}

func (a *App) Run(v View) error {
	a.Window.Option(app.Size(unit.Dp(425), unit.Dp(750)))
	a.Window.Option(app.Title("Title"))
	var ops op.Ops

	v.SendStateToUI(*a.StateApp)
	v.AppBar()

	if a.StateApp.AutoMode {
		v.InitLoadIcon()
		go a.GetStatus()
	}

	for {
		switch e := a.Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			v.HandleEvents(gtx)

			v.Draw(gtx, a.Th)

			v.InvokeModals(gtx, a.Th)

			if !a.status {
				hasData := v.SendStatusToUI(gtx)
				a.status = hasData
			}

			for {
				keyEvent, ok := gtx.Event(
					key.Filter{
						Name: key.NameBack,
					},
				)
				if !ok {
					break
				}

				switch keyEvent := keyEvent.(type) {
				case key.Event:
					switch keyEvent.Name {
					case key.NameBack:
						return nil
					}
				}
			}

			e.Frame(gtx.Ops)
		}
	}
}
