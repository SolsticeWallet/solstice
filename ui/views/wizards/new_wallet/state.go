package new_wallet

import (
	"github.com/solsticewallet/solstice-core/blockchains"
	"github.com/solsticewallet/solstice/i18n"
	"github.com/solsticewallet/solstice/ui/base"
)

type RestoreOrCreate int

const (
	RestoreOrCreateUnknown RestoreOrCreate = -1
	CreateSoftwareWallet   RestoreOrCreate = 0
	RestoreSoftwareWallet  RestoreOrCreate = 1
	CreateHardwareWallet   RestoreOrCreate = 2
	RestoreHardwareWallet  RestoreOrCreate = 3
)

type WizardState struct {
	Network         string
	WalletName      string
	RestoreOrCreate RestoreOrCreate
	Mnemonic        string
	Passphrase      string
	BackupOK        bool
	Password        string
}

func NewWizardState() *WizardState {
	return &WizardState{
		Network:         "",
		WalletName:      "",
		RestoreOrCreate: RestoreOrCreateUnknown,
		Mnemonic:        "",
		Passphrase:      "",
	}
}

func (s WizardState) GetResult() any {
	if s.Network == "" {
		return nil
	}
	if s.WalletName == "" {
		return nil
	}
	if s.Mnemonic == "" {
		return nil
	}

	return nil
}

func (s WizardState) GetSoftwareWallet() any {
	wlt, err := blockchains.NewWallet(blockchains.WalletOpts{
		Network:    s.Network,
		Mnemonic:   s.Mnemonic,
		Passphrase: s.Passphrase,
	})
	if err != nil {
		base.Logger.Error(i18n.T("Err.WalletCreation", i18n.T("Err.Arg.Error"), err.Error()))
		return nil
	}
	return wlt
}
