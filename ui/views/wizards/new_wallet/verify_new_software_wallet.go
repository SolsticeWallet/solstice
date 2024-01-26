package new_wallet

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/custom_widgets"
)

type VerifySoftwareWalletPane struct {
	*base.AbstractWizardPane

	canvas              fyne.CanvasObject
	columnContainers    [4]*fyne.Container
	wordLabels          [24]*widget.Label
	wordEntries         [24]*custom_widgets.Bip39Entry
	passphraseContainer *fyne.Container
	passphraseEntry     *widget.Entry

	backupOk bool

	wizardState *WizardState
}

func NewVerifySoftwareWalletPane() base.WizardPane {
	pane := &VerifySoftwareWalletPane{
		AbstractWizardPane: base.NewAbstractWizardPane(
			i18n.T("WZ.NewWallet.VerifySoftware.Title"),
		),
	}
	return pane
}

// Initialize implements base.WizardPane.
func (p *VerifySoftwareWalletPane) Initialize() (fyne.CanvasObject, error) {
	p.columnContainers = [4]*fyne.Container{
		container.New(layout.NewFormLayout()),
		container.New(layout.NewFormLayout()),
		container.New(layout.NewFormLayout()),
		container.New(layout.NewFormLayout()),
	}

	for i := 0; i < 24; i++ {
		p.wordLabels[i] = widget.NewLabel(fmt.Sprintf("%d:", i+1))
		p.wordLabels[i].Alignment = fyne.TextAlignTrailing
		p.wordEntries[i] = custom_widgets.NewBip39Entry()
		p.wordEntries[i].Validator = func(idx int) fyne.StringValidator {
			return func(v string) error {
				if p.wizardState == nil {
					return errors.New("no state")
				}
				words := strings.Split(p.wizardState.Mnemonic, " ")
				if idx >= len(words) {
					return errors.New("invalid index")
				}
				if strings.ToLower(v) != words[idx] {
					return errors.New("invalid word")
				}
				return nil
			}
		}(i)
		p.wordEntries[i].OnChanged = func(string) {
			p.AbstractWizardPane.NotifyOnChanged()
		}
		p.columnContainers[i/6].Add(p.wordLabels[i])
		p.columnContainers[i/6].Add(p.wordEntries[i])
	}

	mnemonicContainer := container.New(
		layout.NewGridLayout(4),
		p.columnContainers[0],
		p.columnContainers[1],
		p.columnContainers[2],
		p.columnContainers[3],
	)

	topContainer := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel(
			i18n.Tn(
				"WZ.NewWallet.VerifySoftware.Info",
				func() int {
					if p.wizardState == nil {
						return 1
					}
					if p.wizardState.Passphrase == "" {
						return 1
					}
					return 2
				}(),
				"numWords", func() string {
					if p.wizardState == nil {
						return "0"
					}
					return fmt.Sprintf(
						"%d",
						strings.Count(p.wizardState.Mnemonic, " ")+1)
				}(),
			)),
		mnemonicContainer,
	)

	p.passphraseEntry = widget.NewPasswordEntry()
	p.passphraseEntry.Validator = func(v string) error {
		if p.wizardState == nil {
			return errors.New("no state")
		}
		if v != p.wizardState.Passphrase {
			return errors.New("invalid passphrase")
		}
		return nil
	}
	p.passphraseEntry.OnChanged = func(string) {
		p.AbstractWizardPane.NotifyOnChanged()
	}
	p.passphraseContainer = container.New(
		layout.NewFormLayout(),
		widget.NewLabel(i18n.T("FrmLbl.Passphrase")),
		p.passphraseEntry,
	)

	canvas := container.NewBorder(
		topContainer, p.passphraseContainer,
		nil, nil,
	)
	p.canvas = canvas
	return canvas, nil
}

// IsValid implements base.WizardPane.
func (p *VerifySoftwareWalletPane) IsValid() bool {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return false
	}

	p.backupOk = false
	numWords := strings.Count(p.wizardState.Mnemonic, " ") + 1
	for i := 0; i < numWords; i++ {
		if p.wordEntries[i].Validate() != nil {
			return false
		}
	}
	if p.wizardState.Passphrase != "" {
		if p.passphraseEntry.Validate() != nil {
			return false
		}
	}

	p.backupOk = true
	return true
}

// OnHide implements base.WizardPane.
func (*VerifySoftwareWalletPane) OnHide() {
}

// OnShow implements base.WizardPane.
func (p *VerifySoftwareWalletPane) OnShow() {
	if p.wizardState == nil {
		return
	}

	numWords := strings.Count(p.wizardState.Mnemonic, " ") + 1
	p.displayMnemonicEntries(numWords)

	if p.wizardState.Passphrase == "" {
		p.passphraseContainer.Hide()
	} else {
		p.passphraseContainer.Show()
	}
}

// ResetState implements base.WizardPane.
func (p *VerifySoftwareWalletPane) ResetState() {
	p.wizardState.BackupOK = false
}

// SetState implements base.WizardPane.
func (p *VerifySoftwareWalletPane) SetState(state base.WizardState) error {
	var ok bool
	if p.wizardState, ok = state.(*WizardState); !ok {
		return errors.New(i18n.T("Err.ConvertWizardState"))
	}
	return nil
}

// CanTransitionTo implements base.WizardPane.
func (*VerifySoftwareWalletPane) CanTransitionTo(state base.WizardState) bool {
	var ok bool
	var wizardState *WizardState
	if wizardState, ok = state.(*WizardState); !ok {
		return false
	}
	return wizardState.Mnemonic != ""
}

// OnBeforeNext implements base.WizardPane.
func (p *VerifySoftwareWalletPane) OnBeforeNext() bool {
	p.wizardState.BackupOK = p.backupOk
	return true
}

// OnBeforePrevious implements base.WizardPane.
func (p *VerifySoftwareWalletPane) OnBeforePrevious() bool {
	return true
}

func (*VerifySoftwareWalletPane) OnBeforeCancel() bool {
	return true
}

func (*VerifySoftwareWalletPane) OnBeforeFinish() bool {
	return false
}

func (p *VerifySoftwareWalletPane) displayMnemonicEntries(numWords int) {
	p.columnContainers[0].RemoveAll()
	p.columnContainers[1].RemoveAll()
	p.columnContainers[2].RemoveAll()
	p.columnContainers[3].RemoveAll()

	for i := 0; i < numWords; i++ {
		p.wordLabels[i].Show()
		p.wordEntries[i].Show()
	}
	for i := numWords; i < 24; i++ {
		p.wordLabels[i].Hide()
		p.wordEntries[i].Hide()
	}
	for i := 0; i < numWords; i++ {
		p.columnContainers[i/(numWords/4)].Add(p.wordLabels[i])
		p.columnContainers[i/(numWords/4)].Add(p.wordEntries[i])
	}
}
