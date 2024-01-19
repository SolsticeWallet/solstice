//go:generate go run generators/metadata.generator.go cmd/solstice/FyneApp.toml metadata.generated.go solstice

package solstice

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/solsticewallet/solstice-core/log"
)

var App fyne.App
var AppWindow fyne.Window

func InitApp() {
	App = app.New()
	initializeStorage()
	log.CreateDefaultLogger(AppLogDir, GetMetadata().Name)

	AppWindow = App.NewWindow(
		fmt.Sprintf("%s - %s", GetMetadata().Name, GetMetadata().Version))
	AppWindow.Resize(fyne.NewSize(800, 600))
	AppWindow.CenterOnScreen()
}
