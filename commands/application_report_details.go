package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/regmicmahesh/merosharemorelikeidontcare/common"
	"github.com/urfave/cli/v2"
)

type ApplicationReportDetailsResponse struct {
	AccountNumber        string    `json:"accountNumber"`
	Action               string    `json:"action"`
	Amount               float64   `json:"amount"`
	ApplicantFormID      int       `json:"applicantFormId"`
	AppliedDate          time.Time `json:"appliedDate"`
	AppliedKitta         int       `json:"appliedKitta"`
	ClientName           string    `json:"clientName"`
	MaxIssueCloseDate    time.Time `json:"maxIssueCloseDate"`
	MeroShareID          int       `json:"meroShareId"`
	MeroshareRemark      string    `json:"meroshareRemark"`
	ReasonOrRemark       string    `json:"reasonOrRemark"`
	RegisteredBranchName string    `json:"registeredBranchName"`
	Remarks              string    `json:"remarks"`
	StageName            string    `json:"stageName"`
	StatusDescription    string    `json:"statusDescription"`
	StatusName           string    `json:"statusName"`
	SuspectStatusName    string    `json:"suspectStatusName"`
}

func GetApplicationDetails(id string) (*ApplicationReportDetailsResponse, error) {

	client := http.Client{}

	req, err := common.PrepareAuthenticatedRequest("GET", common.APPLICATION_REPORT_DETAILS_URL+id, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {

		req.URL.Path = "/api/meroShare/migrated/applicantForm/report/" + id
		res, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, errors.New("Invalid ID or Application not found")
		}

	}

	var app_reports ApplicationReportDetailsResponse

	json.NewDecoder(res.Body).Decode(&app_reports)

	return &app_reports, nil

}

func ApplicationReportDetailsToJSON(appreports *ApplicationReportDetailsResponse) ([]byte, error) {

	return json.Marshal(appreports)

}

func ApplicationReportDetailsToASCIITable(appreports *ApplicationReportDetailsResponse) *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Field", "Value"})

	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	table.Append([]string{"Account Number", appreports.AccountNumber})
	table.Append([]string{"Action", appreports.Action})
	table.Append([]string{"Amount", fmt.Sprintf("%f", appreports.Amount)})
	table.Append([]string{"Applicant Form ID", fmt.Sprintf("%d", appreports.ApplicantFormID)})
	table.Append([]string{"Applied Date", appreports.AppliedDate.Format("2006-01-02")})
	table.Append([]string{"Applied Kitta", fmt.Sprintf("%d", appreports.AppliedKitta)})
	table.Append([]string{"Client Name", appreports.ClientName})
	table.Append([]string{"Max Issue Close Date", appreports.MaxIssueCloseDate.Format("2006-01-02")})
	table.Append([]string{"MeroShare ID", fmt.Sprintf("%d", appreports.MeroShareID)})
	table.Append([]string{"Meroshare Remark", appreports.MeroshareRemark})
	table.Append([]string{"Reason Or Remark", appreports.ReasonOrRemark})
	table.Append([]string{"Registered Branch Name", appreports.RegisteredBranchName})
	table.Append([]string{"Remarks", appreports.Remarks})
	table.Append([]string{"Stage Name", appreports.StageName})
	table.Append([]string{"Status Description", appreports.StatusDescription})
	table.Append([]string{"Status Name", appreports.StatusName})
	table.Append([]string{"Suspect Status Name", appreports.SuspectStatusName})

	return table

}

var ApplicationReportDetailsCommand = &cli.Command{
	Name:    "app_reports_details",
	Aliases: []string{"ard"},
	Usage:   "ApplicationReports",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "Application Form ID to get details.",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {

		appreports, err := GetApplicationDetails(c.String("id"))
		if err != nil {
			log.Fatal(err)
		}

		switch c.String("output") {
		case "json":
			json, err := ApplicationReportDetailsToJSON(appreports)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
		default:
			table := ApplicationReportDetailsToASCIITable(appreports)
			table.Render()
		}

		return nil
	},
}
