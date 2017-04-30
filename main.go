package main

import (
	"fmt"
	"github.com/modmuss50/FactorioWrapper/utils"
	"os"
	"os/exec"
	"time"
	"bufio"
	"io"
	"log"
	"strings"
	"syscall"
)

func main() {
	fmt.Println("Starting wrapper")

	dataDir := "./data/"
	version := "0.15.4"
	tarBal := fmt.Sprintf("%vfactorio_headless_x64_%v.tar.xz", dataDir, version)
	gameDir := dataDir
	proccessDir := "/FactorioWrapper/data/factorio"

	if !utils.FileExists(dataDir) || ! utils.FileExists(tarBal) {
		utils.MakeDir(dataDir)
		if !utils.DownloadURL(fmt.Sprintf("https://www.factorio.com/get-download/%v/headless/linux64", version), tarBal) {
			fmt.Println("Failed to download")
			os.Exit(1)
		}
		utils.ExtractTarXZ(tarBal, gameDir)
	}

	factorioProcess := getExec(proccessDir)
	factorioInput, err := factorioProcess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer factorioInput.Close()
	factorioOutput, _ := factorioProcess.StdoutPipe()

	scanner := bufio.NewScanner(factorioOutput)
	go func() {
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "ping"){
				io.WriteString(factorioInput, "Pong!\n")
			}
			if strings.Contains(scanner.Text(), "Goodbye"){
				os.Exit(0)
			}
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()


	factorioProcess.Start()
	//factorioProcess.Wait()

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			//io.WriteString(factorioInput, "hello is this working?\n")
		}
	}()

	fmt.Println("Hello")
	readInput(factorioProcess)

}

func readInput(cmd *exec.Cmd) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.HasPrefix(text, "stop") {
		fmt.Println("Stopping server...")
		syscall.Kill(cmd.Process.Pid, syscall.SIGINT)

	}
	if strings.HasPrefix(text, "fstop") {
		fmt.Println("Stopping server")
		cmd.Process.Kill()
		return
	}
	fmt.Println("Command not found!")
	readInput(cmd)
}

func getExec(dir string) *exec.Cmd {
	factorioExec := exec.Command( utils.GetRunPath() + dir + "/bin/x64/factorio", "--start-server", utils.GetRunPath() + dir + "/saves/test.zip")
	return factorioExec
}
