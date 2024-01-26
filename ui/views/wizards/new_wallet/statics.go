package new_wallet

import "github.com/solsticewallet/solstice/i18n"

var (
	mnemonicTypes            []string
	restoreOrCreateLabels    []string
	restoreOrCreateLabelInfo []string
)

func staticInitialize() {
	if len(mnemonicTypes) == 0 {
		mnemonicTypes = []string{
			i18n.T("Lbl.MnemonicTypes.12Words"),
			i18n.T("Lbl.MnemonicTypes.12WordsPwd"),
			i18n.T("Lbl.MnemonicTypes.24Words"),
			i18n.T("Lbl.MnemonicTypes.24WordsPwd"),
		}
	}

	if len(restoreOrCreateLabels) == 0 {
		restoreOrCreateLabels = []string{
			i18n.T("WZ.NewWallet.RestoreCreate.LblNewSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.LblRestoreSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.LblNewHardware"),
			i18n.T("WZ.NewWallet.RestoreCreate.LblRestoreHardware"),
		}
	}

	if len(restoreOrCreateLabelInfo) == 0 {
		restoreOrCreateLabelInfo = []string{
			i18n.T("WZ.NewWallet.RestoreCreate.InfoNewSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.InfoRestoreSoftware"),
			i18n.T("WZ.NewWallet.RestoreCreate.InfoNewHardware"),
			i18n.T("WZ.NewWallet.RestoreCreate.InfoRestoreHardware"),
		}
	}
}
