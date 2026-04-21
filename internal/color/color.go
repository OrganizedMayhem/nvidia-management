// Package color provides lightweight ANSI terminal colour helpers.
// It mirrors the parts of github.com/fatih/color used by this project so the
// dependency can be swapped back in later with no call-site changes.
//
// Colours are automatically disabled when the output file descriptor is not a
// terminal (e.g. when stdout/stderr is piped or redirected).
package color

import (
	"fmt"
	"io"
	"os"
)

// Attribute is an ANSI SGR parameter value.
type Attribute int

const (
	Reset Attribute = 0

	Bold  Attribute = 1
	Faint Attribute = 2

	FgRed    Attribute = 31
	FgGreen  Attribute = 32
	FgYellow Attribute = 33
	FgCyan   Attribute = 36
)

// NoColor can be set to true to suppress all colour output globally.
var NoColor = !isTerminal(os.Stdout)

// Color holds a set of ANSI attributes and exposes Printf/Println/Fprintf helpers.
type Color struct {
	attrs []Attribute
}

// New creates a Color with the given attributes.
func New(attrs ...Attribute) *Color { return &Color{attrs: attrs} }

func (c *Color) sequence() string {
	if NoColor || len(c.attrs) == 0 {
		return ""
	}
	s := "\033["
	for i, a := range c.attrs {
		if i > 0 {
			s += ";"
		}
		s += fmt.Sprintf("%d", a)
	}
	return s + "m"
}

func (c *Color) wrap(s string) string {
	if NoColor {
		return s
	}
	return c.sequence() + s + "\033[0m"
}

// Sprint returns the string wrapped in ANSI escape codes.
func (c *Color) Sprint(a ...any) string { return c.wrap(fmt.Sprint(a...)) }

// Sprintf formats then wraps in ANSI escape codes.
func (c *Color) Sprintf(format string, a ...any) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// Printf prints to stdout with colour.
func (c *Color) Printf(format string, a ...any) (int, error) {
	return fmt.Print(c.Sprintf(format, a...))
}

// Println prints a line to stdout with colour.
func (c *Color) Println(a ...any) (int, error) {
	return fmt.Println(c.Sprint(a...))
}

// Fprintf prints to w with colour.
func (c *Color) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, c.Sprintf(format, a...))
}

// Fprintln prints a line to w with colour.
func (c *Color) Fprintln(w io.Writer, a ...any) (int, error) {
	return fmt.Fprintln(w, c.Sprint(a...))
}

// isTerminal reports whether f is connected to a terminal.
// It uses a simple /dev/tty check that works without importing golang.org/x/sys.
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	// ModeCharDevice is set for ttys; ModeDevice alone covers block devices.
	return fi.Mode()&os.ModeCharDevice != 0
}
