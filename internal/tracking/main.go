package tracking

import (
	"io/ioutil"
	"bendy-bot/internal/markov"
	"bendy-bot/internal"
	"os"
	"encoding/json"
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


func Store(uid string, sentence string) {
	chain := CurrentlyTracking[uid]
	go chain.Parse(sentence)
} 