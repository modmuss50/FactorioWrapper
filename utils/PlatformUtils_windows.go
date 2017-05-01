// +build windows
package utils

import (
	"os/exec"
)

func KillProcess(cmd *exec.Cmd){
	cmd.Process.Kill()
}