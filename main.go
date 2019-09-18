package main

import(
	"bendy-bot/internal"
	_"bendy-bot/internal/markov"
	"bendy-bot/internal/client"
	"bendy-bot/internal/tracking"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"fmt"
)


var botID string

func main() {
	// Authenticate the client
	authcode := internal.OpenFileFromStore("autcode.txt")
	discord, err := discordgo.New("Bot " + authcode)
	if err != nil {
		panic(err)
	}
	
	// Set up the bot and the handlers
	u, _ := discord.User("@me")
	botID = u.ID
	err = discord.Open()
	discord.AddHandler(messageCreate)
	
	defer discord.Close()
}



// messageCreate handlers handles incoming messages and determines if they should be interpreted or stored within storage
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// regex to determine if the message is a command
	commandRegex, _ := regexp.Compile(`^(yo bendy)(\s.+)`)
	
	// This is what deals with our incoming commands
	switch commandRegex.MatchString(m.Content) {
	case true:
		client.HandleIncomingCommand(m)
		break;
	case false:
		tracking.HandleIncomingMessage(m)	
		break;
	}
}