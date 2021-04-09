package sessions

import (
	"nc-shell/src/terminal"
	"strconv"

	logging "github.com/op/go-logging"
)

var sessions = make(map[int]*terminal.Terminal)
var lastSessionID = 0
var log *logging.Logger
var term = &terminal.Terminal{}

func sessionIDExists(idString string) (bool, int) {
	id, err := strconv.Atoi(idString)
	if _, ok := sessions[id]; ok && err == nil {
		return true, id
	}
	return false, 0
}

func newTerminals() {
	for {
		term = &terminal.Terminal{}
		term.Options = OptionsSession
		term.Log = log
		if !term.Options.OnlyWebserver {
			err := term.New()
			// Exit the function if there is an error
			if err != nil {
				// Destroy term
				term = &terminal.Terminal{}
				break
			}
			term.GetOS()
		} else {
			term.Log.Debug("Skipping first stage as the option OnlyWebserver is set to true")
		}
		err := term.PrepareShell()
		// Exit the function if there is an error
		if err != nil {
			// Destroy term
			term = &terminal.Terminal{}
			break
		}
		sessionID := lastSessionID + 1
		sessions[sessionID] = term
		term.Log.Info("Session " + strconv.Itoa(sessionID) + " (" + term.OS + ") available")
		lastSessionID = sessionID
	}
}
