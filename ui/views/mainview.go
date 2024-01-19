package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/ui/base"
)

type MainView struct {
	*base.AbstractTabbedView
}

func NewMainView() base.View {
	return &MainView{
		AbstractTabbedView: base.NewAbstractTabbedView(
			"Main",
			container.TabLocationTop,
		),
	}
}
func (v *MainView) Initialize() (fyne.CanvasObject, error) {
	err := v.configureMainMenu()
	if err != nil {
		return nil, err
	}

	tabs, err := v.AbstractTabbedView.Initialize()
	if err != nil {
		return nil, err
	}
	return tabs, nil
}

// OnHide implements base.View.
func (*MainView) OnHide() {
}

// OnShow implements base.View.
func (*MainView) OnShow() {
}

// Refresh implements base.View.
func (*MainView) Refresh() error {
	return nil
}

// --------------- Main Menu ------------------ //

func (v *MainView) onMenuNewWallet() {

}

func (v *MainView) onMenuOpenWallet() {

}

func (v *MainView) onMenuImportWallet() {

}

func (v *MainView) onMenuExportWallet() {

}

func (v *MainView) onMenuRenameWallet() {

}

func (v *MainView) onMenuDeleteWallet() {

}

func (v *MainView) onMenuCloseTab() {
	v.AbstractTabbedView.RemoveActive()
}

func (v *MainView) configureMainMenu() error {
	menu := solstice.AppWindow.MainMenu()
	if menu == nil {
		menu = fyne.NewMainMenu()
	}

	fileMenu := v.configureFileMenu()

	menu.Items = append(menu.Items,
		fileMenu,
	)
	solstice.AppWindow.SetMainMenu(menu)
	return nil
}

func (v *MainView) configureFileMenu() *fyne.Menu {
	return fyne.NewMenu(
		"File",
		func() *fyne.MenuItem {
			item := fyne.NewMenuItem("New Wallet", v.onMenuNewWallet)
			item.Shortcut = base.NewShortcut(
				"NewWallet",
				fyne.KeyN,
				fyne.KeyModifierSuper,
			)
			return item
		}(),
		func() *fyne.MenuItem {
			item := fyne.NewMenuItem("Open Wallet...", v.onMenuOpenWallet)
			item.Shortcut = base.NewShortcut(
				"OpenWallet",
				fyne.KeyO,
				fyne.KeyModifierSuper,
			)
			return item
		}(),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Import Wallet...", v.onMenuImportWallet),
		fyne.NewMenuItem("Export Wallet...", v.onMenuExportWallet),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Rename Wallet...", v.onMenuRenameWallet),
		fyne.NewMenuItem("Delete Wallet...", v.onMenuDeleteWallet),
		func() *fyne.MenuItem {
			item := fyne.NewMenuItem("Close Tab", v.onMenuCloseTab)
			item.Shortcut = base.NewShortcut(
				"CloseTab",
				fyne.KeyW,
				fyne.KeyModifierSuper,
			)
			return item
		}(),
	)
}
