package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"srv-goldcard/internal/pkg/api"
	"testing"

	"github.com/labstack/echo"
)

var stlRequest = map[string]interface{}{
	"channelId": "6017",
	"clientId":  "9997",
	"flag":      "k",
}

func TestNewSwitchingAPI(t *testing.T) {
	response := api.SwitchingResponse{}
	body := stlRequest

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	switc, _ := api.NewSwitchingAPI(c)
	reqBody, _ := switc.Request("/param/stl", echo.POST, body)
	resp, _ := switc.Do(reqBody, &response)

	fmt.Println(reqBody)
	fmt.Println(resp)
	fmt.Println(response)
}
