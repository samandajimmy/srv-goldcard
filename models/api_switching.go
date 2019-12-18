package models

import (
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo"
)

var (
	switchingHost          = os.Getenv(`SWITCHING_HOST`)
	switchingClientID      = os.Getenv(`SWITCHING_CLIENT_ID`)
	switchingChannelID     = os.Getenv(`SWITCHING_CHANNEL_ID`)
	switchingUserName      = os.Getenv(`SWITCHING_USERNAME`)
	switchingPassword      = os.Getenv(`SWITCHING_PASSWORD`)
	switchingPasswordToken = os.Getenv(`SWITCHING_PASSWORD_TOKEN`)
)

// APIswitching struct represents a request for API Switching
type APIswitching struct {
	Host        *url.URL
	API         API
	Method      string
	AccessToken string
}

// NewSwitchingAPI is function to initiate a Switching API request
func NewSwitchingAPI() (APIswitching, error) {
	apiSwitching := APIswitching{}
	url, err := url.Parse(switchingHost)

	if err != nil {
		return apiSwitching, err
	}

	api, err := NewAPI(switchingHost, echo.MIMEApplicationJSON)

	if err != nil {
		return apiSwitching, err
	}

	err = apiSwitching.setAccessTokenSwitching()

	if err != nil {
		return apiSwitching, err
	}

	apiSwitching.Host = url
	apiSwitching.API = api

	return apiSwitching, nil
}

// Do is a function to execute the http request
func (switc *APIswitching) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := switc.API.Do(req, v)

	if err != nil {
		return resp, err
	}

	return resp, err
}

// Request represent global API Request
func (switc *APIswitching) Request(endpoint, method string, body interface{}) (*http.Request, error) {
	switc.Method = method
	req, err := switc.API.Request(endpoint, method, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+switc.AccessToken)

	return req, nil
}

func (switc *APIswitching) setAccessTokenSwitching() error {
	response := map[string]interface{}{}
	params := map[string]string{"grant_type": "password", "username": switchingClientID, "password": switchingPasswordToken}
	endpoint := "/oauth/token"
	api, err := NewAPI(switchingHost, echo.MIMEApplicationForm)

	if err != nil {
		return err
	}

	req, err := api.Request(endpoint, echo.POST, params)

	if err != nil {
		return err
	}

	req.SetBasicAuth(switchingUserName, switchingPassword)

	_, err = api.Do(req, &response)

	if err != nil {
		return err
	}

	switc.AccessToken = response["access_token"].(string)

	return nil
}
