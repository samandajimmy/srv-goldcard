package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
)

// Client struct represents a request in ApiRequest
type Client struct {
	URL         *url.URL
	UserAgent   string
	HTTPClient  *http.Client
	ContentType string
	Payload     interface{} `json:"payload,omitempty"`
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

// ApiRequest represent global API Request
func (c *Client) ApiRequest(ctx echo.Context, pathName string, method string, body interface{}, strct interface{}) (*http.Response, error) {
	var jsonData []byte
	var mapData map[string]interface{}

	c.URL.Path += pathName

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		switch ct := c.ContentType; ct {
		case "application/x-www-form-urlencoded":
			return c.requestUrlEncoded(method, jsonData, mapData, body, &strct)
		case "application/json":
			return c.requestJson(method, jsonData, body, &strct)
		}
	}

	return nil, err

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

func (c *Client) requestUrlEncoded(method string, jsonData []byte, mapData map[string]interface{}, body interface{}, strct interface{}) (*http.Response, error) {
	data := url.Values{}
	json.Unmarshal(jsonData, &mapData)
	for k, v := range mapData {
		data.Set(k, fmt.Sprintf("%v", v))
	}
	req, err := http.NewRequest(method, c.URL.String(), strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Accept", c.ContentType)
	req.Header.Set("User-Agent", c.UserAgent)
	response, err := c.do(req, &strct)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) requestJson(method string, jsonData []byte, body interface{}, strct interface{}) (*http.Response, error) {
	var buf io.ReadWriter
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

	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("Accept", c.ContentType)
	req.Header.Set("User-Agent", c.UserAgent)

	response, err := c.do(req, &strct)
	if err != nil {
		return nil, err
	}
	return response, nil
}
