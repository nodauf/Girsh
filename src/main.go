package main

import (
	cmd "nc-shell/src/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cmd.Execute()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}
