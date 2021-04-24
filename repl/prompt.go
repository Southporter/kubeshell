package repl

import (
	"fmt"
	"log"

	"github.com/peterh/liner"
	"github.com/ssedrick/kubeshell/repl/state"
)

type Config struct {
	namespace string
	cluster   string
}

type Prompt struct {
	config Config
	liner  *liner.State
	state  *state.State
}

func NewPrompt(s *state.State) Prompt {
	l := liner.NewLiner()
	l.SetCtrlCAborts(true)
	return Prompt{
		state: s,
		config: Config{
			namespace: s.Namespace,
			cluster:   s.CurrentCluster(),
		},
		liner: l,
	}
}

func (p *Prompt) Get() (string, error) {
	prompt := fmt.Sprintf("%s:%s$ ", p.config.cluster, p.config.namespace)
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
