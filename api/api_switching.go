package api

import (
	"encoding/json"
	_apiRequestsUseCase "gade/srv-goldcard/apirequests/usecase"
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

// SwitchingResponse struct represents a response for API Switching
type SwitchingResponse struct {
	ResponseCode string            `json:"responseCode"`
	ResponseDesc string            `json:"responseDesc"`
	Data         string            `json:"data"`
	ResponseData map[string]string `json:"responseData,omitempty"`
}

// APIswitching struct represents a request for API Switching
type APIswitching struct {
	Host        *url.URL
	API         API
	Method      string
	AccessToken string
	ctx         echo.Context
}

// MappingRequestSwitching mapping request switching
func MappingRequestSwitching(req map[string]interface{}) interface{} {
	req["channelId"] = os.Getenv(`SWITCHING_CHANNEL_ID`)
	req["clientId"] = os.Getenv(`SWITCHING_CLIENT_ID`)
	req["flag"] = os.Getenv(`SWITCHING_FLAG`)

	return req
}

// NewSwitchingAPI is function to initiate a Switching API request
func NewSwitchingAPI(c echo.Context) (APIswitching, error) {
	apiSwitching := APIswitching{}
	apiSwitching.ctx = c
	url, err := url.Parse(os.Getenv(`SWITCHING_HOST`))

	if err != nil {
		return apiSwitching, err
	}

	api, err := NewAPI(apiSwitching.ctx, os.Getenv(`SWITCHING_HOST`), echo.MIMEApplicationJSON)

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
		return nil, err
	}

	err = switc.mappingDataResponseSwitching(v)

	if err != nil {
		logger.Make(switc.ctx, nil).Debug(err)

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
	params := map[string]string{"grant_type": "password", "username": os.Getenv(`SWITCHING_CLIENT_ID`), "password": os.Getenv(`SWITCHING_PASSWORD_TOKEN`)}
	endpoint := "/oauth/token"
	api, err := NewAPI(switc.ctx, os.Getenv(`SWITCHING_HOST`), echo.MIMEApplicationForm)

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

// SwitchingPost represent Post Switching API Request
func SwitchingPost(c echo.Context, body interface{}, path string, response interface{}) error {
	switching, err := NewSwitchingAPI(c)

	if err != nil {
		return err
	}

	req, err := switching.Request(path, echo.POST, body)

	if err != nil {
		return err
	}

	r, err := switching.Do(req, response)

	go func() {
		_ = _apiRequestsUseCase.ARUseCase.PostAPIRequest(c, r.StatusCode, switching.API, body, response)
	}()

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	return nil

}

// RetryableSwitchingPost function to retryable request Switching API with post method
func RetryableSwitchingPost(c echo.Context, body interface{}, path string, response interface{}) error {
	fn := func() error {
		return SwitchingPost(c, body, path, response)
	}

	err := RetryablePost(c, "SWITCHING API: POST "+path, fn)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (switc *APIswitching) mappingDataResponseSwitching(v interface{}) error {
	resp := v.(*SwitchingResponse)

	if resp.Data == "" {
		return nil
	}

	err := json.Unmarshal([]byte(resp.Data), &resp.ResponseData)
	resp.Data = ""

	if err != nil {
		logger.Make(switc.ctx, nil).Fatal("Response Data Error Unmarshal")
	}

	return nil
}
