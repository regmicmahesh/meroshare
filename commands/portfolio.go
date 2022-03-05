package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/regmicmahesh/merosharemorelikeidontcare/common"
	"github.com/urfave/cli/v2"
)

type PortfolioData struct {
	MeroShareMyPortfolio []struct {
		CurrentBalance                float64 `json:"currentBalance"`
		LastTransactionPrice          string  `json:"lastTransactionPrice"`
		PreviousClosingPrice          string  `json:"previousClosingPrice"`
		Script                        string  `json:"script"`
		ScriptDesc                    string  `json:"scriptDesc"`
		ValueAsOfLastTransactionPrice string  `json:"valueAsOfLastTransactionPrice"`
		ValueAsOfPreviousClosingPrice string  `json:"valueAsOfPreviousClosingPrice"`
		ValueOfLastTransPrice         float64 `json:"valueOfLastTransPrice"`
		ValueOfPrevClosingPrice       float64 `json:"valueOfPrevClosingPrice"`
	} `json:"meroShareMyPortfolio"`
	TotalItems                         int     `json:"totalItems"`
	TotalValueAsOfLastTransactionPrice string  `json:"totalValueAsOfLastTransactionPrice"`
	TotalValueAsOfPreviousClosingPrice string  `json:"totalValueAsOfPreviousClosingPrice"`
	TotalValueOfLastTransPrice         float64 `json:"totalValueOfLastTransPrice"`
	TotalValueOfPrevClosingPrice       float64 `json:"totalValueOfPrevClosingPrice"`
}

func preparePortfolioPayload(ud *UserDetails) (*bytes.Buffer, error) {
	data := &struct {
		SortBy     string   `json:"sortBy"`
		Demat      []string `json:"demat"`
		ClientCode string   `json:"clientCode"`
		Page       int      `json:"page"`
		Size       int      `json:"size"`
		SortAsc    bool     `json:"sortAsc"`
	}{
		SortBy:     "script",
		Demat:      []string{ud.Demat},
		ClientCode: ud.ClientCode,
		Page:       1,
		Size:       1000,
		SortAsc:    true,
	}
	jsonBuffer := new(bytes.Buffer)
	err := json.NewEncoder(jsonBuffer).Encode(data)
	if err != nil {
		return nil, err
	}
	return jsonBuffer, nil
}

func GetPortfolio() (*PortfolioData, error) {
	client := &http.Client{}

	ud, err := GetUserDetails()
	if err != nil {
		return nil, err
	}

	jsonBuffer, err := preparePortfolioPayload(ud)

	if err != nil {
		return nil, err
	}

	req, err := common.PrepareAuthenticatedRequest("POST", common.PORTFOLIO_URL, jsonBuffer)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var portfolioData = &PortfolioData{}
	err = json.NewDecoder(res.Body).Decode(&portfolioData)
	if err != nil {
		return nil, err
	}
	return portfolioData, nil

}

func PortfolioToASCIITable(pd *PortfolioData) *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Script", "Current Balance", "Previous Closing Price", "Value of Previous Closing Price", "Last Transaction Price", "Value as of LTP"})

	for _, script := range pd.MeroShareMyPortfolio {
		currBalance := strconv.FormatFloat(script.CurrentBalance, 'f', 2, 64)
		valueOfLTP := strconv.FormatFloat(script.ValueOfLastTransPrice, 'f', 2, 64)
		table.Append([]string{script.Script, currBalance, script.PreviousClosingPrice, script.ValueAsOfPreviousClosingPrice, script.LastTransactionPrice, valueOfLTP})
	}

	table.Append([]string{"Total", "", "", pd.TotalValueAsOfPreviousClosingPrice, "", pd.TotalValueAsOfLastTransactionPrice})

	return table

}

func PortfolioToJSON(pd *PortfolioData) string {
	jsonBuffer := new(bytes.Buffer)
	err := json.NewEncoder(jsonBuffer).Encode(pd)
	if err != nil {
		return ""
	}
	return jsonBuffer.String()
}

var PortfolioCommand = &cli.Command{
	Name:  "portfolio",
	Usage: "Get Portfolio",
	Aliases: []string{"p"},
	Action: func(c *cli.Context) error {
		portfolioData, err := GetPortfolio()
		if err != nil {
			return err
		}

		output := c.String("output")

		switch output {
		case "json":
			fmt.Println(PortfolioToJSON(portfolioData))
		default:
			PortfolioToASCIITable(portfolioData).Render()
		}

		return nil

	},
}
