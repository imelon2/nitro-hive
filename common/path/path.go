package path

import (
	"log"
	"path/filepath"
	"runtime"
)

var (
	ACCOUNT_FILE_NALE     = "accounts"
	PRIVATE_KEY_FILE_NALE = "privateKey"
)

func AccountPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("bad path")
	}

	root := filepath.Dir(filename)
	parent := filepath.Dir(filepath.Dir(root))
	path := filepath.Join(parent, "account", ACCOUNT_FILE_NALE)
	return path
}

func PrivateKeyPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("bad path")
	}

	root := filepath.Dir(filename)
	parent := filepath.Dir(filepath.Dir(root))
	path := filepath.Join(parent, "account", PRIVATE_KEY_FILE_NALE)
	return path
}
