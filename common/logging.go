package common

import (
	"os"
	"fmt"
)

func LogAndExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
