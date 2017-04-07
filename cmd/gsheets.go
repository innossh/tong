package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
)

var clientSecretJson string

func init() {
	gSheetsCmd.PersistentFlags().StringVarP(&clientSecretJson, "client-secret-json", "c", "./client_secret.json", "client secret file (default is ./client_secret.json)")
	RootCmd.AddCommand(gSheetsCmd)
}

var gSheetsCmd = &cobra.Command{
	Use:   "gsheets",
	Short: "Google spread sheets",
	Long:  `Create a spread sheet`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
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

		cell := &sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: stdin,
			},
		}
		cells := []*sheets.CellData{cell}
		row := &sheets.RowData{
			Values: cells,
		}
		rows := []*sheets.RowData{row}
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

		fmt.Printf("%#v\n", resp)
	},
}
