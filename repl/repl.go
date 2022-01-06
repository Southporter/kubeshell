package repl

import (
	"os"
	"strings"

	log "k8s.io/klog"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ssedrick/kubeshell/repl/builtins"
	"github.com/ssedrick/kubeshell/repl/cmd"
	"github.com/ssedrick/kubeshell/repl/state"
)

type commandHandler = func(*cmd.Command) error

var commands = map[string]commandHandler{
	"ls":   builtins.Ls,
	"pwd":  builtins.Pwd,
	"cd":   builtins.Cd,
	"exit": builtins.Exit,
}

func handleCommand(command string, state *state.State) {
	if command == "" {
		return
	}
	parts := strings.Split(command, " ")
	c := parts[0]
	if handler, ok := commands[c]; ok {
		cmd := &cmd.Command{
			Args:  parts[1:],
			State: state,
		}
		err := handler(cmd)
		if err != nil {
			color.Red("Command Not Found")
		}
	}
}

func Start(cmd *cobra.Command, args []string) {
	// Set up kube client
	s := state.NewState()
	if err := s.Load(cmd); err != nil {
		log.Infoln("Error loading kubernetes state")
		os.Exit(2)
	}

	// Start loop
	p := NewPrompt(&s)
	defer p.Close()
	for {
		input, err := p.Get()
		if err != nil {
			log.Info("Error reading line", err)
			return
		}
		log.Info("Got input: ", input)

		handleCommand(input, &s)
	}
}
