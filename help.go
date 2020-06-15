package cmd

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

func printHelp(usage, summary, details string, flags *Flags) {
	columns := terminalColumns()
	w := os.Stdout
	fmt.Fprintln(w, usage)
	fmt.Fprintln(w)
	if summary != "" {
		fmt.Fprintln(w, summary)
		fmt.Fprintln(w)
	}
	flags.printHelp(w, columns)
	if details != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, details)
	}
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func terminalColumns() int {
	// try $COLUMNS env variable
	cols, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err == nil {
		return cols
	}

	// try syscall
	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) != -1 {
		return int(ws.Col)
	}

	return 80
}
