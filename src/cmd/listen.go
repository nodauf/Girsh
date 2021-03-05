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
		term.Options.DisableConPTY = disableConPTY
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
