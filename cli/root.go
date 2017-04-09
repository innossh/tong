package cmd

import (
	"github.com/spf13/cobra"
	"github.com/innossh/tong/cli/cmd"
	"github.com/innossh/tong/cli/cmd/gsheets"
)

func SetupRootCmd(rootCmd *cobra.Command) {
	//rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	//rootCmd.PersistentFlags().MarkShorthandDeprecated("help", "Please use --help")
}

func AddCmds(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		cmd.NewVersionCmd(),

		gsheets.NewGsheetsCmd(),
	)
}
