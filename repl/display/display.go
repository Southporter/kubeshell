package display

import (
	"os"

	"golang.org/x/term"
)

type Display struct {
	terminal *term.Terminal
	oldState *term.State
}

type IO struct {
}

func (i *IO) Read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}

func (o *IO) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func NewDisplay(prompt string) Display {
	com := &IO{
	}
	return Display{
		terminal: term.NewTerminal(com, prompt),
	}
}

func (d *Display) Close() {
	if d.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), d.oldState)
	}
}
