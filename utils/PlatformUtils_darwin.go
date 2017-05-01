// +build darwin
package utils

import (
	"os/exec"
	"syscall"
)

func KillProcess(cmd *exec.Cmd){
	syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
}

func HandleDownload(){

}

func GetBinPath() string {
	return "./bin/x64/factorio"
}

func GetProcessDir(version string) string {
	return "/data/factorio"
}