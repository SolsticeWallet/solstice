package new_wallet

import (
	"errors"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
	"golang.org/x/exp/slices"
)

var (
	RestoreOrCreateLabels    []string
	RestoreOrCreateLabelInfo []string
)

type RestoreOrCreatePane struct {
	*base.AbstractWizardPane

	wizardState *WizardState

	canvas      fyne.CanvasObject
	actionGroup *widget.RadioGroup
	infoLabel   *widget.Label
}

func NewRestoreOrCreatePane() base.WizardPane {
	pane := &RestoreOrCreatePane{
		AbstractWizardPane: base.NewAbstractWizardPane(
			i18n.T("WZ.NewWallet.RestoreCreate.Title"),
		),
	}
	return pane
}

// Initialize implements base.WizardPane.
func (p *RestoreOrCreatePane) Initialize() (fyne.CanvasObject, error) {
	p.staticInitialize()
	p.actionGroup = widget.NewRadioGroup(
		RestoreOrCreateLabels,
		p.onActionGroupChanged,
	)
	p.infoLabel = widget.NewLabel("")

	canvas := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel(i18n.T("WZ.NewWallet.RestoreCreate.Info")),
		p.actionGroup,
		p.infoLabel,
	)
	p.canvas = canvas
	return canvas, nil
}

func (*RestoreOrCreatePane) staticInitialize() {
	if len(RestoreOrCreateLabels) == 0 {
		RestoreOrCreateLabels = []string{
			i18n.T("WZ.NewWallet.RestoreCreate.LblNewSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.LblRestoreSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.LblNewHardware"),
			i18n.T("WZ.NewWallet.RestoreCreate.LblRestoreHardware"),
		}
	}

	if len(RestoreOrCreateLabelInfo) == 0 {
		RestoreOrCreateLabelInfo = []string{
			i18n.T("WZ.NewWallet.RestoreCreate.InfoNewSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.InfoRestoreSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.InfoNewHardware"),
			i18n.T("WZ.NewWallet.RestoreCreate.InfoRestoreHardware"),
		}
	}
}

// IsValid implements base.WizardPane.
func (p *RestoreOrCreatePane) IsValid() bool {
	if p.canvas == nil {
		// Not yet fully initialized !!!
		return false
	}
	return p.actionGroup.Selected != ""
}

// OnHide implements base.WizardPane.
func (*RestoreOrCreatePane) OnHide() {
}

// OnShow implements base.WizardPane.
func (p *RestoreOrCreatePane) OnShow() {
	if p.wizardState == nil {
		return
	}

	if p.wizardState.RestoreOrCreate != RestoreOrCreateUnknown {
		p.actionGroup.SetSelected(
			RestoreOrCreateLabels[p.wizardState.RestoreOrCreate])
	}
	p.actionGroup.Refresh()
}

// Refresh implements base.WizardPane.
func (*RestoreOrCreatePane) Refresh() error {
	return nil
}

// ResetState implements base.WizardPane.
func (p *RestoreOrCreatePane) ResetState() {
	p.wizardState.RestoreOrCreate = RestoreOrCreateUnknown
}

// SetState implements base.WizardPane.
func (p *RestoreOrCreatePane) SetState(state any) error {
	var ok bool
	if p.wizardState, ok = state.(*WizardState); !ok {
		return errors.New(i18n.T("Err.ConvertWizardState"))
	}
	return nil
}

func (p *RestoreOrCreatePane) CanTransitionTo(state any) bool {
	var wizardState *WizardState
	var ok bool
	if wizardState, ok = state.(*WizardState); !ok {
		return false
	}
	return wizardState.Network != "" && wizardState.WalletName != ""
}

// OnBeforeNext implements base.WizardPane.
func (p *RestoreOrCreatePane) OnBeforeNext() {
	action := p.actionGroup.Selected
	idx := slices.IndexFunc(
		RestoreOrCreateLabels,
		func(a string) bool { return a == action },
	)
	p.wizardState.RestoreOrCreate = RestoreOrCreate(idx)
}

// OnBeforePrevious implements base.WizardPane.
func (p *RestoreOrCreatePane) OnBeforePrevious() {
}

func (p *RestoreOrCreatePane) onActionGroupChanged(action string) {
	idx := slices.IndexFunc(
		RestoreOrCreateLabels,
		func(a string) bool { return a == action },
	)
	if idx != -1 {
		p.infoLabel.Text = RestoreOrCreateLabelInfo[idx]
		p.infoLabel.Refresh()
	}
	p.AbstractWizardPane.NotifyOnChanged()
}
