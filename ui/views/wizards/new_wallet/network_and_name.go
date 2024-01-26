package new_wallet

import (
	"errors"
	"regexp"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
)

type NetworkAndNamePane struct {
	*base.AbstractWizardPane

	wizardState *WizardState

	canvas          fyne.CanvasObject
	networkSelect   *widget.Select
	walletNameInput *widget.Entry

	walletNameRegexp *regexp.Regexp
}

func NewNetworkAndNamePane() base.WizardPane {
	return &NetworkAndNamePane{
		AbstractWizardPane: base.NewAbstractWizardPane(
			i18n.T("WZ.NewWallet.NetworkName.Title"),
		),
		walletNameRegexp: regexp.MustCompile(`^[A-Za-z]{1}(\w|-){2,}$`),
	}
}

func (p *NetworkAndNamePane) Initialize() (fyne.CanvasObject, error) {
	p.networkSelect = widget.NewSelect(
		solstice.SupportedNetworks,
		func(s string) {
			p.AbstractWizardPane.NotifyOnChanged()
		})
	p.networkSelect.SetSelected(solstice.SupportedNetworks[0])

	p.walletNameInput = widget.NewEntry()
	p.walletNameInput.PlaceHolder = i18n.T("WZ.NewWallet.NetworkName.WalletNamePlaceholder")
	p.walletNameInput.OnChanged = func(string) {
		p.AbstractWizardPane.NotifyOnChanged()
	}
	p.walletNameInput.Validator = func(v string) error {
		if !p.isValidWalletName(v) {
			return errors.New("invalid wallet name")
		}
		return nil
	}

	canvas := container.New(
		layout.NewFormLayout(),
		widget.NewLabel(i18n.T("FrmLbl.Network")),
		p.networkSelect,
		widget.NewLabel(i18n.T("FrmLbl.WalletName")),
		p.walletNameInput)
	p.canvas = canvas
	return canvas, nil
}

func (p *NetworkAndNamePane) IsValid() bool {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return false
	}
	return p.networkSelect.Selected != "" && p.isValidWalletName(p.walletNameInput.Text)
}

func (p *NetworkAndNamePane) OnHide() {

}

func (p *NetworkAndNamePane) OnShow() {
	if p.wizardState == nil {
		return
	}
	if p.wizardState.Network == "" {
		p.wizardState.Network = solstice.SupportedNetworks[0]
	}

	p.networkSelect.SetSelected(p.wizardState.Network)
	p.networkSelect.Refresh()

	p.walletNameInput.Text = p.wizardState.WalletName
	p.walletNameInput.Refresh()
}

func (p *NetworkAndNamePane) ResetState() {
	p.wizardState.Network = solstice.SupportedNetworks[0]
	p.wizardState.WalletName = ""
}

func (p *NetworkAndNamePane) SetState(state base.WizardState) error {
	var ok bool
	if p.wizardState, ok = state.(*WizardState); !ok {
		return errors.New(i18n.T("Err.ConvertWizardState"))
	}
	return nil
}

func (p *NetworkAndNamePane) CanTransitionTo(state base.WizardState) bool {
	return true
}

func (p *NetworkAndNamePane) OnBeforeNext() bool {
	p.wizardState.Network = solstice.SupportedNetworks[p.networkSelect.SelectedIndex()]
	p.wizardState.WalletName = p.walletNameInput.Text
	return true
}

func (p *NetworkAndNamePane) OnBeforePrevious() bool {
	return true
}

func (p *NetworkAndNamePane) OnBeforeCancel() bool {
	return true
}

func (p *NetworkAndNamePane) OnBeforeFinish() bool {
	// We do not allow to finish from here
	return false
}

func (p *NetworkAndNamePane) isValidWalletName(walletName string) bool {
	return p.walletNameRegexp.MatchString(walletName)
}
