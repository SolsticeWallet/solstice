package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/views/wizards"
)

type MainView struct {
	*base.AbstractTabbedView
}

// NewMainView returns a new MainView.
//
// Returns a base.View.
func NewMainView() base.View {
	return &MainView{
		AbstractTabbedView: base.NewAbstractTabbedView(
			"",
			container.TabLocationTop,
		),
	}
}

// Initialize initializes the MainView.
//
// It returns a fyne.CanvasObject and an error.
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
	wizard, err := wizards.ShowWizard(
		wizards.CreateWalletWizardView,
		func() {},
		func() {},
		fyne.NewSize(640, 480),
	)
	_ = wizard
	if err != nil {
		base.Logger.Error(
			i18n.T("Err.NewWallet"),
			i18n.T("Err.Arg.Error"), err)
		return
	}
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
		i18n.T("MM.File"),
		func() *fyne.MenuItem {
			item := fyne.NewMenuItem(
				i18n.T("MM.F.NewWallet"),
				v.onMenuNewWallet)
			item.Shortcut = base.NewShortcut(
				i18n.T("MM.F.NewWallet"),
				fyne.KeyN,
				base.StdModifier)
			return item
		}(),
		func() *fyne.MenuItem {
			item := fyne.NewMenuItem(
				i18n.T("MM.F.OpenWallet"),
				v.onMenuOpenWallet)
			item.Shortcut = base.NewShortcut(
				i18n.T("MM.F.OpenWallet"),
				fyne.KeyO,
				base.StdModifier,
			)
			return item
		}(),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(i18n.T("MM.F.ImportWallet"), v.onMenuImportWallet),
		fyne.NewMenuItem(i18n.T("MM.F.ExportWallet"), v.onMenuExportWallet),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem(i18n.T("MM.F.RenameWallet"), v.onMenuRenameWallet),
		fyne.NewMenuItem(i18n.T("MM.F.DeleteWallet"), v.onMenuDeleteWallet),
		func() *fyne.MenuItem {
			item := fyne.NewMenuItem(
				i18n.T("MM.F.CloseTab"),
				v.onMenuCloseTab)
			item.Shortcut = base.NewShortcut(
				i18n.T("MM.F.CloseTab"),
				fyne.KeyW,
				base.StdModifier,
			)
			return item
		}(),
	)
}
