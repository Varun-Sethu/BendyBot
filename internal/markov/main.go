package markov

import (
	"math/rand"
	"regexp"
	"sync"
	"strings"
	"time"
	"encoding/json"
)



// Markov data structure represents our markov chain
type Markov struct {
	InitState      string
	TerminalState  string
	Strings        []string
	Edges          map[string][]*string

	// Mutex
	mux 		   sync.Mutex
}



// Helper functions and initializers basically
func Build(inp ...string) *Markov {
	markov := new(Markov)
	// parse the input as a json strcut if there are any actual inputs
	if len(inp) == 1 {
		json.Unmarshal([]byte(inp[0]), markov)
		return markov
	}
	
	// Return a default markov state
	*markov = Markov{
		InitState: "INIT",
		TerminalState: "TERMINAL",
		Strings: []string{},
		Edges: make(map[string][]*string),
	}
	return markov
}

func init() {
	rand.Seed(time.Now().UnixNano())
}





// Add node adds an additional state to the markov chain
func (m *Markov) addState(s string) *string {
	m.Strings = append(m.Strings, s)
	return &m.Strings[len(m.Strings)-1]
}


// Add edge adds a link between one state and another
func (m *Markov) addEdge(i string, d *string) {
	if _, ok := m.Edges[i]; !ok {
		m.Edges[i] = []*string{}
	}
	m.Edges[i] = append(m.Edges[i], d)
}






// Function to generate an arbiatry sentence of words
func (m *Markov) Generate() []string {
	sentence := []string{}
	currState := &m.InitState

	// continue adding words to the sentence while the current state is not a terminating state
	for *currState != m.TerminalState && *currState != "." {
		// Attain the random state that will be the next word
		nextStateIndex := rand.Intn(len(m.Edges[*currState]))
		currState = m.Edges[*currState][nextStateIndex]

		sentence = append(sentence, *currState)
	}

	// Cut off the final "TERMINAL" text
	return append(sentence[:len(sentence)-1], ".")
}


// Function to parse input data
func (m *Markov) Parse(sentence string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	
	// Split into words and punctuation
	r := regexp.MustCompile(`[\w']+|[.,!?;]`)
	words := r.FindAllString(sentence, -1)

	curr := &m.InitState
	for _, word := range words {
		// Generate a state to add to the chain
		stringBlock := strings.ToLower(word)
		newState := m.addState(stringBlock)
		m.addEdge(*curr, newState)
		
		curr = newState
	}

	// tie it to the end state so the generate function knows how to terminate
	m.addEdge(*curr, &m.TerminalState)
}





// Function that turns the markov chain to json
func (m *Markov) ToJSON() string {
	json, _ := json.Marshal(m)
	return string(json)
}