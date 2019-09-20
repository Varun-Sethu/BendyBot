package bot

import (
	"io/ioutil"
	"bendy-bot/internal/markov"
	"bendy-bot/internal"
	"github.com/bwmarrin/discordgo"
	"os"
	"encoding/json"
	"strings"
	"regexp"
	"fmt"
	"errors"
)



type botStateInfo struct {
	Channel string
	Active_tracking []string
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
	for _, user := range stateInfo.Active_tracking {
		// Open the file the holds this users dictionary
		dictBytes := internal.OpenFileFromStore(user)	
		currentlyTracking[user] = markov.Build(string(dictBytes))
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

// Function to save the state of a user
func SaveUserState(user string) {
	dataString := currentlyTracking[user].ToJSON()
	userDict, _ := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s.json", user)))
	defer userDict.Close()
	ioutil.WriteFile(userDict.Name(), []byte(dataString), 0644)
}



// Also does what it says :P Handles and incoming message for storage
func HandleIncomingMessage(m *discordgo.MessageCreate) {
	if m.ChannelID == stateInfo.Channel {
		// Parse the message on to be tracked
		if _, ok := currentlyTracking[m.Author.ID]; ok {
			go currentlyTracking[m.Author.ID].Parse(m.Content)
		}
	}
}





// Utility functions that are used by the client and commands
// BeginTrackingUser is a function that implements a tracker for a specific user
func BeginTrackingUser(uid string) error {
	if _, ok := currentlyTracking[uid]; ok {
		return errors.New("Stop being stupid Botond! User is already being tracked >:(")
	}

	// Update the tracking state and read from their storage file
	stateInfo.Active_tracking = append(stateInfo.Active_tracking, uid)
	dictBytes := internal.OpenFileFromStore(uid)

	if len(dictBytes) == 0 {
		// Build a new markov chain if nothing exists in storage for this user
		currentlyTracking[uid] = markov.Build()
	} else {
		// Interpret and parse their data
		currentlyTracking[uid] = markov.Build(string(dictBytes))
	}
	SaveBotState()
	SaveUserState(uid)

	return nil
}



// Function to end the tracking of a specific user
func EndTrackingUser(uid string) {
	// This section really just deletes the user from the active_tracking slice
	i := 0
	for q, v := range stateInfo.Active_tracking {
		if v == uid {
			i = q
		}
	}
	stateInfo.Active_tracking = append(stateInfo.Active_tracking[:i], stateInfo.Active_tracking[i+1:]...)

	SaveBotState()
	SaveUserState(uid)
	delete(currentlyTracking, uid)
}



// Function to generate some text for an individual
func GenerateSentenceForUser(uid string) (string, error) {
	if _, ok := currentlyTracking[uid]; ok {
		return "", errors.New("Stop being stupid Botond! Can't generate a sentence for a user that is actively being tracked >:(")
	}
	userDict, err := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s.json", uid)))
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
func SetTrackingChannel(channelID string) {
	stateInfo.Channel = channelID
	SaveBotState()
}