package terminal

import (
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/op/go-logging"
)

// Terminal object
type Terminal struct {
	OS       string
	Con      net.Conn
	Options  Options
	rows     string
	cols     string
	spawnTTY string
	mutex    sync.Mutex
	Listener net.Listener
	Log      *logging.Logger
}

// Options type to manage the terminal's options, also use by the session
type Options struct {
	Port          int
	Debug         bool
	DisableConPTY bool
	OnlyWebserver bool
	Raw           bool
	TimerBuffer   int
}

// New will initialize the logging configuration and start the listener and wait for client.
func (terminal *Terminal) New() error {
	//terminal.logging()
	return listenAndAcceptConnection(terminal)
}

/*func (terminal *Terminal) logging() {
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

}*/

// GetOS send the command whoami and parse the result. The windows format is COMPUTERNAME\username
func (terminal *Terminal) GetOS() {

	// Wait a bit for the connection to be ready. Windows connection can take some time to be ready
	time.Sleep(2000 * time.Millisecond)
	//Use whoami command and parse outptut
	output := terminal.execute("whoami", []byte{promptLinux1, promptLinux2, promptWindows1})
	if strings.Contains(string(output), "\\") {
		terminal.OS = "windows"
	} else {
		terminal.OS = "linux"
	}
	//fmt.Println(string(execute("env", dst)))

}

// PrepareShell manage the stty raw, spawn the tty for linux and use the ConPTY for windows
func (terminal *Terminal) PrepareShell() error {

	if terminal.OS == "linux" {
		// If the main binary is running on linux
		if runtime.GOOS == "linux" {
			terminal.setSpawnTTY()
		}
		terminal.interactiveReverseShellLinux()

	} else if terminal.OS == "windows" {

		if !terminal.Options.DisableConPTY {
			terminal.interactiveReverseShellWindows()
			terminal.Log.Debug("Starting http server on " + terminal.Con.LocalAddr().String())
			err := terminal.serveHTTPRevShellPowershell()
			if err != nil {
				return err
			}

			if err := listenAndAcceptConnection(terminal); err != nil {
				terminal.Log.Error("Error while listening or accepting a connection " + err.Error())
				return err
			}
		}

	} else if terminal.Options.OnlyWebserver {
		terminal.Log.Debug("Starting http server on 0.0.0.0")
		err := terminal.serveHTTPRevShellPowershell()
		if err != nil {
			return err
		}

		if err = listenAndAcceptConnection(terminal); err != nil {
			terminal.Log.Error("Error while listening or accepting a connection " + err.Error())
			return err
		}
	}
	return nil
	//terminal.Connect()
}

// Connect to a remote stdin and stdout to a net.Conn
func (terminal *Terminal) Connect() int {
	defer func() {
		if runtime.GOOS == "linux" {
			terminal.sttyRawEcho("disable")
		}
	}()
	// If the main binary is running on linux
	if runtime.GOOS == "linux" {
		// Set the terminal to raw mode
		terminal.sttyRawEcho("enable")
	}
	// The terminal is natively in raw mode with go-prompt, we need to disable the raw mode when this is not necessary
	// If the client is windows OS and we disable ConPTY, the raw mode is not needed
	if !terminal.Options.Raw || (terminal.OS == "windows" && terminal.Options.DisableConPTY) {
		terminal.sttyRawEcho("disable")
	}
	var chanToStdout = make(chan int)
	var chanToRemote = make(chan int)
	var kill = make(chan bool)

	terminal.streamCopy(terminal.Con, os.Stdout, false, chanToStdout, kill)
	terminal.streamCopy(os.Stdin, terminal.Con, true, chanToRemote, kill)

	select {
	case status := <-chanToStdout:
		kill <- true
		if status == 0 {
			terminal.Log.Debug("Remote connection is closed")
		} else if status == 1 {
			terminal.Log.Debug("Remote connection is backgrounded")
		}
		return status

	case status := <-chanToRemote:
		kill <- true
		if status == 0 {
			terminal.Log.Debug("Local program is terminated")

		} else if status == 1 {
			terminal.Log.Debug("Connection is backgrounded")
		}
		return status
	}

}
