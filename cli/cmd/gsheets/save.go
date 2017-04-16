package gsheets

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"syscall"

	"github.com/innossh/tong/cli/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/sheets/v4"
)

func NewSaveCmd() *cobra.Command {
	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Create a spread sheet",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validate()
		},
		Run: func(cmd *cobra.Command, args []string) {
			save()
		},
	}
	return saveCmd
}

var stdin []string

// validate validates stdin
func validate() error {
	if terminal.IsTerminal(syscall.Stdin) {
		return errors.New("Unable to read stdin")
	}
	b, _ := ioutil.ReadAll(os.Stdin)
	stdin = regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(b), -1)
	return nil
}

// save creates a new sheet
func save() {
	ctx := context.Background()
	client := getClient(ctx, getConfig())
	sheetService, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	var rows []*sheets.RowData
	for _, line := range stdin {
		var cells []*sheets.CellData
		for _, c := range regexp.MustCompile(",").Split(line, -1) {
			cell := &sheets.CellData{
				UserEnteredValue: &sheets.ExtendedValue{
					StringValue: c,
				},
			}
			cells = append(cells, cell)
		}
		row := &sheets.RowData{
			Values: cells,
		}
		rows = append(rows, row)
	}
	grid := &sheets.GridData{
		RowData: rows,
	}
	grids := []*sheets.GridData{grid}
	sheet := &sheets.Sheet{
		Data: grids,
	}
	rb := &sheets.Spreadsheet{
		Sheets: []*sheets.Sheet{sheet},
	}

	resp, err := sheetService.Spreadsheets.Create(rb).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to create a new sheet. %v", err)
	}
	cmd.OpenBrowser(resp.SpreadsheetUrl)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		log.Fatalf("Failed to get access token. %v", err)
	}
	return config.Client(ctx, tok)
}
