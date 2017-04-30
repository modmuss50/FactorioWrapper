package config

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/modmuss50/FactorioWrapper/utils"
	"os"
)

var (
	FactorioVersion string
	FactorioSaveFileName string
	DiscordToken   string
	DiscordChannel string
)

var DefaultConfig = []byte(`

[factorio]
version = "0.15.4"
save = "world.zip"

[discord]
token = "<Discord Bot Token>"
channelID = "<Discord Channel ID>"
commanderRole = "Bot Admin"

`)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("data")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found, generating defualt and closing...")
		utils.WriteStringToFile(string(DefaultConfig), "./data/config.toml")
		os.Exit(0)

	} else {
		FactorioVersion = viper.GetString("factorio.version")
		FactorioSaveFileName = viper.GetString("factorio.save")
		DiscordToken = viper.GetString("discord.token")
		DiscordChannel = viper.GetString("discord.channelID")
		DiscordChannel = viper.GetString("discord.commanderRole")

	}
}
