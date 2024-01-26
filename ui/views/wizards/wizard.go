package wizards

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
)

type WizardCreator func(fyne.Window, base.WizardCancelCallback, base.WizardConfirmCallback) base.WizardView

func ShowWizardInDialog(
	creator WizardCreator,
	cancelCallback base.WizardCancelCallback,
	confirmCallback base.WizardConfirmCallback,
	size ...fyne.Size,
) (base.WizardView, error) {
	var dlg dialog.Dialog

	wizard := creator(
		solstice.AppWindow,
		func() {
			dlg.Hide()
			cancelCallback()
		},
		func(result any) {
			dlg.Hide()
			confirmCallback(result)
		},
	)

	canvas, err := wizard.Initialize()
	if err != nil {
		base.Logger.Error(
			i18n.T("Err.InitWizard"),
			i18n.T("Err.Arg.Error"), err)
		return nil, err
	}

	dlg = dialog.NewCustomWithoutButtons(
		wizard.Title(),
		canvas,
		solstice.AppWindow)
	if len(size) == 0 {
		dlg.Resize(fyne.NewSize(400, 300))
	} else {
		dlg.Resize(size[0])
	}
	dlg.SetOnClosed(wizard.OnHide)
	wizard.OnShow()
	dlg.Show()

	return wizard, nil
}

func ShowWizardInWindow(
	creator WizardCreator,
	cancelCallback base.WizardCancelCallback,
	confirmCallback base.WizardConfirmCallback,
	size ...fyne.Size,
) (base.WizardView, error) {
	wnd := solstice.App.NewWindow("")

	wizard := creator(
		wnd,
		func() {
			wnd.Close()
			cancelCallback()
		},
		func(result any) {
			wnd.Close()
			confirmCallback(result)
		},
	)

	canvas, err := wizard.Initialize()
	if err != nil {
		base.Logger.Error(
			i18n.T("Err.InitWizard"),
			i18n.T("Err.Arg.Error"), err)
		return nil, err
	}

	wnd.SetTitle(wizard.Title())
	if len(size) == 0 {
		wnd.Resize(fyne.NewSize(400, 300))
	} else {
		wnd.Resize(size[0])
	}
	wnd.CenterOnScreen()

	wnd.SetContent(canvas)
	wnd.SetOnClosed(wizard.OnHide)
	wizard.OnShow()
	wnd.Show()

	return wizard, nil
}
