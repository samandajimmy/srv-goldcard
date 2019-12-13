package models

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

const (
	hostBRI      = "https://sandbox.partner.api.bri.co.id"
	grantType    = "client_credentials"
	clientID     = "YtkJkC2vrGGxsD7KVjSLWAk38cPq4thm"
	clientSecret = "HG2X0VaAP0nRmc4X"
)

// APIbri struct represents a request for API BRI
type APIbri struct {
	Host         *url.URL
	API          API
	AccessToken  string
	BRITimestamp string
	BRISignature string
}

// NewBriAPI is function to initiate a BRI API request
func NewBriAPI() (APIbri, error) {
	apiBri := APIbri{}
	url, err := url.Parse(hostBRI)

	if err != nil {
		return apiBri, err
	}

	api, err := NewAPI(hostBRI, echo.MIMEApplicationJSON)

	if err != nil {
		return apiBri, err
	}

	err = apiBri.getAccessToken()

	if err != nil {
		return apiBri, err
	}

	apiBri.Host = url
	apiBri.API = api

	return apiBri, nil
}

// Request is for blah blah
func (bri *APIbri) Request(endpoint string, method string, body interface{}) (*http.Request, error) {
	req, err := bri.API.Request(endpoint, method, body)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (bri *APIbri) getAccessToken() error {
	response := map[string]interface{}{}
	params := map[string]string{"client_id": clientID, "client_secret": clientSecret}
	endpoint := "/oauth/client_credential/accesstoken?grant_type=client_credentials"
	api, err := NewAPI(hostBRI, echo.MIMEApplicationForm)

	if err != nil {
		return err
	}

	req, err := api.Request(endpoint, echo.POST, params)

	if err != nil {
		return err
	}

	_, err = api.Do(req, &response)

	if err != nil {
		return err
	}

	bri.AccessToken = response["access_token"].(string)

	return nil
}
