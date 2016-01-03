package common

import (
	"os"
	"fmt"
)

func logAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
