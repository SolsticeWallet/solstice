package solstice

import (
	"io/fs"
	"os"
	"path"
)

var AppStorageDir = ""
var AppCacheDir = ""
var AppHomeDir = ""
var AppWalletDir = ""
var AppLogDir = ""

func initializeStorage() {
	initAppStorageDir()
	initAppCacheDir()
	initAppHomeDir()
}

func initAppStorageDir() {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	AppStorageDir = path.Join(dir, GetMetadata().ID)
	if err = ensureDir(AppStorageDir); err != nil {
		panic(err)
	}
}

func initAppCacheDir() {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	AppCacheDir = path.Join(dir, GetMetadata().ID)
	if err = ensureDir(AppCacheDir); err != nil {
		panic(err)
	}

}

func initAppHomeDir() {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	AppHomeDir = path.Join(dir, HiddenFolderName(GetMetadata().Name))
	if err = ensureDir(AppHomeDir); err != nil {
		panic(err)
	}

	AppWalletDir = path.Join(AppHomeDir, "wallets")
	if err = ensureDir(AppWalletDir); err != nil {
		panic(err)
	}

	AppLogDir = path.Join(AppHomeDir, "logs")
	if err = ensureDir(AppLogDir); err != nil {
		panic(err)
	}
}

func ensureDir(d string, p ...fs.FileMode) (err error) {
	dm := fs.FileMode(0755)
	if len(p) > 0 {
		dm = p[0]
	}

	if _, err = os.Stat(d); os.IsNotExist(err) {
		err = os.Mkdir(d, dm)
	}
	return err
}
