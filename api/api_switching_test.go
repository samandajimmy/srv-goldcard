package api_test

import (
	"fmt"
	"gade/srv-goldcard/api"
	"testing"

	"github.com/labstack/echo"
)

var stlRequest = map[string]interface{}{
	"channelId": "6017",
	"clientId":  "9997",
	"flag":      "k",
}

func TestNewSwitchingAPI(t *testing.T) {
	response := map[string]interface{}{}
	body := stlRequest

	switc, _ := api.NewSwitchingAPI(nil)
	req, _ := switc.Request("/param/stl", echo.POST, body)
	resp, _ := switc.Do(req, &response)

	fmt.Println(req)
	fmt.Println(resp)
	fmt.Println(response)
}
