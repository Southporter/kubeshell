package builtins

import (
	"errors"
	"strings"

	"github.com/ssedrick/kubeshell/repl/cmd"
)

func Ls(cmd *cmd.Command) error {
	parts := strings.Split(cmd.State.CurrentDirectory(), "/")
	switch len(parts) {
	case 0:
		// print all resource types
	case 1:
		// if namespace, print all namespaces
		// if resouce, get all of that resource regardless of namespace
	case 2:
		if parts[0] != "namespaces" {
			return errors.New("Unknown base folder")
		}
	}
	return nil
}
