package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo"
)

var (
	hostBRI      = os.Getenv(`BRI_HOST`)
	grantType    = os.Getenv(`BRI_GRANT_TYPE`)
	clientID     = os.Getenv(`BRI_CLIENT_ID`)
	clientSecret = os.Getenv(`BRI_CLIENT_SECRET`)
)

// APIbri struct represents a request for API BRI
type APIbri struct {
	Host         *url.URL
	API          API
	Method       string
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

	err = apiBri.setAccessToken()

	if err != nil {
		return apiBri, err
	}

	apiBri.Host = url
	apiBri.API = api

	return apiBri, nil
}

// Request represent BRI API Request
func (bri *APIbri) Request(endpoint string, method string, body interface{}) (*http.Request, error) {
	bri.Method = method
	req, err := bri.API.Request(endpoint, method, body)

	if err != nil {
		return nil, err
	}

	bri.BRITimestamp = time.Now().UTC().Format(DateTimeFormatZone)
	err = bri.setBriSignature(endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+bri.AccessToken)
	req.Header.Set("BRI-Signature", bri.BRISignature)
	req.Header.Set("BRI-Timestamp", bri.BRITimestamp)

	return req, nil
}

// Do is a function to execute the http request
func (bri *APIbri) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := bri.API.Do(req, v)

	if err != nil {
		return resp, err
	}

	return resp, err
}

func (bri *APIbri) setBriSignature(endpoint string, body interface{}) error {
	var jsonData []byte
	var bodyStr string
	jsonData, err := json.Marshal(body)

	if err != nil {
		return err
	}

	if string(jsonData) != "null" {
		bodyStr = string(jsonData)
	}

	queryParams := url.Values{}
	queryParams.Set("path", endpoint)
	queryParams.Set("verb", bri.Method)
	queryParams.Set("token", "Bearer "+bri.AccessToken)
	queryParams.Set("timestamp", bri.BRITimestamp)
	queryParams.Set("body", bodyStr)
	queryParamStr := bri.getUnorderedURLQuery(queryParams)

	key := []byte(clientSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(queryParamStr))
	sEnc := base64.StdEncoding.EncodeToString(h.Sum(nil))
	bri.BRISignature = sEnc

	return nil
}

func (bri *APIbri) setAccessToken() error {
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

func (bri APIbri) getUnorderedURLQuery(queryParams url.Values) string {
	return "path=" + queryParams.Get("path") + "&verb=" + queryParams.Get("verb") +
		"&token=" + queryParams.Get("token") + "&timestamp=" + queryParams.Get("timestamp") +
		"&body=" + queryParams.Get("body")
}
