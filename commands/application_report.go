package commands

import (
	"bytes"
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

type ApplicationReportResponse struct {
	Object []struct {
		CompanyShareID  int    `json:"companyShareId"`
		SubGroup        string `json:"subGroup"`
		Scrip           string `json:"scrip"`
		CompanyName     string `json:"companyName"`
		ShareTypeName   string `json:"shareTypeName"`
		ShareGroupName  string `json:"shareGroupName"`
		StatusName      string `json:"statusName"`
		ApplicantFormID int    `json:"applicantFormId"`
		AllotmentStatus string `json:"allotmentStatus"`
	} `json:"object"`
	TotalCount int `json:"totalCount"`
}

type ApplicationReportRequest struct {
	FilterFieldParams []struct {
		Key   string `json:"key"`
		Alias string `json:"alias"`
	} `json:"filterFieldParams"`
	Page                    int    `json:"page"`
	Size                    int    `json:"size"`
	SearchRoleViewConstants string `json:"searchRoleViewConstants"`
	FilterDateParams        []struct {
		Key       string `json:"key"`
		Condition string `json:"condition"`
		Alias     string `json:"alias"`
		Value     string `json:"value"`
	} `json:"filterDateParams"`
}

//{"filterFieldParams":[{"key":"companyShare.companyIssue.companyISIN.script","alias":"Scrip"},{"key":"companyShare.companyIssue.companyISIN.company.name","alias":"Company Name"}],"page":1,"size":200,"searchRoleViewConstants":"VIEW_APPLICANT_FORM_COMPLETE","filterDateParams":[{"key":"appliedDate","condition":"","alias":"","value":""},{"key":"appliedDate","condition":"","alias":"","value":""}]}

func GetApplications(activeOnly, withAllotmentStatus bool) (*ApplicationReportResponse, error) {

	client := http.Client{}

	request := &ApplicationReportRequest{
		Page:                    1,
		Size:                    200,
		SearchRoleViewConstants: "VIEW_APPLICANT_FORM_COMPLETE",
		FilterFieldParams: []struct {
			Key   string `json:"key"`
			Alias string `json:"alias"`
		}{
			{
				Key:   "companyShare.companyIssue.companyISIN.script",
				Alias: "Scrip",
			},
			{
				Key:   "companyShare.companyIssue.companyISIN.company.name",
				Alias: "Company Name",
			},
		},
		FilterDateParams: []struct {
			Key       string `json:"key"`
			Condition string `json:"condition"`
			Alias     string `json:"alias"`
			Value     string `json:"value"`
		}{
			{
				Key:       "appliedDate",
				Condition: "",
				Alias:     "",
				Value:     "",
			},
			{
				Key:       "appliedDate",
				Condition: "",
				Alias:     "",
				Value:     "",
			},
		},
	}

	requestJSON, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := common.PrepareAuthenticatedRequest("POST", common.APPLICATION_REPORT_URL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Error fetching applications")
	}

	var app_reports ApplicationReportResponse

	json.NewDecoder(res.Body).Decode(&app_reports)
	if !activeOnly {
		req.URL.Path = "/api/meroShare/migrated/applicantForm/search/"

		res, err = client.Do(req)
		if err != nil {
			return nil, err
		}

		var temp_app_reports ApplicationReportResponse

		json.NewDecoder(res.Body).Decode(&temp_app_reports)

		app_reports.TotalCount += temp_app_reports.TotalCount
		app_reports.Object = append(app_reports.Object, temp_app_reports.Object...)

	}

	if withAllotmentStatus {
		for i, object := range app_reports.Object {
			data, err := GetApplicationDetails(fmt.Sprintf("%d", object.ApplicantFormID))
			if err != nil {
				return nil, err
			}

			app_reports.Object[i].AllotmentStatus = data.StatusName

		}
	}

	return &app_reports, nil

}

func ApplicationReportsToJSON(appreports *ApplicationReportResponse) ([]byte, error) {

	return json.Marshal(appreports)

}

func ApplicationReportsToASCIITable(appreports *ApplicationReportResponse) *tablewriter.Table {

	table := tablewriter.NewWriter(os.Stdout)

	if len(appreports.Object) == 0 {
		table.Append([]string{"No applications found"})
		return table
	}

	table.SetHeader([]string{"Company Share ID", "Sub Group", "Scrip", "Company Name", "Share Type Name", "Share Group Name", "Status Name", "Applicant Form ID", "AllotmentStatus"})

	for _, object := range appreports.Object {
		table.Append([]string{
			fmt.Sprintf("%d", object.CompanyShareID),
			object.SubGroup,
			object.Scrip,
			object.CompanyName,
			object.ShareTypeName,
			object.ShareGroupName,
			object.StatusName,
			fmt.Sprintf("%d", object.ApplicantFormID),
			object.AllotmentStatus,
		})

	}
	return table
}

var ApplicationReportCommand = &cli.Command{
	Name:    "app_reports",
	Usage:   "ApplicationReports",
	Aliases: []string{"ar"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "active-only",
			Aliases: []string{"ac"},
			Usage:   "View only active applications.",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "with-allotment-status",
			Aliases: []string{"wa"},
			Usage:   "View applications along with allotment status.",
			Value:   false,
		},
	},
	Action: func(c *cli.Context) error {

		appreports, err := GetApplications(c.Bool("active-only"), c.Bool("with-allotment-status"))
		if err != nil {
			log.Fatal(err)
		}

		switch c.String("output") {
		case "json":
			json, err := ApplicationReportsToJSON(appreports)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(json))
		default:
			table := ApplicationReportsToASCIITable(appreports)
			table.Render()
		}

		return nil
	},
}
