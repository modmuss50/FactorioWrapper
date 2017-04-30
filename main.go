package main

import (
	"fmt"
	"github.com/modmuss50/FactorioWrapper/utils"
	"github.com/modmuss50/FactorioWrapper/config"
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

	config.LoadConfig()

	dataDir := "./data/"
	version := config.FactorioVersion
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
	utils.TextInput = factorioInput
	if err != nil {
		log.Fatal(err)
	}
	defer factorioInput.Close()
	factorioOutput, _ := factorioProcess.StdoutPipe()

	scanner := bufio.NewScanner(factorioOutput)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			if strings.Contains(text, "ping") {
				io.WriteString(factorioInput, "Pong!\n")
			}
			if strings.Contains(text, "changing state from(CreatingGame) to(InGame)") {
				utils.LoadDiscord(config.DiscordToken)
				utils.SendStringToDiscord("Server started on factorio version " + version, config.DiscordChannel)
			}
			if strings.Contains(text, "[JOIN]") {
				utils.SendStringToDiscord(text[26:], config.DiscordChannel)
			}
			if strings.Contains(text, "[CHAT]") {
				if !strings.Contains(text, " [CHAT] <server>:") {
					utils.SendStringToDiscord(text[26:], config.DiscordChannel)
				}
			}
			if strings.Contains(text, "[LEAVE]") {
				utils.SendStringToDiscord(text[27:], config.DiscordChannel)
			}
			if strings.Contains(text, "Goodbye") {
				utils.SendStringToDiscord("Server closed", config.DiscordChannel)
				utils.DiscordClient.Close()
				os.Exit(0)
			}
			fmt.Printf("\t > %s\n", text)
		}
	}()

	factorioProcess.Start()

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			//io.WriteString(factorioInput, "hello is this working?\n")
		}
	}()

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
	factorioExec := exec.Command(utils.GetRunPath()+dir+"/bin/x64/factorio", "--start-server", utils.GetRunPath()+dir+"/saves/" + config.FactorioSaveFileName)
	return factorioExec
}
