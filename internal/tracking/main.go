package tracking

import (
	"io/ioutil"
	"bendy-bot/internal/markov"
	"bendy-bot/internal"
	"github.com/bwmarrin/discordgo"
	"os"
	"encoding/json"
	"fmt"
)



type BotStateInfo struct {
	Channel string
	Active_tracking []string
}
var StateInfo BotStateInfo
var CurrentlyTracking map[string]markov.Markov



func init() {
	// Read the config files into these variables
	jsonFile, _ := os.Open("stateinfo.json")
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, &StateInfo)

	// Load up the data from the data storage for the users that are currently within the active_tracking section
	for _, user := range StateInfo.Active_tracking {
		// Open the file the holds this users dictionary
		dictBytes := internal.OpenFileFromStore(user)
	
		data := internal.FromGOB64(string(dictBytes))
		CurrentlyTracking[user] = data
	}
}



// Function that just saves the current state of the bot and the specific user that must be saved
func SaveBotState(user string) {
	// if a user string is provided then save that user's state
	if user != "" {
		dataString := internal.ToGOB64(CurrentlyTracking[user])
		userDict, _ := os.Open(internal.GetAbsFile(fmt.Sprintf("data/%s.dict", user)))
		defer userDict.Close()
		ioutil.WriteFile(userDict.Name(), []byte(dataString), 0644)
	}

	// Save the current bot's state as well
	botStateJson, _ := json.Marshal(StateInfo)
	f, _ := os.Open(internal.GetAbsFile("stateinfo.json"))
	defer f.Close()
	ioutil.WriteFile(f.Name(), botStateJson, 0644)
}



// Also does what it says :P Handles and incoming message for storage
func HandleIncomingMessage(m *discordgo.MessageCreate) {
	if m.ChannelID == StateInfo.Channel {
		// Parse the message on to be tracked
		if _, ok := CurrentlyTracking[m.Author.ID]; ok {
			chain := CurrentlyTracking[m.Author.ID]
			go chain.Parse(m.Content)
		}
	}
}