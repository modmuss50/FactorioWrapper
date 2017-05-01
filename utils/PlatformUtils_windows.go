// +build windows
package utils

import (
	"os/exec"
	"fmt"
	"os"
)

func KillProcess(cmd *exec.Cmd){
	cmd.Process.Kill()
}


func HandleDownload(dataDir string, version string){
	zipArchive := fmt.Sprintf("%vFactorio_x64_%v.zip", dataDir, version)
	//binDir := GetProcessDir(version) + "/bin/"
	if !FileExists(zipArchive) {
		fmt.Println(fmt.Sprintf("Please download Factorio_x64_%v.zip from https://www.factorio.com/download/ and place it into the data folder", version))
		os.Exit(0)
	}
	if !FileExists(dataDir)  {
		delDir := dataDir + "bin/"
		if FileExists(delDir){
			DeleteDir(delDir)
		}
		delDir = dataDir + "data/"
		if FileExists(delDir){
			DeleteDir(delDir)
		}

		ExtractZip(zipArchive, dataDir)
		fmt.Println("Done.")
	}
}


func GetBinPath() string {
	fmt.Println(GetRunPath())
	return GetRunPath() +  "/bin/x64/factorio.exe"
}

func GetProcessDir(version string) string {
	return "/data/Factorio_" + version
}