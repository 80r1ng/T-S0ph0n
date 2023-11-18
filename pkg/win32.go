package pkg

import (
	"os"
	"syscall"
)

var (
	netapi32 = syscall.MustLoadDLL("netapi32.dll")
)

func captureErr(err error) {
	if err.Error() == "The operation completed successfully." {
		return
	} else {
		os.Exit(0)
	}
}
