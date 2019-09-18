package client

import (
	"bendy-bot/internal/markov"
	"bendy-bot/internal/tracking"
	"bendy-bot/internal"
	"github.com/bwmarrin/discordgo"
	"os"
	"io/ioutil"
	"fmt"
	"strings"
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
	if _, ok := tracking.CurrentlyTracking[uid]; ok {
		return "Stop being stupid Botond! User is already being tracked >:("
	}

	// Update the tracking state and read from their storage file
	tracking.StateInfo.Active_tracking = append(tracking.StateInfo.Active_tracking, uid)
	dictBytes := internal.OpenFileFromStore(uid)

	if len(dictBytes) == 0 {
		// Build a new markov chain if nothing exists in storage for this user
		tracking.CurrentlyTracking[uid] = markov.Build()
	} else {
		// Interpret and parse their data
		data := internal.FromGOB64(string(dictBytes))
		tracking.CurrentlyTracking[uid] = data
	}

	return "Yay! I'm now tracking: " + uid
}



// Function that handles a request from a user to end the tracking of another user
func endtrack(input []string, m *discordgo.MessageCreate) string {
	uid := m.Mentions[0].ID
	if !validate(input, m, 1, 1) {
		return "You have to/can only mention 1 person :(."
	}

	// This section really just deletes the user from the active_tracking slice
	tState := tracking.StateInfo.Active_tracking
	i := 0
	for q, v := range tState {
		if v == uid {
			i = q
		}
	}
	tState = append(tState[:i], tState[i+1:]...)


	delete(tracking.CurrentlyTracking, uid)
	tracking.SaveBotState(uid)	

	return "Successfully ended tracking for: " + uid
}



// Function that handles the request to generate a sentence from another user's markov chain
func generate(input []string, m *discordgo.MessageCreate) string {
	uid := m.Mentions[0].ID
	if !validate(input, m, 1, 1) {
		return "You have to/can only mention 1 person :(."
	}

	if _, ok := tracking.CurrentlyTracking[uid]; ok {
		return "Stop being stupid Botond! Can't generate a sentence for a user that is actively being tracked >:("
	}
	userDict, err := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s.dict", uid)))
	defer userDict.Close()
	if err != nil {
		return "Stop being stupid Botond! You've never ever tracked this person!"
	}
	dictBytes, _ := ioutil.ReadAll(userDict)


	data := internal.FromGOB64(string(dictBytes))
	sentence := strings.Title(strings.Join(data.Generate(), " "))
	
	return sentence + "."
}



// Function that handles the request from a user to change the channel that is actively being tracked from
func setChannel(input []string, m *discordgo.MessageCreate) string {
	if !validate(input, m, 1, 0) {
		return "Invalid format >:("
	}
	channelID := m.ChannelID

	tracking.StateInfo.Channel = channelID
	return "Successfully changed channel to: " + channelID
}