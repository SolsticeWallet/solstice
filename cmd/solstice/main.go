package main

import (
	"flag"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/views"
)

var flagLanguage string
var flagDebug bool
var flagDebugI18N bool

func init() {
	flag.StringVar(&flagLanguage, "lang", "", i18n.T("Flag.Language.Usage"))
	flag.BoolVar(&flagDebug, "debug", false, i18n.T("Flag.Debug.Usage"))
	flag.BoolVar(&flagDebugI18N, "debug-i18n", false, i18n.T("Flag.DebugI18N.Usage"))
	flag.Parse()

	if flagLanguage != "" {
		i18n.OverrideLanguage(flagLanguage)
	}
	i18n.SetDebug(flagDebugI18N)
	solstice.InitApp(flagDebug)
}

func main() {
	slog.Info(i18n.T("Info.SolsticeStarting"))

	vws := setupMainWindow()
	defer func() {
		for _, v := range vws {
			v.OnHide()
		}
	}()

	solstice.AppWindow.SetMaster()
	solstice.AppWindow.Show()
	solstice.App.Run()
}

func setupMainWindow() []base.View {
	mainView := views.NewMainView()

	vws := []base.View{
		mainView,
	}

	mainCanvas := container.NewBorder(
		nil, nil,
		nil, nil,
		func() fyne.CanvasObject {
			c, err := mainView.Initialize()
			if err != nil {
				slog.Default().Error(
					i18n.T("Err.InitWindow"),
					i18n.T("Err.Arg.Error"), err)
				panic(err)
			}
			return c
		}(),
	)
	solstice.AppWindow.SetContent(mainCanvas)

	for _, v := range vws {
		v.OnShow()
	}
	return vws
}
