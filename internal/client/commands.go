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
	"json"
)



// Utility functions
func validate(input []string, m *discordgo.MessageCreate, inputs int, mentions int) bool {
	if len(input) != inputs + 1 || len(m.Mentions) != mentions {
		return false
	}

	return true
}





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


func endtrack(input []string, m *discordgo.MessageCreate) string {
	uid := m.Mentions[0].ID
	if !validate(input, m, 1, 1) {
		return "You have to/can only mention 1 person :(."
	}
	tState := *tracking.StateInfo.Active_tracking

	// Delete from active_tracking
	i := 0
	for q, v := range tState {
		if v == uid {
			i = q
		}
	}
	tState = append(tState[:i], tState[i+1:]...)


	// Read from the currently tracking map and save the data
	dataString := internal.ToGOB64(tracking.CurrentlyTracking[uid])
	userDict, _ := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s.dict", uid)))
	defer userDict.Close()
	ioutil.WriteFile(userDict.Name(), []byte(dataString), 0644)
	delete(tracking.CurrentlyTracking, uid)


	// save the bot state as well
	botStateJson, _ := json.Marshal(tracking.StateInfo)
	f, _ := os.Open(internal.GetAbsFile("stateinfo.json"))
	defer f.Close()
	ioutil.WriteFile(f.Name(), botStateJson)

	return "Successfully ended tracking for: " + uid
}



func setChannel(input []string, m *discordgo.MessageCreate) string {
	if !validate(input, m, 1, 0) {
		return "Invalid format >:("
	}
	channelID := m.ChannelID

	tracking.StateInfo.Channel = channelID
	return "Successfully changed channel to: " + channelID
}