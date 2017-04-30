package main

import (
	"fmt"
	"github.com/modmuss50/FactorioWrapper/utils"
	"os"
)

func main() {
	fmt.Println("Starting wrapper")

	dataDir := "./data/"
	version := "0.15.4"
	tarBal := fmt.Sprintf("%vfactorio_headless_x64_%v.tar.xz", dataDir, version)
	gameDir := dataDir

	if !utils.FileExists(dataDir) || ! utils.FileExists(tarBal) {
		utils.MakeDir(dataDir)
		if !utils.DownloadURL(fmt.Sprintf("https://www.factorio.com/get-download/%v/headless/linux64", version), tarBal) {
			fmt.Println("Failed to download")
			os.Exit(1)
		}
		utils.ExtractTarXZ(tarBal, gameDir)
	}


}
