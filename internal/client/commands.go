package client

import (
	"bendy-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)



// Utility functions
func validate(input []string, m *discordgo.MessageCreate, inputs int, mentions int) bool {
	if len(input) != inputs + 1 || len(m.Mentions) != mentions {
		return false
	}

	return true
}





// function that handles a request from a user to begin the tracking of another user
func track(input []string, m *discordgo.MessageCreate) string {
	// Generate and pass user id + determine if the inputs are valid
	uid := m.Mentions[0].ID
	if !validate(input, m, 1, 1) {
		return "You have to/can only mention 1 person :(."
	}
	
	err := bot.BeginTrackingUser(uid)
	if err != nil {
		return err.Error()
	}

	return "Yay! I'm now tracking: " + uid
}



// Function that handles a request from a user to end the tracking of another user
func endtrack(input []string, m *discordgo.MessageCreate) string {
	uid := m.Mentions[0].ID
	if !validate(input, m, 1, 1) {
		return "You have to/can only mention 1 person :(."
	}

	bot.EndTrackingUser(uid)	

	return "Successfully ended tracking for: " + uid
}



// Function that handles the request to generate a sentence from another user's markov chain
func generate(input []string, m *discordgo.MessageCreate) string {
	uid := m.Mentions[0].ID
	if !validate(input, m, 1, 1) {
		return "You have to/can only mention 1 person :(."
	}

	sentence, err := bot.GenerateSentenceForUser(uid)
	if err != nil {
		return err.Error()
	}
	
	return sentence + "."
}



// Function that handles the request from a user to change the channel that is actively being tracked from
func setChannel(input []string, m *discordgo.MessageCreate) string {
	if !validate(input, m, 0, 0) {
		return "Invalid format >:("
	}
	chanelID := m.ChannelID
	bot.SetTrackingChannel(chanelID)
	return "Successfully changed channel to: " + chanelID
}