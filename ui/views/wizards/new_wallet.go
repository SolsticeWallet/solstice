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
	cancelCallback base.WizardCancelCallback,
	confirmCallback base.WizardConfirmCallback,
) base.WizardView {
	rootPane := new_wallet.NewNetworkAndNamePane()
	wizard := NewNewWalletWizardView(
		parentWindow,
		rootPane,
		cancelCallback,
		confirmCallback,
	)

	restoreOrCreatePane := new_wallet.NewRestoreOrCreatePane()
	wizard.AddStep(rootPane, restoreOrCreatePane, func() bool { return true })

	createSoftPane := new_wallet.NewCreateSoftwareWalletPane()
	wizard.AddStep(restoreOrCreatePane, createSoftPane, func() bool {
		return createSoftPane.CanTransitionTo(wizard.GetState())
	})

	verifySoftPane := new_wallet.NewVerifySoftwareWalletPane()
	wizard.AddStep(createSoftPane, verifySoftPane, func() bool {
		return verifySoftPane.CanTransitionTo(wizard.GetState())
	})

	restoreSoftPane := new_wallet.NewRestoreSoftwareWalletPane()
	wizard.AddStep(restoreOrCreatePane, restoreSoftPane, func() bool {
		return restoreSoftPane.CanTransitionTo(wizard.GetState())
	})

	encryptPane := new_wallet.NewEncryptWalletPane()
	wizard.AddStep(verifySoftPane, encryptPane, func() bool {
		return encryptPane.CanTransitionTo(wizard.GetState())
	})
	wizard.AddStep(restoreSoftPane, encryptPane, func() bool {
		return encryptPane.CanTransitionTo(wizard.GetState())
	})
	return wizard
}

func NewNewWalletWizardView(
	parentWindow fyne.Window,
	rootPane base.WizardPane,
	cancelCallback base.WizardCancelCallback,
	confirmCallback base.WizardConfirmCallback,
) base.WizardView {
	wz := new(NewWalletWizardView)
	wz.wizardState = new_wallet.NewWizardState()
	wz.AbstractWizardView = base.NewAbstractWizardView(
		parentWindow,
		i18n.T("WZ.NewWallet.Title"),
		rootPane,
		func() base.WizardState { return wz.wizardState },
		cancelCallback,
		confirmCallback,
	)
	return wz
}

// GetState implements base.WizardView.
func (w *NewWalletWizardView) GetState() base.WizardState {
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
