package bot

import (
	"bendy-bot/internal"
	"bendy-bot/internal/markov"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)



type botStateInfo struct {
	Channel map[string]string
	Active_tracking map[string][]string
}
var stateInfo botStateInfo
var currentlyTracking map[string]*markov.Markov



func init() {
	// Read the config files into these variables
	jsonFile, _ := os.Open(internal.GetAbsFile("stateinfo.json"))
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, &stateInfo)

	currentlyTracking = make(map[string]*markov.Markov)

	// Load up the data from the data storage for the users that are currently within the active_tracking section
	for server, _ := range stateInfo.Active_tracking {
		for _, user := range stateInfo.Active_tracking[server] {
			// Open the file the holds this users dictionary
			dictBytes := internal.OpenFileFromStore(fmt.Sprintf("%s-%s", user, server))	
			currentlyTracking[fmt.Sprintf("%s-%s", user, server)] = markov.Build(string(dictBytes))
		}
	}
}



// Function that just saves the current state of the bot and the specific user that must be saved
func SaveBotState() {
	// Save the current bot's state as well
	botStateJson, _ := json.Marshal(stateInfo)
	err := ioutil.WriteFile(internal.GetAbsFile("stateinfo.json"), botStateJson, 0644)
	if err != nil {
		panic(err)
	}
}

// Function to save the state of a user, for privacy reasons a guild id must be provided and only that guild can access data from that id
func SaveUserState(guildID string, user string) {
	dataString := currentlyTracking[fmt.Sprintf("%s-%s", user, guildID)].ToJSON()
	userDict, _ := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s-%s.json", user, guildID)))
	defer userDict.Close()
	ioutil.WriteFile(userDict.Name(), []byte(dataString), 0644)
}



// Also does what it says :P Handles and incoming message for storage
func HandleIncomingMessage(m *discordgo.MessageCreate) {
	if m.ChannelID == stateInfo.Channel[m.GuildID] {
		// Parse the message on to be tracked
		if _, ok := currentlyTracking[fmt.Sprintf("%s-%s", m.Author.ID, m.GuildID)]; ok {
			go currentlyTracking[fmt.Sprintf("%s-%s", m.Author.ID, m.GuildID)].Parse(m.Content)
		}
	}
}





// Utility functions that are used by the client and commands
// BeginTrackingUser is a function that implements a tracker for a specific user
func BeginTrackingUser(guildID string, uid string) error {
	if _, ok := currentlyTracking[fmt.Sprintf("%s-%s", uid, guildID)]; ok {
		return errors.New("Stop being stupid Botond! User is already being tracked >:(")
	}

	// Update the tracking state and read from their storage file
	stateInfo.Active_tracking[guildID] = append(stateInfo.Active_tracking[guildID], uid)
	dictBytes := internal.OpenFileFromStore(fmt.Sprintf("%s-%s", uid, guildID))

	if len(dictBytes) == 0 {
		// Build a new markov chain if nothing exists in storage for this user
		currentlyTracking[fmt.Sprintf("%s-%s", uid, guildID)] = markov.Build()
	} else {
		// Interpret and parse their data
		currentlyTracking[fmt.Sprintf("%s-%s", uid, guildID)] = markov.Build(string(dictBytes))
	}
	SaveBotState()
	SaveUserState(guildID, uid)

	return nil
}



// Function to end the tracking of a specific user
func EndTrackingUser(guildID string, uid string) {
	// This section really just deletes the user from the active_tracking slice
	i := 0
	for q, v := range stateInfo.Active_tracking[guildID] {
		if v == uid {
			i = q
		}
	}
	stateInfo.Active_tracking[guildID] = append(stateInfo.Active_tracking[guildID][:i], stateInfo.Active_tracking[guildID][i+1:]...)

	SaveBotState()
	SaveUserState(guildID, uid)
	delete(currentlyTracking, fmt.Sprintf("%s-%s", uid, guildID))
}



// Function to generate some text for an individual
func GenerateSentenceForUser(guildID string, uid string) (string, error) {
	if _, ok := currentlyTracking[uid]; ok {
		return "", errors.New("Stop being stupid Botond! Can't generate a sentence for a user that is actively being tracked >:(")
	}
	userDict, err := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s-%s.json", uid, guildID)))
	defer userDict.Close()
	if err != nil {
		return "", errors.New("Stop being stupid Botond! You've never ever tracked this person!")
	}
	dictBytes, _ := ioutil.ReadAll(userDict)
	// Attain the chain that is connected to that user
	chain := markov.Build(string(dictBytes))
	words := chain.Generate()

	// Create an actual sentence
	words[0] = strings.Title(words[0])
	sentence := ""

	// Actually construct a sentnece with proper grammar
	for i, word := range words {
		// Determine if they are just punctuation or an actual word
		if mathced, _ := regexp.Match(`[].,!?;]`, []byte(word)); mathced {
			sentence = sentence + word
		} else {
			sentence = sentence + " " + word
			// Add a full stop if necessary
			if i == len(words) - 1 {
				sentence = sentence + "."
			}
		}
	}


	return sentence, nil
}



// Function to set the current tracking channel for a user
func SetTrackingChannel(guildID string, channelID string) {
	stateInfo.Channel[guildID] = channelID
	SaveBotState()
}