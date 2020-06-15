package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type helpPrintable interface {
	summary() string
	details() string
	usage() string
	printDefinitions(w io.Writer, columns int)
}

func printHelp(h helpPrintable) {
	// TODO allow multiple paragraphs in summary and details
	columns := terminalColumns()
	w := os.Stdout
	fmt.Fprintln(w, h.usage())
	fmt.Fprintln(w)
	summary := h.summary()
	if summary != "" {
		for _, line := range wrapText(summary, columns) {
			fmt.Fprintln(w, line)
		}
		fmt.Fprintln(w)
	}
	h.printDefinitions(w, columns)
	details := h.details()
	if details != "" {
		for _, line := range wrapText(details, columns) {
			fmt.Fprintln(w, line)
		}
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

func wrapText(text string, columns int) []string {
	// split text into words and convert to []rune
	// TODO split by whitespace, not just space character
	words := strings.Split(text, " ")
	runeWords := make([][]rune, len(words))
	for i, word := range words {
		runeWords[i] = []rune(word)
	}

	// join words into lines
	lines := []string{}
	current := []rune{}
	firstWord := true
	for _, word := range runeWords {
		if len(current)+1+len(word) > columns {
			lines = append(lines, string(current))
			current = word
			continue
		}
		if !firstWord {
			current = append(current, ' ')
		}
		current = append(current, word...)
		firstWord = false
	}
	lines = append(lines, string(current))

	return lines
}

type definition struct {
	terms []string
	text  string
}

type termLines struct {
	separate []string
	inline   string
}

func (d *definition) formatTerms(maxCols int) termLines {
	joined := strings.Join(d.terms, ", ")
	last := len(d.terms) - 1
	if len([]rune(d.terms[last])) > maxCols {
		return termLines{
			separate: d.terms,
		}
	}
	if len([]rune(joined)) > maxCols {
		return termLines{
			separate: d.terms[:last],
			inline:   d.terms[last],
		}
	}
	return termLines{
		inline: joined,
	}
}

func printDefinitions(w io.Writer, defs []*definition, columns int) {
	// set a maximum for left column
	maxLeftCols := (columns - 4) / 2
	if maxLeftCols > 25 {
		maxLeftCols = 25
	}

	// get text for left column
	terms := []termLines{}
	for _, def := range defs {
		terms = append(terms, def.formatTerms(maxLeftCols))
	}

	// find out size of left and right column
	leftCols := 0
	for _, t := range terms {
		if t.inline == "" {
			continue
		}
		cols := len([]rune(t.inline))
		if cols > leftCols {
			leftCols = cols
		}
	}
	rightCols := columns - 4 - leftCols
	if rightCols > 80 {
		rightCols = 80
	}

	// print
	for i, def := range defs {
		flagDef := terms[i]
		for _, line := range flagDef.separate {
			fmt.Fprintf(w, "  %s\n", line)
		}
		usageLines := wrapText(def.text, rightCols)
		fmt.Fprintf(w, "  %-*s  %s\n", leftCols, flagDef.inline, usageLines[0])
		for _, line := range usageLines[1:] {
			fmt.Fprintf(w, "%*s%s\n", 2+leftCols+2, "", line)
		}
	}
}
