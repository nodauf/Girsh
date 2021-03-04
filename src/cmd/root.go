package cmd

import (
	"fmt"
	"nc-shell/src/menu"
	terminal "nc-shell/src/terminal"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var port int
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "girsh",
	Short: "Generate a reverseshell oneliners and listen",
	Long: `Generate a reverseshell oneliners (credits shellerator).
	And listen then run stty raw -echo and send the python command to spawn a tty shell if it's Linux
	or use ConPTY if it's windows`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		port := menu.Menu(port)
		term := &terminal.Terminal{}
		term.Options.Port = port
		term.Options.Debug = debug

		term.New()
		term.GetOS()
		term.Shell()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 1234, "port to listen (default is 1234)")

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug output")

}
