package main

import (
	"fmt"
	"os"

	"github.com/innossh/tong/cli"
	"github.com/spf13/cobra"
)

func main() {
	tongCmd := NewTongCmd()
	if err := tongCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewTongCmd() *cobra.Command {
	tongCmd := &cobra.Command{
		Use:   "tong",
		Short: "Tong is very useful",
		Long:  "Tong is a command line application to simplify usual long commands",
	}
	cmd.SetupRootCmd(tongCmd)
	cmd.AddCmds(tongCmd)
	return tongCmd
}
