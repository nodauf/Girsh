package prompt

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

var commands = []prompt.Suggest{
	{Text: "sessions", Description: "Print actives sessions"},
	{Text: "connect", Description: "Connect to an active session"},
	{Text: "menu", Description: "Start the reverse shell menu"},
	{Text: "options", Description: "Manage current options (default print them)"},
	{Text: "start", Description: "Start the listener"},
	{Text: "stop", Description: "Stop the listener"},
	{Text: "restart", Description: "Restart the listener"},
	{Text: "help", Description: "Help menu"},

	{Text: "exit", Description: "Exit this program"},
}

// Options subcommand

var optionsSubCommand = []prompt.Suggest{
	{Text: "debug", Description: "Manage debug option"},
	{Text: "port", Description: "Manage port listener option"},
	{Text: "conpty", Description: "Manage conpty option"},
	{Text: "raw", Description: "Manage the activation of raw terminal"},
	{Text: "timerBuffer", Description: "Time to wait to clear terminal's buffer when executing command (in ms)"},
}

var conptySubCommand = []prompt.Suggest{
	{Text: "disableconpty", Description: "Disable the use of ConPty, the reverse shell will not be interactive"},
	{Text: "onlywebserver", Description: "Will not send the powershell command. A webserver will start to serve Invoke-ConPtyShell.ps1 on every endpoint"},
}

// Help subcommand

var helpSubCommand = []prompt.Suggest{
	{Text: "sessions", Description: "Help about sessions"},
	{Text: "connect", Description: "Help about connect command"},
	{Text: "options", Description: "Help about options"},
}

func complete(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")

	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}
	first := strings.ToLower(args[0])
	switch first {
	case "help":
		second := args[1]
		if len(args) == 2 {
			return prompt.FilterHasPrefix(helpSubCommand, second, true)
		}
	case "options":
		second := strings.ToLower(args[1])
		if len(args) == 2 {
			return prompt.FilterHasPrefix(optionsSubCommand, second, true)
		} else if len(args) > 2 {
			switch second {
			case "conpty":
				third := args[2]
				if len(args) == 3 {
					return prompt.FilterHasPrefix(conptySubCommand, third, true)
				}
			}
		}
	}
	return []prompt.Suggest{}
}
