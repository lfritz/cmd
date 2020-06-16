package cmd

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// terminalColumns finds out the number of columns in the terminal, or returns a default.
func terminalColumns() int {
	// try $COLUMNS env variable
	cols, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err == nil {
		return cols
	}

	// try syscall
	ws := new(struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	})
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) != -1 {
		return int(ws.Col)
	}

	return 80
}
