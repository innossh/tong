package gsheets

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"syscall"

	"github.com/innossh/tong/cli/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func NewSaveCmd(clientSecretJson string) *cobra.Command {
	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Create a spread sheet",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validate()
		},
		Run: func(cmd *cobra.Command, args []string) {
			save(clientSecretJson)
		},
	}
	return saveCmd
}

var stdin []string

func validate() error {
	if terminal.IsTerminal(syscall.Stdin) {
		return errors.New("Unable to read stdin")
	}
	b, _ := ioutil.ReadAll(os.Stdin)
	stdin = regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(b), -1)
	return nil
}

func save(clientSecretJson string) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecretJson)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, DriveScope, SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

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
