package main

import (
	"fmt"
	"github.com/modmuss50/FactorioWrapper/utils"
	"github.com/modmuss50/FactorioWrapper/config"
	"os"
	"bufio"
	"io"
	"log"
	"strings"
	"time"
	"os/exec"
	"strconv"
)

var (
	input io.WriteCloser
	FactorioProcess *exec.Cmd
	ScheduleRestart bool
)



func main() {
	fmt.Println("Starting wrapper")

	dataDir := "./data/"
	if !utils.FileExists(dataDir){
		utils.MakeDir(dataDir)
	}

	fmt.Println("Loading config")
	config.LoadConfig()

	version := config.FactorioVersion
	processDir := utils.GetProcessDir(version)

	utils.HandleDownload(dataDir, version)

	fmt.Println("Connecting to discord")
	utils.DiscordAdmin = config.DiscordAdmin
	utils.LoadDiscord(config.DiscordToken)

	//Generates a blank world is one doesnt exist
	if(!doesSaveExist(processDir)){
		fmt.Println("Generating world...")
		factorioProcess := getExec(processDir, "--create")
		factorioOutput, _ := factorioProcess.StdoutPipe()
		scanner := bufio.NewScanner(factorioOutput)
		go func() {
			for scanner.Scan() {
				text := scanner.Text()
				fmt.Printf("\t > %s\n", text)
			}
		}()
		factorioProcess.Start()
		factorioProcess.Wait()
		fmt.Println("Done")
		os.Exit(1)
	}


	startGame(processDir, version)

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for range ticker.C {
			if utils.RequestRestart {
				utils.RequestRestart = false
				RequestRestart()
			}

		}
	}()

	readInput()

}

func startGame(processDir string, version string) *exec.Cmd {
	if FactorioProcess != nil {
		fmt.Println("Is the game allready running?")
	}
	fmt.Println("Starting game...")
	factorioProcess := getExec(processDir, "--start-server")

	factorioInput, err := factorioProcess.StdinPipe()
	input = factorioInput
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
			fmt.Printf("\t > %s\n", text)
			if strings.Contains(text, "changing state from(CreatingGame) to(InGame)") {
				utils.ChannelID = config.DiscordChannel
				utils.SendStringToDiscord("Server started on factorio version " + version, config.DiscordChannel)
			} else
			if strings.Contains(text, "changing state from(CreatingGame) to(InitializationFailed)") || strings.Contains(text, "Couldn't acquire exclusive lock for") {
				fmt.Println("Game failed to start")
				fmt.Println("Is a server allready running? Check Task Manager")
				factorioProcess.Process.Kill()
				os.Exit(1)
			} else
			if strings.Contains(text, "[JOIN]") {
				utils.SendStringToDiscord(text[26:], config.DiscordChannel)
			} else
			if strings.Contains(text, "[CHAT]") {
				if !strings.Contains(text, " [CHAT] <server>:") {
					utils.SendStringToDiscord(text[26:], config.DiscordChannel)
				}
			} else
			if strings.Contains(text, "[LEAVE]") {
				utils.SendStringToDiscord(text[27:], config.DiscordChannel)
			} else
			if strings.Contains(text, "Goodbye") {
				if ScheduleRestart {
					handleRestart(processDir, version)
				} else {
					utils.SendStringToDiscord("Server closed", config.DiscordChannel)
					time.Sleep(1 * time.Second)
					os.Exit(0)
				}

			}
		}
	}()

	fmt.Println("Launching process")
	er := factorioProcess.Start()
	if er != nil {
		log.Fatal(er)
		os.Exit(1)
	}
	FactorioProcess = factorioProcess
	return factorioProcess
}

func readInput() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if strings.HasPrefix(text, "stop") {
		fmt.Println("Stopping server...")
		utils.KillProcess(FactorioProcess)
		FactorioProcess = nil

	}
	if strings.HasPrefix(text, "fstop") {
		fmt.Println("Stopping server")
		FactorioProcess.Process.Kill()
		FactorioProcess = nil
		return
	}

	ranCmd := false
	if strings.HasPrefix(text, "cmd") {
		ranCmd = true
		io.WriteString(input, strings.Replace(text, "cmd ", "", -1))
	} else
	if strings.HasPrefix(text, "restartnow") {
		ranCmd = true
		RequestInstantRestart()
	} else
	if strings.HasPrefix(text, "restart") {
		ranCmd = true
		RequestRestart()
	}

	if !ranCmd {
		fmt.Println("Command not found!: " + text)
	}

	readInput()
}

func SendGloabMessage(message string){
	fmt.Println(message)
	io.WriteString(input, message + "\n")
	utils.SendStringToDiscord(message, config.DiscordChannel)
}

func RequestRestart(){
	SendGloabMessage("Server Restarting in 30 seconds")
	time.Sleep(15 * time.Second)
	SendGloabMessage("Server Restarting in 15 seconds")
	time.Sleep(5 * time.Second)
	for i:= 10; i>0 ; i--{
		SendGloabMessage("Server Restarting in " + strconv.Itoa(i) +" seconds!")
		time.Sleep(time.Second)
	}
	RequestInstantRestart()
}

func RequestInstantRestart(){
	fmt.Println("Restarting server...")
	utils.SendStringToDiscord("Server restarting...", config.DiscordChannel)
	ScheduleRestart = true
	utils.KillProcess(FactorioProcess)
}

func handleRestart(processDir string, version string) {
	ScheduleRestart = false
	fmt.Println("Restating game")
	FactorioProcess = startGame(processDir, version)
}

func getExec(dir string, mode string) *exec.Cmd {
	fullDir := utils.GetBinPath()
	fullDir = utils.FormatPath(fullDir)
	factorioExec := exec.Command(fullDir, mode, "./saves/" + config.FactorioSaveFileName)
	factorioExec.Dir = "." + dir
	return factorioExec
}


func doesSaveExist(dir string) bool {
	return utils.FileExists("." + dir + "/saves/" + config.FactorioSaveFileName)
}