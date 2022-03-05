package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/regmicmahesh/merosharemorelikeidontcare/common"
	"github.com/urfave/cli/v2"
)

type CapitalResponse []struct {
	Code string `json:"code"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetCapitals() (*CapitalResponse, error) {

	client := http.Client{}

	req, err := common.PrepareUnauthenticatedRequest("GET", common.CAPITALS_URL, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Error in getting capitals")
	}

	var capitals CapitalResponse

	json.NewDecoder(res.Body).Decode(&capitals)

	return &capitals, nil

}

func CapitalsToJSON(capitals *CapitalResponse) ([]byte, error) {

	return json.Marshal(capitals)

}

func CapitalsToASCIITable(capitals *CapitalResponse) *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"ID", "Code", "Name"})

	table.SetColWidth(100)

	for _, capital := range *capitals {
		table.Append([]string{
			fmt.Sprintf("%d", capital.ID),
			capital.Code,
			capital.Name,
		})
	}

	return table

}

var CapitalsCommand = &cli.Command{
	Name:    "capitals",
	Usage:   "Get all the capitals",
	Aliases: []string{"c"},
	Action: func(c *cli.Context) error {

		capitals, err := GetCapitals()
		if err != nil {
			log.Fatal(err)
			return err
		}

		output := c.String("output")

		switch output {
		case "json":
			json, err := CapitalsToJSON(capitals)
			if err != nil {
				log.Fatal(err)
				return err
			}
			fmt.Println(string(json))
		default:
			table := CapitalsToASCIITable(capitals)
			table.Render()
		}

		return nil

	},
}
