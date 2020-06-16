package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func helpAndExit(usage, summary, details string, defs []definitionList) {
	s := formatHelp(usage, summary, details, defs)
	fmt.Fprintf(os.Stdout, s)
	os.Exit(0)
}

func formatHelp(usage, summary, details string, defs []definitionList) string {
	columns := terminalColumns()
	sections := []string{}
	sections = append(sections, wrapParagraphs(usage, columns))
	if summary != "" {
		sections = append(sections, wrapParagraphs(summary, columns))
	}
	for _, d := range defs {
		if len(d.definitions) == 0 {
			continue
		}
		sections = append(sections, d.format(columns))
	}
	if details != "" {
		sections = append(sections, wrapParagraphs(details, columns))
	}
	return strings.Join(sections, "\n")
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

var whitespaceRe = regexp.MustCompile(`\s+`)

func wrapParagraphs(text string, columns int) string {
	b := new(strings.Builder)
	paragraphs := strings.Split(text, "\n\n")
	for _, p := range paragraphs {
		lines := wrapText(p, columns)
		for _, line := range lines {
			fmt.Fprintln(b, line)
		}
	}
	return b.String()
}

func wrapText(text string, columns int) []string {
	// split text into words and convert to []rune
	words := whitespaceRe.Split(text, -1)
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

type definitionList struct {
	title       string
	definitions []*definition
}

func (d *definitionList) format(columns int) string {
	b := new(strings.Builder)

	// set a maximum for left column
	maxLeftCols := (columns - 4) / 2
	if maxLeftCols > 25 {
		maxLeftCols = 25
	}

	// get text for left column
	terms := []termLines{}
	for _, def := range d.definitions {
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

	// print text
	fmt.Fprintf(b, "%s:\n", d.title)
	for i, def := range d.definitions {
		flagDef := terms[i]
		for _, line := range flagDef.separate {
			fmt.Fprintf(b, "  %s\n", line)
		}
		usageLines := wrapText(def.text, rightCols)
		fmt.Fprintf(b, "  %-*s  %s\n", leftCols, flagDef.inline, usageLines[0])
		for _, line := range usageLines[1:] {
			fmt.Fprintf(b, "%*s%s\n", 2+leftCols+2, "", line)
		}
	}

	return b.String()
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
