package markov

import (
	"math/rand"
	"regexp"
	"sync"
	"strings"
	"time"
)

//https://www.youtube.com/watch?v=IvKynoaD7PA 

// Markov data structure represents our markov chain
type Markov struct {
	initState      string
	terminalState  string
	strings        []string
	edges          map[string][]*string

	// Mutex
	mux 		   sync.Mutex
}





// Helper functions and initializers basically
func Build() Markov {
	return Markov{
		initState: "INIT",
		terminalState: "TERMINAL",
		strings: []string{},
		edges: make(map[string][]*string),
	}
}
func init() {
	rand.Seed(time.Now().UnixNano())
}





// Add node adds an additional state to the markov chain
func (m *Markov) addState(s string) *string {
	m.strings = append(m.strings, s)
	return &m.strings[len(m.strings)-1]
}

// Add edge adds a link between one state and another
func (m *Markov) addEdge(i string, d *string) {
	if _, ok := m.edges[i]; !ok {
		m.edges[i] = []*string{}
	}
	m.edges[i] = append(m.edges[i], d)
}






// Function to generate an arbiatry sentence of words
func (m *Markov) Generate() []string {
	sentence := []string{}
	currState := &m.initState

	// continue adding words to the sentence while the current state is not a terminating state
	for *currState != m.terminalState && *currState != "." {
		// Attain the random state that will be the next word
		nextStateIndex := rand.Intn(len(m.edges[*currState]))
		currState = m.edges[*currState][nextStateIndex]

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

	curr := &m.initState
	for _, word := range words {
		// Generate a state to add to the chain
		stringBlock := strings.ToLower(word)
		newState := m.addState(stringBlock)
		m.addEdge(*curr, newState)
		
		curr = newState
	}

	// tie it to the end state so the generate function knows how to terminate
	m.addEdge(*curr, &m.terminalState)
}
