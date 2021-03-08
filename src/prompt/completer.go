package prompt

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

var commands = []prompt.Suggest{
	{Text: "sessions", Description: "Print actives sessions"},
	{Text: "connect", Description: "Connect to an active session"},
	{Text: "options", Description: "Manage current options (default print them)"},
	{Text: "start", Description: "Start the listener"},
	{Text: "stop", Description: "Stop the listener"},
	{Text: "restart", Description: "Restart the listener"},
	{Text: "help", Description: "Help menu"},

	{Text: "exit", Description: "Exit this program"},
}

var optionsSubCommand = []prompt.Suggest{
	{Text: "debug", Description: "Manage debug option"},
	{Text: "port", Description: "Manage port listener option"},
	{Text: "conpty", Description: "Manage conpty option"},
}

var helpSubCommand = []prompt.Suggest{
	{Text: "sessions", Description: "Print actives sessions"},
}

func complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}
	first := args[0]
	switch first {
	case "help":
		second := args[1]
		if len(args) == 2 {
			return prompt.FilterHasPrefix(helpSubCommand, second, true)
		}
	case "options":
		second := args[1]
		if len(args) == 2 {
			return prompt.FilterHasPrefix(optionsSubCommand, second, true)
		}
	}
	return []prompt.Suggest{}
}
