package prompt

import "fmt"

func helpSessions() {
	fmt.Println(`sessions
  Output:
  1 => windows 192.168.0.13:45722
  2 => linux 192.168.0.13:45722`)
}
func helpConnect() {
	fmt.Println(`connect <id>`)
}

func helpOptions() {
	fmt.Println(`debug: enable/disable debug output
port: update listener port
raw: true (default) the terminal will be set to raw mode. Otherwise will stay in cooked mode
bufferTimer: Time to wait to clear terminal's buffer when executing command (in ms). Mostly use for Windows.
conpty
	disableconpty: In the case of conpty causing issue on your reverse shell you could disable it but your reverse shell will not be interactive
	onlywebserver: if you have already a powershell commande execution you can use this option to serve the ConPty scripts and get your interactive reverse shell`)
}

func help() {
	fmt.Println("sessions: Manage route to socks servers")
	fmt.Println("connect: Manage route to socks servers")
	fmt.Println("menu: Start the reverse shell menu")
	fmt.Println("options: Manage current options (default print them)")
	fmt.Println("start: Start the listener")
	fmt.Println("stop: Stop the listener")
	fmt.Println("restart: Restart the listener")
	fmt.Println("help: help command")
}
