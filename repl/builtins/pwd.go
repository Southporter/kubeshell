package builtins

import "github.com/ssedrick/kubeshell/repl/cmd"

func Pwd(cmd *cmd.Command) error {
	println(cmd.State.CurrentDirectory())
	return nil
}
