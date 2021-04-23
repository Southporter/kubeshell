package repl

import (
	"fmt"
	"log"

	"github.com/peterh/liner"
)

type Config struct {
	namespace string
	cluster   string
}

type Prompt struct {
	config Config
	liner  *liner.State
	state  *State
}

func NewPrompt(s *State) Prompt {
	l := liner.NewLiner()
	l.SetCtrlCAborts(true)
	return Prompt{
		state: s,
		config: Config{
			s.namespace,
			s.cluster,
		},
		liner: l,
	}
}

func (p *Prompt) Get() (string, error) {
	prompt := fmt.Sprintf("%s:%s$ ", "cluster", "namespace")
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
