//go:build !(android || ios)

package wizards

import (
	"fyne.io/fyne/v2"
	"github.com/solsticewallet/solstice/ui/base"
)

func ShowWizard(
	WizardCreator WizardCreator,
	cancelCallback base.WizardCancelCallback,
	confirmCallback base.WizardConfirmCallback,
	size ...fyne.Size,
) (base.WizardView, error) {
	return ShowWizardInWindow(
		WizardCreator,
		cancelCallback, confirmCallback,
		size...,
	)
}
