package new_wallet

import (
	"errors"
	"fmt"
	"strings"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice-core/blockchains/ethereum/utils"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"github.com/solsticewallet/solstice/ui/custom_widgets"
)

var MnemonicTypes []string

type CreateSoftwareWalletPane struct {
	*base.AbstractWizardPane

	canvas              fyne.CanvasObject
	mnemonicTypeSelect  *widget.Select
	columnContainers    [4]*fyne.Container
	wordLabels          [24]*widget.Label
	wordEntries         [24]*custom_widgets.Entry
	passphraseContainer *fyne.Container
	passphraseEntry     *widget.Entry

	wizardState *WizardState
	mnemonic    string
}

func NewCreateSoftwareWalletPane() base.WizardPane {
	pane := &CreateSoftwareWalletPane{
		AbstractWizardPane: base.NewAbstractWizardPane(
			i18n.T("WZ.NewWallet.CreateSoftware.Title"),
		),
	}
	return pane
}

// Initialize implements base.WizardPane.
func (p *CreateSoftwareWalletPane) Initialize() (fyne.CanvasObject, error) {
	p.staticInitialize()

	mnemonicSelectLayout := container.NewBorder(
		nil, nil,
		nil, container.NewHBox(
			widget.NewButtonWithIcon("", theme.QuestionIcon(), func() {
				dlgContent := container.New(
					layout.NewVBoxLayout(),
					widget.NewLabel(
						i18n.T("WZ.NewWallet.CreateSoftware.InfoBip39")),
					widget.NewHyperlink(
						solstice.ExternalLinks["Bip39Info"].String(),
						solstice.ExternalLinks["Bip39Info"]),
				)
				dlg := dialog.NewCustom(
					i18n.T("WZ.NewWallet.CreateSoftware.DlgBip39Info.Title"),
					i18n.T("Lbl.Close"),
					dlgContent,
					p.GetParentWindow(),
				)
				dlg.Show()
			}),
			widget.NewButtonWithIcon(
				i18n.T("Lbl.Generate"),
				theme.ViewRefreshIcon(),
				func() {
					p.generateNewMnemonic()
				}),
		),
	)
	p.mnemonicTypeSelect = widget.NewSelect(
		MnemonicTypes,
		p.onMnemonicTypeChanged,
	)
	p.mnemonicTypeSelect.SetSelectedIndex(0)
	mnemonicSelectLayout.Add(p.mnemonicTypeSelect)

	topForm := container.New(
		layout.NewFormLayout(),
		widget.NewLabel(i18n.T("WZ.NewWallet.CreateSoftware.LblMnemonicType")),
		mnemonicSelectLayout,
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
		p.wordEntries[i] = custom_widgets.NewEntry()
		p.wordEntries[i].SetReadonly(true)
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
		widget.NewLabel(i18n.T("WZ.NewWallet.CreateSoftware.Info")),
		mnemonicContainer,
	)

	p.passphraseEntry = widget.NewPasswordEntry()
	p.passphraseEntry.OnChanged = func(string) {
		p.AbstractWizardPane.NotifyOnChanged()
	}
	p.passphraseContainer = container.New(
		layout.NewFormLayout(),
		widget.NewLabel(i18n.T("WZ.NewWallet.CreateSoftware.LblPassphrase")),
		p.passphraseEntry,
	)

	canvas := container.NewBorder(
		topContainer, p.passphraseContainer,
		nil, nil,
	)
	p.canvas = canvas
	return canvas, nil
}

func (*CreateSoftwareWalletPane) staticInitialize() {
	if len(MnemonicTypes) == 0 {
		MnemonicTypes = []string{
			i18n.T("Lbl.MnemonicTypes.12Words"),
			i18n.T("Lbl.MnemonicTypes.12WordsPwd"),
			i18n.T("Lbl.MnemonicTypes.24Words"),
			i18n.T("Lbl.MnemonicTypes.24WordsPwd"),
		}
	}
}

// IsValid implements base.WizardPane.
func (p *CreateSoftwareWalletPane) IsValid() bool {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return false
	}

	idx := p.mnemonicTypeSelect.SelectedIndex()
	if idx%2 == 1 {
		return p.passphraseEntry.Text != ""
	}

	return true
}

// OnHide implements base.WizardPane.
func (p *CreateSoftwareWalletPane) OnHide() {
}

// OnShow implements base.WizardPane.
func (p *CreateSoftwareWalletPane) OnShow() {
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
func (p *CreateSoftwareWalletPane) ResetState() {
	p.wizardState.Mnemonic = ""
	p.wizardState.Passphrase = ""
}

// SetState implements base.WizardPane.
func (p *CreateSoftwareWalletPane) SetState(state any) error {
	var ok bool
	if p.wizardState, ok = state.(*WizardState); !ok {
		return errors.New(i18n.T("Err.ConvertWizardState"))
	}
	return nil
}

func (p *CreateSoftwareWalletPane) CanTransitionTo(state any) bool {
	var ok bool
	var wizardState *WizardState
	if wizardState, ok = state.(*WizardState); !ok {
		return false
	}
	return wizardState.RestoreOrCreate == CreateSoftwareWallet
}

func (p *CreateSoftwareWalletPane) OnBeforeNext() {
	p.wizardState.Mnemonic = p.mnemonic

	p.wizardState.Passphrase = ""
	if !p.passphraseEntry.Hidden {
		p.wizardState.Passphrase = p.passphraseEntry.Text
	}
}

func (p *CreateSoftwareWalletPane) OnBeforePrevious() {

}

func (p *CreateSoftwareWalletPane) onMnemonicTypeChanged(mnemonicType string) {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return
	}

	p.generateNewMnemonic()
	p.passphraseEntry.Text = ""
	p.setUIState()
	p.AbstractWizardPane.NotifyOnChanged()
}

func (p *CreateSoftwareWalletPane) setUIState() {
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

func (p *CreateSoftwareWalletPane) generateNewMnemonic() {
	idx := p.mnemonicTypeSelect.SelectedIndex()
	bitSize := 128
	if idx > 1 {
		bitSize = 256
	}

	mnemonic, err := utils.NewMnemonic(bitSize)
	//p.wizardState.Mnemonic, err = utils.NewMnemonic(bitSize)
	if err != nil {
		base.Logger.Error(
			i18n.T("Err.MnemonicGeneration"),
			i18n.T("Err.Arg.Error"), err)
		return
	}
	p.setMnemonic(mnemonic)
}

func (p *CreateSoftwareWalletPane) setMnemonic(mnemonic string) {
	p.mnemonic = mnemonic
	words := strings.Split(mnemonic, " ")
	for i, word := range words {
		p.wordEntries[i].SetText(word)
		p.wordEntries[i].Refresh()
	}
}
