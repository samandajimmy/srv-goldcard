package api

import (
	_apiRequestsUseCase "gade/srv-goldcard/apirequests/usecase"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo"
)

// PdsResponse struct to store response from PDS API
type PdsResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Errors  string                   `json:"errors"`
	Data    []map[string]interface{} `json:"data,omitempty"`
	User    map[string]interface{}   `json:"user,omitempty"`
	Token   string                   `json:"token,omitempty"`
}

// APIpds struct represents a request for API PDS
type APIpds struct {
	Host *url.URL
	API  API
	ctx  echo.Context
}

// NewPdsAPI is function to initiate a PDS API request
func NewPdsAPI(c echo.Context, contentType string) (APIpds, error) {
	apiPds := APIpds{}
	apiPds.ctx = c

	url, err := url.Parse(os.Getenv(`PDS_API_HOST`))

	if err != nil {
		return apiPds, err
	}

	api, err := NewAPI(apiPds.ctx, os.Getenv(`PDS_API_HOST`), contentType)

	if err != nil {
		return apiPds, err
	}

	apiPds.Host = url
	apiPds.API = api

	return apiPds, nil
}

// PdsPost function to request PDS API with post method
func PdsPost(c echo.Context, endpoint string, reqBody, resp interface{}, contentType string) error {
	pds, err := NewPdsAPI(c, contentType)

	if err != nil {
		return err
	}

	req, err := pds.Request(endpoint, echo.POST, reqBody)

	if err != nil {
		return err
	}

	r, err := pds.Do(req, resp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		resp = PdsResponse{
			Status:  "error",
			Message: err.Error(),
		}
	}

	res := resp.(*PdsResponse)

	if res.Status != "success" {
		r.StatusCode = http.StatusBadRequest
	}

	go func() {
		_ = _apiRequestsUseCase.ARUseCase.PostAPIRequest(c, r.StatusCode, pds.API, reqBody, resp)
	}()

	if r.StatusCode != http.StatusOK {
		return models.DynamicErr(models.ErrPdsAPIRequest, []interface{}{res.Status,
			res.Message})
	}

	if res.Status != "success" {
		return models.DynamicErr(models.ErrPdsAPIRequest, []interface{}{res.Status,
			res.Message})
	}

	return nil
}

// RetryablePdsPost function to retryable request PDS API with post method
func RetryablePdsPost(c echo.Context, endpoint string, reqBody interface{}, resp interface{}, contentType string) error {
	fn := func() error {
		return PdsPost(c, endpoint, reqBody, resp, contentType)
	}

	err := RetryablePost(c, "PDS API: POST "+endpoint, fn)

	if err != nil {
		return err
	}

	return nil
}

// Request represent PDS API Request
func (pds *APIpds) Request(endpoint string, method string, body interface{}) (*http.Request, error) {
	req, err := pds.API.Request(endpoint, method, body)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv(`PDS_API_BASIC_USER`), os.Getenv(`PDS_API_BASIC_PASS`))

	return req, nil
}

// Do is a function to execute the http request
func (pds *APIpds) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := pds.API.Do(req, v)

	if err != nil {
		return resp, err
	}

	return resp, err
}
