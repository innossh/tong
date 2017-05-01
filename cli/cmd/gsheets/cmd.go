package gsheets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	AuthURL           = "https://accounts.google.com/o/oauth2/auth"
	TokenURL          = "https://accounts.google.com/o/oauth2/token"
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

func getAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes: []string{
			DriveScope,
			SpreadsheetsScope,
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  AuthURL,
			TokenURL: TokenURL,
		},
		ClientID:    ClientId,
		RedirectURL: "http://" + TmpLocalServer,
	}
}

// getTokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func getTokenFromFile() (*oauth2.Token, error) {
	file, err := tokenCacheFile()
	if err != nil {
		return nil, fmt.Errorf("Unable to get path to cached credential file.\n%v\n", err)
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	if t.Expiry.Before(time.Now()) {
		err = errors.New("Token has already been expired.")
	}
	return t, err
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, TokenCacheDir)
	err = os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape(TokenCacheName)), err
}
