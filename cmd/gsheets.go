package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(gSheetsCmd)
}

var gSheetsCmd = &cobra.Command{
	Use:   "gsheets",
	Short: "Google spread sheets",
	Long:  `Create a spread sheet`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: create a spread sheet
		fmt.Println(stdin)
	},
}
