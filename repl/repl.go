package repl

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
)

func Start(cmd *cobra.Command, args []string) {
	// Set up kube client

	// Start loop
	l := liner.NewLiner()
	defer l.Close()
	l.SetCtrlCAborts(true)
	for {
		prompt := fmt.Sprintf("%s:%s$ ", "cluster", "namespace")
		input, err := l.Prompt(prompt)
		if err == liner.ErrPromptAborted {
			log.Print("Got err prompt aborted")
			return
		} else if err != nil {
			log.Print("Error reading line", err)
			return
		}
		log.Print("Got input: ", input)
		color.Blue(input)
	}
}
