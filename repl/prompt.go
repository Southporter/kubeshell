package repl

import (
	"fmt"
	"log"

	"github.com/peterh/liner"
	"github.com/ssedrick/kubeshell/repl/state"
)

type Prompt struct {
	state   *state.State
	liner *liner.State
}

func NewPrompt(s *state.State) Prompt {
	return Prompt{
		state:   s,
		liner: liner.NewLiner(),
	}
}

func (p *Prompt) Get() (string, error) {
	prompt := fmt.Sprintf("%s:%s$ ", p.state.CurrentCluster(), p.state.Namespace)
	if input, err := p.liner.Prompt(prompt); err == nil {
		return input, err
	} else if err == liner.ErrPromptAborted {
		log.Print("Got err aborted")
		return "", err
	} else {
		log.Print("Got a different error", err)
		return "", nil
	}
}

func (p *Prompt) Close() {
	p.liner.Close()
}
