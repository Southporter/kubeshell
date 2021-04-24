package cmd

import "github.com/ssedrick/kubeshell/repl/state"

type Command struct {
	Args  []string
	State *state.State
}
