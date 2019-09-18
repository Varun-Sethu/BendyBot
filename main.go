package main

import(
	_"bendy-bot/internal/markov"
	"bendy-bot/internal/client"
	"bendy-bot/internal/tracking"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"fmt"
)


var botID string

func main() {
	discord, err := discordgo.New("Bot " + "NjIyNDMwMTg3MjcwODMyMTI4.XXzzUw.D0P2NRR-mpAIOOpx6Ck8PDJoX-c")
	if err != nil {
		panic(err)
	}
	u, err := discord.User("@me")
	if err != nil {
		panic(err)
	}
	
	
	botID = u.ID
	err = discord.Open()
	if err != nil {
		panic(err)
	}

	discord.AddHandler(messageCreate)
	defer discord.Close()

	<-make(chan struct{})


}




func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	r, _ := regexp.Compile(`^(yo bendy)(\s.+)`)
	
	// This is what deals with our incoming commands
	if matches := r.MatchString(m.Content); matches {
		// Extract the actual command from the message
		command := r.FindAllStringSubmatch(m.Content, 2)[0][2]

		tokens, err := client.Parse(m, command)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		} else {
			s.ChannelMessageSend(m.ChannelID, client.Interpret(tokens, m))
		}
	} else {
		if m.ChannelID == tracking.StateInfo.Channel {
			// Parse the message on to be tracked
			if _, ok := tracking.CurrentlyTracking[m.Author.ID]; ok {
				tracking.Store(m.Author.ID, m.Content)
			}
		}
	}
}