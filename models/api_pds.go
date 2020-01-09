package models

import (
	"net/http"
	"net/url"
	"os"
)

// APIpds struct represents a request for API PDS
type APIpds struct {
	Host *url.URL
	API  API
}

// NewPdsAPI is function to initiate a PDS API request
func NewPdsAPI(contentType string) (APIpds, error) {
	apiPds := APIpds{}
	url, err := url.Parse(os.Getenv(`PDS_API_HOST`))

	if err != nil {
		return apiPds, err
	}

	api, err := NewAPI(os.Getenv(`PDS_API_HOST`), contentType)

	if err != nil {
		return apiPds, err
	}

	apiPds.Host = url
	apiPds.API = api

	return apiPds, nil
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
