package pipeline

import (
	"os/exec"
	"strings"
	"unicode"
)

// Returns new pipeline for string.
func NewFromString(exp string) *Pipeline {
	cmds := ParseCommand(exp)
	return NewPipeline(cmds...)
}

// Parses command string expression. Returns
// a slice of exec.Cmds.
func ParseCommand(exp string) []*exec.Cmd {
	rawCommands := mapTrim(strings.Split(exp, "|"), strings.TrimSpace)
	commands := make([]*exec.Cmd, len(rawCommands))
	for i, cmd := range rawCommands {
		args := mapTrim(split(cmd), strings.TrimSpace, trimQuotes)
		commands[i] = exec.Command(args[0], args[1:]...)
	}
	return commands
}

// Splits a given if on spaces unless the space
// is surronded by double quotes.
func split(str string) []string {
	isSpliting := true
	return strings.FieldsFunc(str, func(char rune) bool {
		if char == '"' {
			isSpliting = !isSpliting
		}
		return unicode.IsSpace(char) && isSpliting
	})
}

// Any function that takes a string as a parameter
// and returns a string.
type TrimFunc func(string) string

// Iterates over a slice of strings. Applies each TrimFunc
// sequentially to each string of slice. Returns new
// modified slice.
func mapTrim(strs []string, trimFns ...TrimFunc) []string {
	ret := make([]string, len(strs))
	for i, val := range strs {
		rep := val
		for _, trimFn := range trimFns {
			rep = trimFn(rep)
		}
		ret[i] = rep
	}
	return ret
}

// Trims leading or trailing quotes.
func trimQuotes(str string) string {
	return strings.TrimFunc(str, func(char rune) bool {
		switch char {
		case '"', '\'', '`':
			return true
		}
		return false
	})
}
