package utils

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"io"
	"strings"
	"regexp"
	"log"
)

var (
	Token         string
	BotID         string
	Connected     bool
	DiscordClient *discordgo.Session
	TextInput io.WriteCloser
	RXMentions *regexp.Regexp
	RXUserID *regexp.Regexp
	ChannelID string
)

//https://discordapp.com/oauth2/authorize?client_id=308269172314603520&scope=bot&permissions=0
func LoadDiscord(token string) {
	if strings.Contains( token,"<Discord Bot Token>") {
		return
	}
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
	Connected = true

	r, err := regexp.Compile(`<@&?!?([0-9]+)>`)
	if err != nil {
		log.Fatal(err)
	}
	RXMentions = r

	r2, err := regexp.Compile(`\d+`)
	if err != nil {
		log.Fatal(err)
	}
	RXUserID = r2

	DiscordClient.UpdateStatus(0, "Factorio")
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.Username + ":" + m.Content)
	if m.Author.ID == BotID {
		return
	}
	if m.ChannelID != ChannelID {
		return
	}

	io.WriteString(TextInput, "<" + m.Author.Username + "> " + m.Content + "\n")
}

func SendStringToDiscord(message string, channelID string){
	if ! Connected {
		return
	}

	DiscordClient.ChannelMessageSend(channelID, escapeMentions(message))
}

func escapeMentions(message string) string {
	str := message
	if RXMentions.MatchString(message) {
		matches := RXMentions.FindAllString(message, -1)
		for _,match := range matches {
			user,_ := DiscordClient.User(RXUserID.FindString(match))
			fmt.Println(RXUserID.FindString(match))
			username := "unknown"
			if user != nil {
				username = user.Username
			}
			str = strings.Replace(str, match, username, -1)
		}
	}
	str = strings.Replace(str, "@", "", -1)
	return str
}
