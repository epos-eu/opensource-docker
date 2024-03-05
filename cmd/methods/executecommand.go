//go:build !windows
// +build !windows

package methods

import (
	"os"
	"os/exec"
)

func ExecuteCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		PrintError("Error on executing command, cause: " + err.Error())
		return err
	}
	return nil
}
