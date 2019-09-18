package client

import (
	"regexp"
	"errors"
	"strings"
	"github.com/bwmarrin/discordgo"
)




// function to parse the given input and turn it into a series of tokens
func Parse(m *discordgo.MessageCreate, command string) ([]string, error) {
	// Split the input into a series of tokens separated by whitespace
	construction := regexp.MustCompile(`[\s]+`).Split(command[1:], -1)
	// all the possible commands
	commandList := map[string]int{"track": 0, "generate": 0, "endtrack": 0, "set-channel": 0}

	var output []string

	// Determine if the first token is an actual command and if it ain't then chuck an error
	if _, exists  := commandList[construction[0]]; exists{
		output = []string{construction[0]}
	} else {
		return nil, errors.New("Looks like that was the wrong input. Expect a Piper Warrior flying through your window in a few minutes...")
	}

	// Append the rest of the tokesn which are paramaters
	for _, val := range construction[1:] {
		output = append(output, strings.TrimSpace(val))
	}

	return output, nil
}


func Interpret(input []string, m *discordgo.MessageCreate) string {
	switch input[0] {
	case "track":
		return track(input, m)
	case "generate":
		return generate(input, m)
	case "endtrack":
		return endtrack(input, m)
	case "set-channel":
		return setChannel(input, m)
	}

	return ""
}