package builtins

import (
	"context"
	"errors"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/ssedrick/kubeshell/repl/cmd"
	"github.com/ssedrick/kubeshell/repl/display"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog"
)

func printAllApiResources(cmd *cmd.Command, namespaced bool) error {
	client, err := cmd.State.ToDiscoveryClient()
	if err != nil {
		return err
	}
	list, err := client.ServerResources()
	if err != nil {
		return err
	}

	displayOptions := display.DefaultOptions()
	displayOptions.ChangeDirection(display.TopToBottom)
	displayOptions.ChangePadding(display.NewWhitespacePadding(1))
	resources := display.NewGridWithOptions(displayOptions)
	for _, item := range list {
		if len(item.APIResources) != 0 {
			for _, resource := range item.APIResources {
				if len(resource.Verbs) != 0 {
					if namespaced == resource.Namespaced {
						log.V(2).Infof("resource: %v; group: %s", resource, resource.Group)
						cell := display.NewCell(color.InBlue(resource.Name))
						log.V(2).Infoln("Cell: ", resource.Name)
						resources.AddCell(cell)
					}
				}
			}
		}
	}
	resources.Sort()
	screen := display.NewScreen(resources)
	return screen.Print()
}

func printAllNamespaces(cmd *cmd.Command, namespace string) error {
	client := cmd.State.ToClient()

	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	displayOptions := display.DefaultOptions()
	displayOptions.ChangeDirection(display.TopToBottom)
	displayOptions.ChangePadding(display.NewWhitespacePadding(1))
	resources := display.NewGridWithOptions(displayOptions)
	for _, resource := range namespaces.Items {
		log.V(2).Infof("Resource: %v", resource)
		cell := display.NewCell(color.InBlue(resource.Name))
		resources.AddCell(cell)
	}

	resources.Sort()
	screen := display.NewScreen(resources)
	return screen.Print()
}

func pathToFolders(cwd string) []string {
	parts := strings.Split(cwd, "/")
	for len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	log.Info("parts", parts, len(parts))
	return parts
}

func printResources(name string) error {
	log.Warning("print resources not implemented")
	return nil
}

func Ls(cmd *cmd.Command) error {
	parts := pathToFolders(cmd.State.CurrentDirectory())
	switch len(parts) {
	case 0:
		// print all resource types
		return printAllApiResources(cmd, false)
	case 1:
		if parts[0] == "namespaces" {
			// Fetch all namespaces
			return printAllNamespaces(cmd, parts[0])
		} else {
			// fetch non-namespaced resources
			return printResources(parts[0])
		}
	case 2:
		if parts[0] != "namespaces" {
			return errors.New("unknown base folder")
		}
		return printAllApiResources(cmd, true)
	}
	return nil
}
