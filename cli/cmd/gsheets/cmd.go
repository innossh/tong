package gsheets

import (
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	DriveScope        = "https://www.googleapis.com/auth/drive"
	SpreadsheetsScope = "https://www.googleapis.com/auth/spreadsheets"
	TmpLocalServer    = "localhost:10080"
	ClientId          = "1008608304541-sn085h0v9lg987c114skdf6s8km9rq4i.apps.googleusercontent.com"
)

func NewGsheetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gsheets",
		Short: "Google spread sheets",
	}

	cmd.AddCommand(
		NewAuthCmd(),
		NewSaveCmd(),
	)
	return cmd
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes: []string{
			DriveScope,
			SpreadsheetsScope,
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		ClientID:    ClientId,
		RedirectURL: "http://" + TmpLocalServer,
	}
}
