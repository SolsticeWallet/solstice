package new_wallet

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
	PasswordHash    string
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
