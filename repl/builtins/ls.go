package builtins

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ssedrick/kubeshell/repl/cmd"
	"github.com/ssedrick/kubeshell/repl/display"
)

func printAllApiResources(cmd *cmd.Command, namespaced bool) error {
	client, err := cmd.State.ToDiscoveryClient()
	if err != nil {
		return err
	}
	list, err := client.ServerPreferredResources()
	if err != nil {
		return err
	}

  displayOptions := display.DefaultOptions()
  displayOptions.ChangeDirection(display.TopToBottom)
  resources := display.NewGridWithOptions(displayOptions)
	for _, item := range list {
		if len(item.APIResources) != 0 {
			for _, resource := range item.APIResources {
				if len(resource.Verbs) != 0 {
					if namespaced == resource.Namespaced {
            cell := display.NewCell(resource.Name)
            resources.AddCell(&cell)
					}
				}
			}
		}
	}
  screen := display.NewScreen(&resources)
  return screen.Print()
}

func pathToFolders(cwd string) []string {
	parts := strings.Split(cwd, "/")
	for len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	log.Print("parts", parts, len(parts))
	return parts
}

func Ls(cmd *cmd.Command) error {
	parts := pathToFolders(cmd.State.CurrentDirectory())
	switch len(parts) {
	case 0:
		// print all resource types
		err := printAllApiResources(cmd, false)
		log.Print("error", err)
		return err
	case 1:
		if parts[0] == "namespaces" {
			// Fetch all namespaces
		} else {
			// fetch non-namespaced resources
		}
	case 2:
		if parts[0] != "namespaces" {
			return errors.New("unknown base folder")
		}
	}
	return nil
}
