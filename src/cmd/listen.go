/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	terminal "nc-shell/src/terminal"

	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen and spawn a fully interactive for windows and linux client",
	Long: `Listen and run stty raw -echo and send the python command to spawn a tty shell if it's Linux
	or use ConPTY if it's windows`,
	Run: func(cmd *cobra.Command, args []string) {

		//network.ListenAndAcceptConnection(port)
		term := &terminal.Terminal{}
		term.Options.Port = port
		term.Options.Debug = debug
		term.New()
		term.GetOS()
		term.Shell()

	},
}

func init() {
	rootCmd.AddCommand(listenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
