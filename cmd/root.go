package cmd

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
)

var stdin string
var RootCmd = &cobra.Command{
	Use:   "tong",
	Short: "Tong is very useful",
	Long:  "Tong is a command line application to simply usual long commands",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if terminal.IsTerminal(syscall.Stdin) {
			return errors.New("Unable to read stdin")
		}
		b, _ := ioutil.ReadAll(os.Stdin)
		stdin = string(b)
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize()
}
