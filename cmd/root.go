package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tong",
	Short: "Tong is very useful",
	Long: "Tong is a command line application to simply usual long commands",
	Run: func(cmd *cobra.Command, args []string) { },
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
