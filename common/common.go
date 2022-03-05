package common

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var AUTH_TOKEN = ""

const BASE_URL = "https://webbackend.cdsc.com.np/api/"
const GET_CAPITALS_URL = BASE_URL + "meroShare/capital/"
const LOGIN_URL = BASE_URL + "meroShare/auth/"
const OWN_DETAILS_URL = BASE_URL + "meroShare/ownDetail/"
const CAPITALS_URL = BASE_URL + "meroShare/capital/"
const PORTFOLIO_URL = BASE_URL + "meroShareView/myPortfolio/"
const APPLICATION_REPORT_URL = BASE_URL + "meroShare/applicantForm/active/search/"
const OLD_APPLICATION_REPORT_URL = BASE_URL + "meroShare/migrated/applicantForm/search/"
const APPLICATION_REPORT_DETAILS_URL = BASE_URL + "meroShare/applicantForm/report/detail/"

var MSHAREKEY string

func GetKey() {
	if os.Getenv("MSHAREKEY") != "" {
		MSHAREKEY = os.Getenv("MSHAREKEY")
		if len(MSHAREKEY) != 25 {
			log.Fatal("Must be of 25 length")
			os.Exit(1)
		}
	} else {
		MSHAREKEY = "abc&1*~#3^^#s00=(^^7%b35"
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientId int    `json:"clientId"`
}

func init() {
	GetKey()
}

func Login(username, password string, clientID int) (string, error) {

	id := strconv.Itoa(clientID)

	res, err := http.Post(LOGIN_URL, "application/json", strings.NewReader(`{"username": "`+username+`", "password": "`+password+`", "clientId": `+id+`}`))
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("Failed to login.")
	}

	AUTH_TOKEN = res.Header.Get("Authorization")

	return AUTH_TOKEN, nil
}

func Rehydrate() (*Credentials, error) {
	creds := &Credentials{}
	homedir, _ := os.UserHomeDir()
	file, err := os.Open(path.Join(homedir, "credentials.json"))
	if err != nil {
		return nil, errors.New("Failed to open credentials file")
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(creds)
	creds.Password, err = Decrypt(creds.Password, MSHAREKEY)
	if err != nil {
		return nil, errors.New("Failed to encrypt credentials")
	}
	if err != nil {
		return nil, errors.New("Failed to write credentials")
	}
	return creds, nil

}

func Hydrate(username, password string, clientID int) (*Credentials, error) {
	creds := &Credentials{
		Username: username,
		Password: password,
		ClientId: clientID,
	}
	homedir, _ := os.UserHomeDir()
	file, err := os.Create(path.Join(homedir, "credentials.json"))
	if err != nil {
		return nil, errors.New("Failed to create credentials file")
	}
	creds.Password, err = Encrypt(creds.Password, MSHAREKEY)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Failed to encrypt credentials")
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(creds)
	if err != nil {
		return nil, errors.New("Failed to write credentials")
	}

	return creds, nil
}

func PrepareAuthenticatedRequest(method, url string, body io.Reader) (*http.Request, error) {

	if AUTH_TOKEN == "" {
		creds, err := Rehydrate()

		if err != nil {
			return nil, err
		}

		AUTH_TOKEN, err = Login(creds.Username, creds.Password, creds.ClientId)

		if err != nil {
			return nil, err
		}

	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", AUTH_TOKEN)
	req.Header.Add("Content-Type", "application/json")

	return req, nil

}

func PrepareUnauthenticatedRequest(method, url string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil

}
