package new_wallet

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/solsticewallet/solstice/ui/base"
)

type EnryptwalletPane struct {
	*base.AbstractWizardPane

	wizardState *WizardState

	canvas        fyne.CanvasObject
	passwordEntry *widget.Entry
}
