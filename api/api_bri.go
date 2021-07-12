package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	_apiRequestsUseCase "gade/srv-goldcard/apirequests/usecase"

	"github.com/labstack/echo"
)

// BriResponse struct to store response from BRI API
type BriResponse struct {
	ResponseCode    string                   `json:"responseCode"`
	ResponseMessage string                   `json:"responseMessage"`
	ResponseData    interface{}              `json:"responseData,omitempty"`
	Data            []map[string]interface{} `json:"data,omitempty"`
	DataOne         map[string]interface{}   `json:"dataOne,omitempty"`
	DataString      string                   `json:"dataString,omitempty"`
	Status          map[string]interface{}   `json:"status,omitempty"`
}

// SetRC to get bri api response code
func (br *BriResponse) SetRC() {
	if br.Status == nil {
		return
	}

	if _, ok := br.Status["code"].(string); !ok {
		logger.Make(nil, nil).Fatal(models.ErrSetVar)
	}

	br.ResponseCode = br.Status["code"].(string)

	if _, ok := br.Status["desc"].(string); !ok {
		logger.Make(nil, nil).Fatal(models.ErrSetVar)
	}

	br.ResponseMessage = br.Status["desc"].(string)
}

// BriRequest struct to store request payload BRI API needed
type BriRequest struct {
	RequestData interface{} `json:"requestData"`
}

// APIbri struct represents a request for API BRI
type APIbri struct {
	Host         *url.URL
	API          API
	Method       string
	Endpoint     string
	AccessToken  string
	BRITimestamp string
	BRISignature string
	ctx          echo.Context
}

// NewBriAPI is function to initiate a BRI API request
func NewBriAPI(c echo.Context) (APIbri, error) {
	apiBri := APIbri{}
	apiBri.ctx = c
	url, err := url.Parse(os.Getenv(`BRI_HOST`))

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return apiBri, err
	}

	api, err := NewAPI(apiBri.ctx, os.Getenv(`BRI_HOST`), echo.MIMEApplicationJSON)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return apiBri, err
	}

	err = apiBri.setAccessToken()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return apiBri, err
	}

	apiBri.Host = url
	apiBri.API = api

	return apiBri, nil
}

// BriPost function to request BRI API with post method
func BriPost(c echo.Context, endpoint string, reqBody, resp interface{}) error {
	oldEndpoint := endpoint
	endpoint = os.Getenv(`BRI_BASE_PATH`) + oldEndpoint
	bri, err := NewBriAPI(c)

	if err != nil {
		return err
	}

	req, err := bri.Request(endpoint, echo.POST, reqBody)

	if err != nil {
		return err
	}

	r, err := bri.Do(req, resp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	res := resp.(*BriResponse)
	res.SetRC()

	go func() {
		_ = _apiRequestsUseCase.ARUseCase.PostAPIRequest(c, r.StatusCode, bri.API, reqBody, resp)
	}()

	if r.StatusCode != http.StatusOK {
		return models.DynamicErr(models.ErrBriAPIRequest, []interface{}{res.ResponseCode,
			res.ResponseMessage})
	}

	if res.ResponseCode != "00" && !isWhitelisted(oldEndpoint) {
		return models.DynamicErr(models.ErrBriAPIRequest, []interface{}{res.ResponseCode,
			res.ResponseMessage})
	}

	return nil
}

// RetryableBriPost function to retryable request BRI API with post method
func RetryableBriPost(c echo.Context, endpoint string, reqBody interface{}, resp interface{}) error {
	fn := func() error {
		return BriPost(c, endpoint, reqBody, resp)
	}

	err := RetryablePost(c, "BRI API: POST "+endpoint, fn)

	if err != nil {
		return err
	}

	return nil
}

// Request represent BRI API Request
func (bri *APIbri) Request(endpoint string, method string, body interface{}) (*http.Request, error) {
	// show request log
	debugStart := fmt.Sprintf("Start to request BRI API: %s %s", method, endpoint)
	logger.MakeWithoutReportCaller(bri.ctx, nil).Info(debugStart)
	bri.Method = method
	bri.Endpoint = endpoint
	bearerToken := "Bearer " + bri.AccessToken
	req, err := bri.API.Request(endpoint, method, body)

	if err != nil {
		return nil, err
	}

	bri.BRITimestamp = time.Now().UTC().Format(models.DateTimeFormatZone)
	err = bri.setBriSignature(endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", bearerToken)
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

	// remapping response struct
	bri.remappingBriResponseData(v)

	// show response log
	debugEnd := fmt.Sprintf("End of request BRI API: %s %s", bri.Method, bri.Endpoint)
	logger.MakeWithoutReportCaller(bri.ctx, nil).Info(debugEnd)

	return resp, err
}

func (bri *APIbri) remappingBriResponseData(v interface{}) {
	rr := reflect.ValueOf(v)
	rrdi := rr.Elem().FieldByName("ResponseData")
	rrd := reflect.ValueOf(rrdi.Interface())

	if rrdi.IsZero() {
		return
	}

	if rrd.Kind() == reflect.Slice {
		rd := rr.Elem().FieldByName("Data")
		destType := reflect.TypeOf([]map[string]interface{}{})
		slc := reflect.MakeSlice(destType, rrd.Len(), rrd.Len())
		rd.Set(slc)

		for i := 0; i < slc.Len(); i++ {
			dt := rrd.Index(i).Interface().(map[string]interface{})
			rd.Index(i).Set(reflect.ValueOf(dt))
		}
		return
	}

	if rrd.Kind() == reflect.String {
		rr.Elem().FieldByName("DataString").Set(rrd)
		return
	}

	rr.Elem().FieldByName("DataOne").Set(rrd)
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

	key := []byte(os.Getenv(`BRI_CLIENT_SECRET`))
	h := hmac.New(sha256.New, key)
	_, err = h.Write([]byte(queryParamStr))

	if err != nil {
		logger.Make(bri.ctx, nil).Debug(err)
		return nil
	}

	sEnc := base64.StdEncoding.EncodeToString(h.Sum(nil))
	bri.BRISignature = sEnc

	return nil
}

func (bri *APIbri) setAccessToken() error {
	response := map[string]interface{}{}
	params := map[string]string{"client_id": os.Getenv(`BRI_CLIENT_ID`), "client_secret": os.Getenv(`BRI_CLIENT_SECRET`)}
	endpoint := "/oauth/client_credential/accesstoken?grant_type=client_credentials"
	api, err := NewAPI(bri.ctx, os.Getenv(`BRI_HOST`), echo.MIMEApplicationForm)

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

	if _, ok := response["access_token"].(string); !ok {
		return models.ErrSetVar
	}

	bri.AccessToken = response["access_token"].(string)

	return nil
}

func (bri APIbri) getUnorderedURLQuery(queryParams url.Values) string {
	return "path=" + queryParams.Get("path") + "&verb=" + queryParams.Get("verb") +
		"&token=" + queryParams.Get("token") + "&timestamp=" + queryParams.Get("timestamp") +
		"&body=" + queryParams.Get("body")
}
