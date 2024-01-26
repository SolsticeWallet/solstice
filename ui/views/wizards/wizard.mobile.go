//go:build andoid || ios

package wizareds

import (
	"fyne.io/fyne/v2"
	"github.com/solsticewallet/solstice"
	"github.com/solsticewallet/solstice/ui/base"
)

func ShowWizard(
	creator WizardCreator,
	cancelCallback base.WizardCancelCallback,
	confirmCallback base.WizardConfirmCallback,
	size ...fyne.Size,
) (base.WizardView, error) {
	return ShowWizardInDialog(
		creator,
		cancelCallback, confirmCallback,
		solstice.AppWindow.Content().Size(),
	)
}
