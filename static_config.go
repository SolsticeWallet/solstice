package solstice

import "net/url"

var SupportedNetworks = []string{
	"Ethereum",
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
