package display

import (
	"io"
	"os"

	"golang.org/x/term"
)

type Display struct {
	terminal *term.Terminal
	oldState *term.State
}

type IO struct {
	Reader io.Reader
	Writer io.Writer
}

type input struct {
}

func Read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}

type output struct {
}

func (o *output) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func NewDisplay(prompt string) Display {
	com := IO{
		Reader: Read,
		Writer: output{},
	}
	return Display{
		terminal: &term.NewTerminal(com, prompt),
	}
}

func (d *Display) Close() {
	if d.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), d.oldState)
	}
}
