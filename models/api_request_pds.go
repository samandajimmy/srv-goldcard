package models

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

// Client struct represents a request in ApiRequest
type Client struct {
	URL         *url.URL
	UserAgent   string
	HTTPClient  *http.Client
	ContentType string
}

// NewClientRequest for create new client request
func NewClientRequest(baseURL string, contentType string) (Client, error) {
	url, err := url.Parse(baseURL)

	if err != nil {
		return Client{}, err
	}

	return Client{
		URL:         url,
		UserAgent:   "goldcard",
		HTTPClient:  &http.Client{},
		ContentType: contentType,
	}, nil
}

// ApiRequest global API Request
func (c *Client) ApiRequest(ctx echo.Context, pathName string, method string, body interface{}, strct interface{}) (*http.Response, error) {
	var buf io.ReadWriter
	c.URL.Path += pathName

	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.URL.String(), buf)

	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", c.ContentType)
	req.Header.Set("User-Agent", c.UserAgent)

	response, err := c.do(req, &strct)

	if err != nil {
		return nil, err
	}

	return response, nil

}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)

	return resp, err
}
