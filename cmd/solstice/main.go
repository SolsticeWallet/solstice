package main

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/views"
)

func init() {
	solstice.InitApp()
}

func main() {
	slog.Info("Solstice starting")

	vws := setupMainWindow()
	defer func() {
		for _, v := range vws {
			v.OnHide()
		}
	}()

	solstice.AppWindow.Show()
	solstice.App.Run()
}

func setupMainWindow() []base.View {
	mainView := views.NewMainView()

	vws := []base.View{
		mainView,
	}

	layout := container.NewBorder(
		nil, nil,
		nil, nil,
		func() fyne.CanvasObject {
			c, err := mainView.Initialize()
			if err != nil {
				slog.Default().Error(
					"window initialization failed",
					"error", err)
				panic(err)
			}
			return c
		}(),
	)
	solstice.AppWindow.SetContent(layout)

	for _, v := range vws {
		v.OnShow()
	}
	return vws
}
