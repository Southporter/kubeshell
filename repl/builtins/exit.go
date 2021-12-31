package builtins

import (
	"os"

	"github.com/ssedrick/kubeshell/repl/cmd"
	log "k8s.io/klog"
)

func Exit(cmd *cmd.Command) error {
	log.Infoln("Exiting")
	os.Exit(0)
	return nil
}
