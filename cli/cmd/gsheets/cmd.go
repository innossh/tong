package gsheets

import (
	"github.com/spf13/cobra"
)

var clientSecretJson string

func NewGsheetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gsheets",
		Short: "Google spread sheets",
	}
	cmd.PersistentFlags().StringVarP(&clientSecretJson, "client-secret-json", "c", "./client_secret.json", "client secret file (default is ./client_secret.json)")

	cmd.AddCommand(
		NewAuthCmd(clientSecretJson),
		NewSaveCmd(clientSecretJson),
	)
	return cmd
}
