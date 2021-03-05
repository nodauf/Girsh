package terminal

import (
	"net"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/op/go-logging"
)

// Terminal object
type Terminal struct {
	OS      string
	Con     net.Conn
	Options struct {
		Port          int
		Debug         bool
		DisableConPTY bool
	}
	rows     string
	cols     string
	spawnTTY string
	mutex    sync.Mutex
	listener net.Listener
	log      *logging.Logger
}

// New will initialize the logging configuration and start the listener and wait for client.
func (terminal *Terminal) New() {
	terminal.logging()
	listenAndAcceptConnection(terminal)
}

func (terminal *Terminal) logging() {
	terminal.log = logging.MustGetLogger("nc-shell")

	logger := logging.NewLogBackend(os.Stderr, "", 0)
	var loggerLeveled logging.LeveledBackend

	if terminal.Options.Debug {
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

// GetOS send the command whoami and parse the result. The windows format is COMPUTERNAME\username
func (terminal *Terminal) GetOS() {
	//Use env ou set command and parse outptut
	output := terminal.execute("whoami")
	if strings.Contains(string(output), "\\") {
		terminal.OS = "windows"
	} else {
		terminal.OS = "linux"
	}
	//fmt.Println(string(execute("env", dst)))

}

// Shell manage the stty raw, spawn the tty for linux and use the ConPTY for windows and connect the stdin and stdout with the con
func (terminal *Terminal) Shell() {
	defer func() {
		if runtime.GOOS == "linux" {
			terminal.sttyRawEcho("disable")
		}
	}()
	if terminal.OS == "linux" {
		// If the main binary is running on linux
		if runtime.GOOS == "linux" {
			// Set the terminal to raw mode
			terminal.sttyRawEcho("enable")
			terminal.setSpawnTTY()
		}
		terminal.interactiveReverseShellLinux()

	} else if terminal.OS == "windows" {

		if !terminal.Options.DisableConPTY {
			// If the main binary is running on linux
			if runtime.GOOS == "linux" {
				// Set the terminal to raw mode
				terminal.sttyRawEcho("enable")
			}
			terminal.interactiveReverseShellWindows()
			terminal.log.Debug("Starting http server on " + terminal.Con.LocalAddr().String())
			terminal.serveHTTPRevShellPowershell()
			listenAndAcceptConnection(terminal)
		}

	}
	chanToStdout := terminal.streamCopy(terminal.Con, os.Stdout, false)
	chanToRemote := terminal.streamCopy(os.Stdin, terminal.Con, true)
	select {
	case <-chanToStdout:
		terminal.log.Debug("Remote connection is closed")

	case <-chanToRemote:
		terminal.log.Debug("Local program is terminated")

	}

}
