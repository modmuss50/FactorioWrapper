package utils

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"io"
)

var (
	Token         string
	BotID         string
	Connected     bool
	DiscordClient *discordgo.Session
	TextInput io.WriteCloser
)

//https://discordapp.com/oauth2/authorize?client_id=308269172314603520&scope=bot&permissions=0
func LoadDiscord(token string) {
	Token = token
	dg, err := discordgo.New("Bot " + Token)
	DiscordClient = dg
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}
	BotID = u.ID

	dg.AddHandler(handleMessage)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	Connected = true
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.Username + ":" + m.Content)
	if m.Author.ID == BotID {
		return
	}
	io.WriteString(TextInput, "<" + m.Author.Username + "> " + m.Content + "\n")
}

func SendStringToDiscord(message string, channelID string){
	DiscordClient.ChannelMessageSend(channelID, message)
}
