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

type RestoreSoftwareWalletPane struct {
	*base.AbstractWizardPane

	canvas              fyne.CanvasObject
	mnemonicTypeSelect  *widget.Select
	columnContainers    [4]*fyne.Container
	wordLabels          [24]*widget.Label
	wordEntries         [24]*custom_widgets.Bip39Entry
	passphraseContainer *fyne.Container
	passphraseEntry     *widget.Entry

	wizardState *WizardState
	mnemonic    string
}

func NewRestoreSoftwareWalletPane() base.WizardPane {
	pane := &RestoreSoftwareWalletPane{
		AbstractWizardPane: base.NewAbstractWizardPane(
			i18n.T("WZ.NewWallet.RestoreSoftware.Title"),
		),
	}
	return pane
}

// Initialize implements base.WizardPane.
func (p *RestoreSoftwareWalletPane) Initialize() (fyne.CanvasObject, error) {
	staticInitialize()

	p.mnemonicTypeSelect = widget.NewSelect(
		mnemonicTypes,
		p.onMnemonicTypeChanged,
	)
	p.mnemonicTypeSelect.SetSelectedIndex(0)

	topForm := container.New(
		layout.NewFormLayout(),
		widget.NewLabel(i18n.T("FrmLbl.MnemonicType")),
		p.mnemonicTypeSelect,
	)

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
		topForm,
		widget.NewLabel(i18n.Tn(
			"WZ.NewWallet.RestoreSoftware.Info",
			func() int {
				idx := p.mnemonicTypeSelect.SelectedIndex()
				return idx%2 + 1
			}(),
			"numWords", func() string {
				idx := p.mnemonicTypeSelect.SelectedIndex()
				return fmt.Sprintf("%d", 12+(idx/2)*12)
			}(),
		)),
		mnemonicContainer,
	)

	p.passphraseEntry = widget.NewPasswordEntry()
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
func (p *RestoreSoftwareWalletPane) IsValid() bool {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return false
	}

	idx := p.mnemonicTypeSelect.SelectedIndex()
	end := 12
	if idx > 1 {
		end = 24
	}
	for i := 0; i < end; i++ {
		if p.wordEntries[i].Text == "" {
			return false
		}
	}
	if idx%2 == 1 {
		return p.passphraseEntry.Text != ""
	}
	return true
}

// OnHide implements base.WizardPane.
func (*RestoreSoftwareWalletPane) OnHide() {
}

// OnShow implements base.WizardPane.
func (p *RestoreSoftwareWalletPane) OnShow() {
	if p.wizardState == nil {
		return
	}

	typeIdx := 0
	if p.wizardState.Mnemonic != "" {
		if strings.Count(p.wizardState.Mnemonic, " ")+1 > 12 {
			typeIdx = 2
		}
		if p.wizardState.Passphrase != "" {
			typeIdx += 1
		}
	}
	p.mnemonicTypeSelect.SetSelectedIndex(typeIdx)

	if p.wizardState.Mnemonic != "" {
		p.setMnemonic(p.wizardState.Mnemonic)
		p.passphraseEntry.Text = p.wizardState.Passphrase
		p.passphraseEntry.Refresh()
	}
	p.setUIState()
}

// ResetState implements base.WizardPane.
func (p *RestoreSoftwareWalletPane) ResetState() {
	p.wizardState.Mnemonic = ""
	p.wizardState.Passphrase = ""
	p.wizardState.BackupOK = false
}

// SetState implements base.WizardPane.
func (p *RestoreSoftwareWalletPane) SetState(state base.WizardState) error {
	var ok bool
	if p.wizardState, ok = state.(*WizardState); !ok {
		return errors.New(i18n.T("Err.ConvertWizardState"))
	}
	return nil
}

// CanTransitionTo implements base.WizardPane.
func (p *RestoreSoftwareWalletPane) CanTransitionTo(state base.WizardState) bool {
	var ok bool
	var wizardState *WizardState
	if wizardState, ok = state.(*WizardState); !ok {
		return false
	}
	return wizardState.RestoreOrCreate == RestoreSoftwareWallet
}

// OnBeforeNext implements base.WizardPane.
func (p *RestoreSoftwareWalletPane) OnBeforeNext() bool {
	p.wizardState.Mnemonic = p.mnemonic

	p.wizardState.Passphrase = ""
	if !p.passphraseEntry.Hidden {
		p.wizardState.Passphrase = p.passphraseEntry.Text
	}
	p.wizardState.BackupOK = true
	return true
}

// OnBeforePrevious implements base.WizardPane.
func (*RestoreSoftwareWalletPane) OnBeforePrevious() bool {
	return true
}

func (*RestoreSoftwareWalletPane) OnBeforeCancel() bool {
	return true
}

func (*RestoreSoftwareWalletPane) OnBeforeFinish() bool {
	return false
}

func (p *RestoreSoftwareWalletPane) onMnemonicTypeChanged(mnemonicType string) {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return
	}

	p.clearMnemonic()
	p.passphraseEntry.Text = ""
	p.setUIState()
	p.AbstractWizardPane.NotifyOnChanged()
}

func (p *RestoreSoftwareWalletPane) setUIState() {
	idx := p.mnemonicTypeSelect.SelectedIndex()

	p.columnContainers[0].RemoveAll()
	p.columnContainers[1].RemoveAll()
	p.columnContainers[2].RemoveAll()
	p.columnContainers[3].RemoveAll()

	if idx < 2 {
		for i := 12; i < 24; i++ {
			p.wordLabels[i].Hide()
			p.wordEntries[i].Hide()
		}

		for i := 0; i < 12; i++ {
			p.columnContainers[i/3].Add(p.wordLabels[i])
			p.columnContainers[i/3].Add(p.wordEntries[i])
		}
	} else {
		for i := 12; i < 24; i++ {
			p.wordLabels[i].Show()
			p.wordEntries[i].Show()
		}

		for i := 0; i < 24; i++ {
			p.columnContainers[i/6].Add(p.wordLabels[i])
			p.columnContainers[i/6].Add(p.wordEntries[i])
		}
	}

	if idx%2 == 1 {
		p.passphraseContainer.Show()
	} else {
		p.passphraseContainer.Hide()
	}

	p.canvas.Refresh()
}

func (p *RestoreSoftwareWalletPane) clearMnemonic() {
	for i := 0; i < 24; i++ {
		p.wordEntries[i].SetText("")
		p.wordEntries[i].Refresh()
	}
	p.mnemonic = ""
}

func (p *RestoreSoftwareWalletPane) setMnemonic(mnemonic string) {
	p.mnemonic = mnemonic
	words := strings.Split(mnemonic, " ")
	for i, word := range words {
		p.wordEntries[i].SetText(word)
		p.wordEntries[i].Refresh()
	}
}
