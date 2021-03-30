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
