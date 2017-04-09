package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Tong",
		Long:  `All software has versions. This is Tong's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Tong v0.2.0")
		},
	}
	return versionCmd
}
