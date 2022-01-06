package builtins

import (
	"strings"

	"github.com/ssedrick/kubeshell/repl/cmd"
)

func popDir(cwd string) string {
	parts := strings.Split(cwd, "/")
	parts = parts[:len(parts)-1]
	return strings.Join(parts, "/")
}

func handleChange(cwd, arg string) (string, error) {
	if strings.HasPrefix(arg, "/") {
		return arg, nil
	}
	if arg == ".." {
		return popDir(cwd), nil
	}
	if arg == "." {
		return cwd, nil
	}
	if strings.HasPrefix(arg, "../") {
		cwd = popDir(cwd)
		arg = strings.TrimPrefix(arg, "../")
		return handleChange(cwd, arg)
	}
	if strings.HasPrefix(arg, "./") {
		arg = strings.TrimPrefix(arg, "./")
		return handleChange(cwd, arg)
	}
	if arg == "" {
		return cwd, nil
	}
	cwd += "/" + arg
	return cwd, nil
}

/****
 * Currently naive. Does not check to see if directory exists
 */
func Cd(cmd *cmd.Command) error {
	cwd := cmd.State.CurrentDirectory()
	arg := cmd.Args[0]
	if arg == "" {
		cmd.State.SetCurrentDirectory("/namespaces/" + cmd.State.Namespace)
		return nil
	}
	cwd, err := handleChange(cwd, arg)
	if err != nil {
		return err
	}
	cmd.State.SetCurrentDirectory(cwd)
	return nil
}
