package models

import (
	"encoding/json"
	"gade/srv-goldcard/logger"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo"
)

var (
	// SwitchingRCInquiryAllow RC 14
	SwitchingRCInquiryAllow = "14"
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
	url, err := url.Parse(os.Getenv(`SWITCHING_HOST`))

	if err != nil {
		return apiSwitching, err
	}

	api, err := NewAPI(os.Getenv(`SWITCHING_HOST`), echo.MIMEApplicationJSON)

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
func (switc *APIswitching) Request(endpoint, method string, body map[string]interface{}) (*http.Request, error) {
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
	params := map[string]string{"grant_type": "password", "username": os.Getenv(`SWITCHING_CLIENT_ID`), "password": os.Getenv(`SWITCHING_PASSWORD_TOKEN`)}
	endpoint := "/oauth/token"
	api, err := NewAPI(os.Getenv(`SWITCHING_HOST`), echo.MIMEApplicationForm)

	if err != nil {
		return err
	}

	req, err := api.Request(endpoint, echo.POST, params)

	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv(`SWITCHING_USERNAME`), os.Getenv(`SWITCHING_PASSWORD`))

	_, err = api.Do(req, &response)

	if err != nil {
		return err
	}

	switc.AccessToken = response["access_token"].(string)

	return nil
}

//SwitchingPost represent Post Switching API Request
func SwitchingPost(body map[string]interface{}, path string) (map[string]json.RawMessage, error) {
	var response map[string]json.RawMessage
	body["channelId"] = os.Getenv(`SWITCHING_CHANNEL_ID`)
	body["clientId"] = os.Getenv(`SWITCHING_CLIENT_ID`)
	body["flag"] = "K"

	switc, err := NewSwitchingAPI()

	if err != nil {
		return nil, err
	}

	req, err := switc.Request(path, echo.POST, body)

	if err != nil {
		return nil, err
	}

	_, err = switc.Do(req, &response)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return nil, err
	}

	return response, nil

}
