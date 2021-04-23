package repl

import (
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type commandHandler = func(args []string) error

func listCommand(args []string) error {
	println("namespace\tpods\tservices")
	println("")
	return nil
}

var commands = map[string]commandHandler{
	"ls": listCommand,
}

func handleCommand(command string) {
	if command == "" {
		return
	}
	parts := strings.Split(command, " ")
	c := parts[0]
	if handler, ok := commands[c]; ok {
		err := handler(parts[1:])
		if err != nil {
			color.Red("Command Not Found")
		}
	}
}

func Start(cmd *cobra.Command, args []string) {
	// Set up kube client
	s := NewState()
	if err := s.Load(cmd); err != nil {
		log.Println("Error loading kubernetes state")
		os.Exit(2)
	}

	// Start loop
	p := NewPrompt(&s)
	defer p.Close()
	for {
		input, err := p.Get()
		if err != nil {
			log.Print("Error reading line", err)
			return
		}
		log.Print("Got input: ", input)

		handleCommand(input)
	}
}
