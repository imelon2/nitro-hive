package path_test

import (
	"fmt"
	"testing"

	"github.com/imelon2/nitro-hive/common/path"
)

func Test_Path(t *testing.T) {
	aPath := path.AccountPath()
	pPath := path.PrivateKeyPath()

	fmt.Printf("account path     : %s\n", aPath)
	fmt.Printf("private key path : %s\n", pPath)
}
