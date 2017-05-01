// +build linux
package utils

import (
	"os/exec"
	"syscall"
)

func KillProcess(cmd *exec.Cmd){
	syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
}