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
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	initDirWithMode(
		&AppStorageDir, configDir, GetMetadata().ID, fs.FileMode(0700))

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	initDir(&AppCacheDir, cacheDir, GetMetadata().ID)

	initHomeDirs()
}

func initHomeDirs() {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	initDirWithMode(&AppHomeDir, dir, HiddenFolderName(GetMetadata().Name), fs.FileMode(0700))
	initDirWithMode(&AppWalletDir, AppHomeDir, "wallets", fs.FileMode(0700))
	initDirWithMode(&AppLogDir, AppHomeDir, "logs", fs.FileMode(0700))
}

func initDir(target *string, baseDir string, subDir string) {
	initDirWithMode(target, baseDir, subDir, fs.FileMode(0755))
}

func initDirWithMode(target *string, baseDir string, subDir string, mode fs.FileMode) {
	*target = path.Join(baseDir, subDir)
	if err := ensureDir(*target, mode); err != nil {
		panic(err)
	}
}

func ensureDir(d string, p fs.FileMode) (err error) {
	if _, err = os.Stat(d); os.IsNotExist(err) {
		err = os.Mkdir(d, p)
	}
	return err
}
