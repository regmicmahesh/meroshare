package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type UserDetails struct {
	Address                string    `json:"address"`
	Boid                   string    `json:"boid"`
	ClientCode             string    `json:"clientCode"`
	Contact                string    `json:"contact"`
	CreatedApproveDate     time.Time `json:"createdApproveDate"`
	CreatedApproveDateStr  string    `json:"createdApproveDateStr"`
	CustomerTypeCode       string    `json:"customerTypeCode"`
	Demat                  string    `json:"demat"`
	DematExpiryDate        string    `json:"dematExpiryDate"`
	Email                  string    `json:"email"`
	ExpiredDate            time.Time `json:"expiredDate"`
	ExpiredDateStr         string    `json:"expiredDateStr"`
	Gender                 string    `json:"gender"`
	ID                     int       `json:"id"`
	ImagePath              string    `json:"imagePath"`
	MeroShareEmail         string    `json:"meroShareEmail"`
	Name                   string    `json:"name"`
	PasswordChangeDate     time.Time `json:"passwordChangeDate"`
	PasswordChangedDateStr string    `json:"passwordChangedDateStr"`
	PasswordExpiryDate     time.Time `json:"passwordExpiryDate"`
	PasswordExpiryDateStr  string    `json:"passwordExpiryDateStr"`
	ProfileName            string    `json:"profileName"`
	RenderDashboard        bool      `json:"renderDashboard"`
	RenewedDate            time.Time `json:"renewedDate"`
	RenewedDateStr         string    `json:"renewedDateStr"`
	Username               string    `json:"username"`
}

var userDetails = &UserDetails{}

var AuthToken string

func login(username, password string, clientID int) (string, error) {

	id := strconv.Itoa(clientID)

	res, err := http.Post(LOGIN_URL, "application/json", strings.NewReader(`{"username": "`+username+`", "password": "`+password+`", "clientId": `+id+`}`))
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("Failed to login.")
	}
	return res.Header.Get("Authorization"), nil
}

func loadToken() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fullPath := path.Join(homedir, "credentials.json")
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}

	fmt.Println("✔️ Successfully loaded credentials.")

	creds := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		ClientId int    `json:"clientId"`
	}{}

	json.NewDecoder(file).Decode(&creds)

	AuthToken, err = login(creds.Username, creds.Password, creds.ClientId)

	fmt.Println("✔️ Successfully Logged In.")

	if err != nil {
		return err
	}

	fmt.Println("✔️ Successfully Loaded Token.")

	return nil

}

func getOwnDetails(render bool) (*UserDetails, error) {

	err := loadToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", OWN_DETAILS_URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", AuthToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Failed to get user details.")
	}
	json.NewDecoder(res.Body).Decode(userDetails)
	fmt.Println("✔️ Successfully loaded user details.")

	if render {
		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"Field", "Value"})
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

		table.Render()
	}

	return userDetails, nil
}

func printPortfolio() error {

	_, err := getOwnDetails(false)
	if err != nil {
		return err
	}

	client := &http.Client{}


	data := &struct {
		SortBy     string   `json:"sortBy"`
		Demat      []string `json:"demat"`
		ClientCode string   `json:"clientCode"`
		Page       int      `json:"page"`
		Size       int      `json:"size"`
		SortAsc    bool     `json:"sortAsc"`
	}{
		SortBy:     "script",
		Demat:      []string{userDetails.Demat},
		ClientCode: userDetails.ClientCode,
		Page:       1,
		Size:       1000,
		SortAsc:    true,
	}

	jsonBuffer := new(bytes.Buffer)
	err = json.NewEncoder(jsonBuffer).Encode(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", PORTFOLIO_URL, jsonBuffer)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", AuthToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("✔️ Successfully loaded portfolio.")


	resStruct := &struct {
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
	}{}

	err = json.NewDecoder(res.Body).Decode(resStruct)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Script", "Current Balance", "Previous Closing Price", "Value of Previous Closing Price", "Last Transaction Price", "Value as of LTP"})

	for _, script := range resStruct.MeroShareMyPortfolio {
		currBalance := strconv.FormatFloat(script.CurrentBalance, 'f', 2, 64)
		valueOfLTP := strconv.FormatFloat(script.ValueOfLastTransPrice, 'f', 2, 64)
		table.Append([]string{script.Script, currBalance, script.PreviousClosingPrice, script.ValueAsOfPreviousClosingPrice, script.LastTransactionPrice, valueOfLTP})
	}

	table.Append([]string{"Total", "", "", resStruct.TotalValueAsOfPreviousClosingPrice, "", resStruct.TotalValueAsOfLastTransactionPrice})

	table.Render()

	return nil

}
