package gsheets

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/innossh/tong/cli/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const (
	TokenCacheDir  = ".tong/credentials"
	TokenCacheName = "accounts.google.com-oauth2-token.json"
)

func NewAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Retrieve authorization token for Gogle Sheets API",
		Run: func(cmd *cobra.Command, args []string) {
			auth(getConfig())
		},
	}
	return authCmd
}

func auth(config *oauth2.Config) {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
		return
	}
	_, err = tokenFromFile(cacheFile)
	if err == nil {
		return
	}
	tok, err := getTokenFromWeb(config)
	if err != nil {
		log.Fatalf("Failed to get access token. %v", err)
		return
	}
	saveToken(cacheFile, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	l, err := net.Listen("tcp", TmpLocalServer)
	if err != nil {
		log.Fatalf("failed listen %v", err)
		return nil, err
	}
	defer l.Close()

	stateBytes := make([]byte, 16)
	_, err = rand.Read(stateBytes)
	if err != nil {
		log.Fatalf("failed stateBytes %v", err)
		return nil, err
	}

	state := fmt.Sprintf("%x", stateBytes)
	cmd.OpenBrowser(config.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "token")))
	quit := make(chan *oauth2.Token)
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			w.Write([]byte(`<script>location.href = "/close?" + location.hash.substring(1);</script>`))
		} else {
			w.Write([]byte(`<script>window.open("about:blank","_self").close()</script>`))
			w.(http.Flusher).Flush()
			quit <- &oauth2.Token{
				AccessToken: req.URL.Query().Get("access_token"),
				TokenType:   req.URL.Query().Get("token_type"),
			}
		}
	}))

	return <-quit, nil
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, TokenCacheDir)
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape(TokenCacheName)), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
