package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/regmicmahesh/merosharemorelikeidontcare/common"
	"github.com/urfave/cli/v2"
)

type UserDetails struct {
	Address        string `json:"address"`
	Boid           string `json:"boid"`
	ClientCode     string `json:"clientCode"`
	Contact        string `json:"contact"`
	Demat          string `json:"demat"`
	Email          string `json:"email"`
	Gender         string `json:"gender"`
	ID             int    `json:"id"`
	MeroShareEmail string `json:"meroShareEmail"`
	Name           string `json:"name"`
	ProfileName    string `json:"profileName"`
	Username       string `json:"username"`
}

func GetUserDetails() (*UserDetails, error) {

	client := http.Client{}

	req, err := common.PrepareAuthenticatedRequest("GET", common.OWN_DETAILS_URL, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	userDetails := &UserDetails{}

	if res.StatusCode != 200 {
		return nil, errors.New("Error getting user details")
	}

	json.NewDecoder(res.Body).Decode(userDetails)

	return userDetails, nil

}

func DetailsToJSON(userDetails *UserDetails) ([]byte, error) {
	return json.Marshal(userDetails)
}

func DetailsToASCIITable(userDetails *UserDetails) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"ID", strconv.Itoa(userDetails.ID)})
	table.Append([]string{"Username", userDetails.Username})
	table.Append([]string{"Name", userDetails.Name})
	table.Append([]string{"Email", userDetails.Email})
	table.Append([]string{"Address", userDetails.Address})
	table.Append([]string{"Gender", userDetails.Gender})
	table.Append([]string{"Boid", userDetails.Boid})
	table.Append([]string{"Client Code", userDetails.ClientCode})
	table.Append([]string{"Contact", userDetails.Contact})
	table.Append([]string{"Demat", userDetails.Demat})

	return table

}

var DetailsCommand = &cli.Command{
	Name:  "details",
	Usage: "Get user details",
	Action: func(c *cli.Context) error {

		userDetails, err := GetUserDetails()
		if err != nil {
			fmt.Println(err)
			return err
		}

		table := DetailsToASCIITable(userDetails)
		table.Render()

		return nil
	},
}

