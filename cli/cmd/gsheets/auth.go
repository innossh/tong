package gsheets

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/innossh/tong/cli/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	TokenCacheDir  = ".tong/credentials"
	TokenCacheName = "accounts.google.com-oauth2-token.json"
)

func NewAuthCmd() *cobra.Command {
	var force bool
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Retrieve authorization token for Gogle Sheets API",
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth(getAuthConfig(), force)
		},
	}
	authCmd.Flags().BoolVarP(&force, "--force", "f", false, "recreate access token even if the credential file exists")
	return authCmd
}

func auth(config *oauth2.Config, force bool) error {
	if !force {
		_, err := getTokenFromFile()
		if err == nil {
			fmt.Println("This client has already been authorized.")
			return nil
		}
	}

	t, err := getTokenFromWeb(config)
	if err != nil {
		return err
	}
	return saveToken(t)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	l, err := net.Listen("tcp", TmpLocalServer)
	if err != nil {
		return nil, fmt.Errorf("Unable to listen a local port, %s\n%v\n", TmpLocalServer, err)
	}
	defer l.Close()

	stateBytes := make([]byte, 16)
	_, err = rand.Read(stateBytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to make state bytes\n%v\n", err)
	}

	state := fmt.Sprintf("%x", stateBytes)
	err = cmd.OpenBrowser(config.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "token")))
	if err != nil {
		return nil, fmt.Errorf("Unable to open web browser\n%v\n", err)
	}

	quit := make(chan *oauth2.Token)
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			w.Write([]byte(`<script>location.href = "/close?" + location.hash.substring(1);</script>`))
		} else {
			w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
			w.(http.Flusher).Flush()
			expiry, err := strconv.Atoi(req.URL.Query().Get("expires_in"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get a valid expiry\n%v\n", err)
				expiry = 0
			}
			quit <- &oauth2.Token{
				AccessToken: req.URL.Query().Get("access_token"),
				TokenType:   req.URL.Query().Get("token_type"),
				Expiry:      time.Now().Add(time.Duration(expiry) * time.Second),
			}
		}
	}))

	return <-quit, nil
}

// saveToken creates a file and store the
// token in it.
func saveToken(token *oauth2.Token) error {
	file, err := tokenCacheFile()
	if err != nil {
		return fmt.Errorf("Unable to get path to cached credential file.\n%v\n", err)
	}

	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
