package wizards

import (
	"fyne.io/fyne/v2"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/views/wizards/new_wallet"
)

type NewWalletWizardView struct {
	*base.AbstractWizardView

	wizardState *new_wallet.WizardState
}

func CreateWalletWizardView(
	parentWindow fyne.Window,
	cancelCallback func(),
	confirmCallback func(),
) base.WizardView {
	rootPane := new_wallet.NewNetworkAndNamePane()
	wizard := NewNewWalletWizardView(
		parentWindow,
		rootPane,
		cancelCallback,
		confirmCallback,
	)

	pane01 := new_wallet.NewRestoreOrCreatePane()
	wizard.AddStep(rootPane, pane01, func() bool { return true })

	pane02 := new_wallet.NewCreateSoftwareWalletPane()
	wizard.AddStep(pane01, pane02, func() bool {
		return pane02.CanTransitionTo(wizard.GetState())
	})

	pane03 := new_wallet.NewVerifySoftwareWalletPane()
	wizard.AddStep(pane02, pane03, func() bool {
		return pane03.CanTransitionTo(wizard.GetState())
	})
	return wizard
}

func NewNewWalletWizardView(
	parentWindow fyne.Window,
	rootPane base.WizardPane,
	cancelCallback func(),
	confirmCallback func(),
) base.WizardView {
	wz := new(NewWalletWizardView)
	wz.wizardState = new_wallet.NewWizardState()
	wz.AbstractWizardView = base.NewAbstractWizardView(
		parentWindow,
		i18n.T("WZ.NewWallet.Title"),
		rootPane,
		func() any { return wz.wizardState },
		cancelCallback,
		confirmCallback,
	)
	return wz
}

// GetState implements base.WizardView.
func (w *NewWalletWizardView) GetState() any {
	return w.wizardState
}

// OnHide implements base.WizardView.
func (w *NewWalletWizardView) OnHide() {
	w.AbstractWizardView.OnHide()
}

// OnShow implements base.WizardView.
func (w *NewWalletWizardView) OnShow() {
	w.AbstractWizardView.OnShow()
}

// Refresh implements base.WizardView.
func (*NewWalletWizardView) Refresh() error {
	return nil
}
