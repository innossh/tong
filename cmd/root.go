package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"regexp"
	"syscall"
)

var stdin []string
var RootCmd = &cobra.Command{
	Use:   "tong",
	Short: "Tong is very useful",
	Long:  "Tong is a command line application to simply usual long commands",
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

func validate() error {
	if terminal.IsTerminal(syscall.Stdin) {
		return errors.New("Unable to read stdin")
	}
	b, _ := ioutil.ReadAll(os.Stdin)
	stdin = regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(b), -1)
	return nil
}
