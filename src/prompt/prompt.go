package prompt

import (
	"fmt"
	"nc-shell/src/menu"
	"nc-shell/src/sessions"
	"os"
	"os/exec"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

func handleExit() {
	// workaround for the bug https://github.com/c-bata/go-prompt/issues/147
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func executor(in string) {
	command := strings.Split(in, " ")
	first := command[0]
	switch strings.ToLower(first) {
	case "help":
		if len(command) > 1 {
			second := command[1]
			switch strings.ToLower(second) {
			case "sessions":
				helpSessions()

			case "connect":
				helpConnect()
			}
		} else {
			help()
		}
	case "menu":
		menu.Menu(sessions.OptionsSession.Port)
	case "sessions":
		sessions.PrintSessions()
	case "connect":
		sessions.Connect(command[1])
	case "start":
		sessions.Start()
	case "stop":
		sessions.Stop()
	case "restart":
		sessions.Restart()
	case "options":
		if len(command) > 1 {
			second := command[1]
			switch strings.ToLower(second) {
			case "debug":
				if len(command) > 2 {
					sessions.SetDebug(command[2])
				} else {
					sessions.PrintDebugOptions()
				}
			case "port":
				if len(command) > 2 {
					sessions.SetPort(command[2])
				} else {
					sessions.PrintPortOptions()
				}
			case "disableconpty":
				if len(command) > 2 {
					sessions.SetDisableConPTY(command[2])
				} else {
					sessions.PrintDisableConPTYOptions()
				}
			}
		} else {
			sessions.PrintOptions()
		}
	case "exit":
		os.Exit(0)
	// Not doing anything for just a new line
	case "":
	default:
		fmt.Println("Invalid command")
	}
}

// Prompt run the custom prompt to manage sessions
func Prompt() {
	defer handleExit()
	sessions.Logger()
	sessions.Start()
	p := prompt.New(
		executor,
		complete,
		prompt.OptionPrefix("Girsh> "),
		prompt.OptionPrefixTextColor(prompt.Red),
		prompt.OptionTitle("Girsh"),
	)
	p.Run()
}
