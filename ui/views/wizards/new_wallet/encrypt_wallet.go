package new_wallet

import (
	"errors"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/resources/iconres"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

type EncryptWalletPane struct {
	*base.AbstractWizardPane

	wizardState *WizardState

	canvas               fyne.CanvasObject
	passwordEntry        *widget.Entry
	passwordStrengthIcon *widget.Icon
}

func NewEncryptWalletPane() base.WizardPane {
	return &EncryptWalletPane{
		AbstractWizardPane: base.NewAbstractWizardPane(
			i18n.T("WZ.NewWallet.EncryptWallet.Title"),
		),
	}
}

// Initialize implements base.WizardPane.
func (p *EncryptWalletPane) Initialize() (fyne.CanvasObject, error) {
	p.passwordEntry = widget.NewPasswordEntry()
	p.passwordEntry.OnChanged = func(string) {
		p.checkPasswordEntropy()
		p.AbstractWizardPane.NotifyOnChanged()
	}

	p.passwordStrengthIcon = widget.NewIcon(iconres.IconSignallightgrayPng)

	canvas := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel(i18n.T("WZ.NewWallet.EncryptWallet.Info")),
		container.New(
			layout.NewFormLayout(),
			widget.NewLabel(i18n.T("FrmLbl.Password")),
			container.NewBorder(
				nil, nil,
				nil, p.passwordStrengthIcon,
				p.passwordEntry,
			),
		),
	)
	p.canvas = canvas
	return canvas, nil
}

// IsValid implements base.WizardPane.
func (p *EncryptWalletPane) IsValid() bool {
	return p.canvas != nil
}

// OnHide implements base.WizardPane.
func (*EncryptWalletPane) OnHide() {
}

// OnShow implements base.WizardPane.
func (p *EncryptWalletPane) OnShow() {
	p.passwordEntry.Text = ""
	p.passwordEntry.Refresh()
}

// ResetState implements base.WizardPane.
func (p *EncryptWalletPane) ResetState() {
	p.wizardState.Password = ""
}

// SetState implements base.WizardPane.
func (p *EncryptWalletPane) SetState(state base.WizardState) error {
	var ok bool
	if p.wizardState, ok = state.(*WizardState); !ok {
		return errors.New(i18n.T("Err.ConvertWizardState"))
	}
	return nil
}

// CanTransitionTo implements base.WizardPane.
func (p *EncryptWalletPane) CanTransitionTo(state base.WizardState) bool {
	var ok bool
	var wizardState *WizardState
	if wizardState, ok = state.(*WizardState); !ok {
		return false
	}
	return wizardState.BackupOK
}

// OnBeforeNext implements base.WizardPane.
func (*EncryptWalletPane) OnBeforeNext() bool {
	return true
}

// OnBeforePrevious implements base.WizardPane.
func (*EncryptWalletPane) OnBeforePrevious() bool {
	return true
}

func (*EncryptWalletPane) OnBeforeCancel() bool {
	return true
}

func (p *EncryptWalletPane) OnBeforeFinish() bool {
	p.wizardState.Password = p.passwordEntry.Text
	return true
}

func (p *EncryptWalletPane) checkPasswordEntropy() {
	pwd := p.passwordEntry.Text
	if pwd == "" {
		p.passwordStrengthIcon.SetResource(
			iconres.IconSignallightgrayPng,
		)
		p.passwordStrengthIcon.Refresh()
		return
	}

	pwdStrength := solstice.PasswordStrengthFromEntropy(passwordvalidator.GetEntropy(pwd))
	resource := iconres.IconSignallightgrayPng
	switch pwdStrength {
	case solstice.PasswordWeak:
		resource = iconres.IconSignallightredPng
	case solstice.PasswordFair:
		resource = iconres.IconSignallightorangePng
	case solstice.PasswordStrong:
		resource = iconres.IconSignallightgreenPng
	}
	p.passwordStrengthIcon.SetResource(resource)
	p.passwordStrengthIcon.Refresh()
}
