package gsheets

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	var delimiter string
	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Create a spread sheet",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return save(delimiter)
		},
	}
	// TODO: Add mysql flag to parse mysql result
	saveCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "delemiter to parse the input")
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
func save(delimiter string) error {
	ctx := context.Background()
	client, err := getClient(ctx, getAuthConfig())
	if err != nil {
		return fmt.Errorf("Unable to get client with access token.\n%v\n", err)
	}
	sheetService, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets Client.\n%v\n", err)
	}

	// TODO: Refactoring
	var rows []*sheets.RowData
	for _, line := range stdin {
		var cells []*sheets.CellData
		for _, c := range regexp.MustCompile(delimiter).Split(line, -1) {
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
		return fmt.Errorf("Unable to create a new sheet.\n%v\n", err)
	}
	cmd.OpenBrowser(resp.SpreadsheetUrl)
	return nil
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tok, err := getTokenFromFile()
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, tok), err
}
