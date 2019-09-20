package main

import(
	"bendy-bot/internal"
	_"bendy-bot/internal/markov"
	"bendy-bot/internal/client"
	"bendy-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)


var botID string

func main() {
	// Authenticate the client
	discord, err := discordgo.New("Bot " + getAuthCode())
	if err != nil {
		panic(err)
	}
	
	// Set up the bot and the handlers
	u, err := discord.User("@me")
	if err != nil {
		panic(err)
	}
	botID = u.ID
	err = discord.Open()
	discord.AddHandler(messageCreate)

	// Keep the bot running
	fmt.Println("Bot is now running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	
	defer discord.Close()
}



// Function to get the bot's oauth code
func getAuthCode() string {
	f, _ := os.Open(internal.GetAbsFile("authcode.txt"))
	codeBytes, _ := ioutil.ReadAll(f)
	defer f.Close()

	return string(codeBytes)
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
		client.HandleIncomingCommand(s, m)
		break;
	case false:
		bot.HandleIncomingMessage(m)	
		break;
	}
}