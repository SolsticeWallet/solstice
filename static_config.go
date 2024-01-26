package solstice

import (
	"net/url"

	"github.com/solsticewallet/solstice-core/blockchains/networks"
)

var SupportedNetworks = []string{
	networks.Ethereum,
}

var ExternalLinks = map[string]*url.URL{}

func init() {
	var err error

	if ExternalLinks["Bip39Info"], err = url.Parse(
		"https://bitcoinwiki.org/wiki/mnemonic-phrase",
	); err != nil {
		panic(err)
	}
}

var passwordEntropyLevels = []float64{
	50.0, // Below is weak
	60.0, // Below is fair, Equal or higher is strong
}

type PasswordStrength int

const (
	PasswordWeak PasswordStrength = iota
	PasswordFair
	PasswordStrong
)

func PasswordStrengthFromEntropy(entropy float64) PasswordStrength {
	if entropy < passwordEntropyLevels[0] {
		return PasswordWeak
	} else if entropy < passwordEntropyLevels[1] {
		return PasswordFair
	} else {
		return PasswordStrong
	}
}
