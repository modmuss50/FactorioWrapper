// +build linux
package utils

import (
	"os/exec"
	"syscall"
	"fmt"
	"os"
)

func KillProcess(cmd *exec.Cmd){
	syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
}

func HandleDownload(dataDir string, version string){
	tarBal := fmt.Sprintf("%vfactorio_headless_x64_%v.tar.xz", dataDir, version)
	if !FileExists(dataDir) || ! FileExists(tarBal) {
		MakeDir(dataDir)
		if !DownloadURL(fmt.Sprintf("https://www.factorio.com/get-download/%v/headless/linux64", version), tarBal) {
			fmt.Println("Failed to download")
			os.Exit(1)
		}

		delDir := dataDir + "bin/"
		if FileExists(delDir){
			DeleteDir(delDir)
		}
		delDir = dataDir + "data/"
		if FileExists(delDir){
			DeleteDir(delDir)
		}

		ExtractTarXZ(tarBal, dataDir)
		fmt.Println("Done.")
	}
}

func GetBinPath() string {
	return "./bin/x64/factorio"
}

func GetProcessDir(version string) string {
	return "/data/factorio"
}