package main

import (
	app2 "easyVending/internal/app"
	"easyVending/internal/ui"
	"log"
	"os"

	"gioui.org/app"
)

func main() {
	go func() {
		application, err := app2.NewApp()
		if err != nil {
			log.Fatal(err)
		}

		u := ui.NewUI(application)

		if err = application.Run(u); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
