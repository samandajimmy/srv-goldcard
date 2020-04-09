package api_test

import (
	"fmt"
	"gade/srv-goldcard/api"
	"testing"

	"github.com/labstack/echo"
)

func TestNewPdsAPI(t *testing.T) {
	response := map[string]interface{}{}
	body := map[string]interface{}{
		"message": "ini message nya",
		"noHp":    "081511150290",
	}

	pds, _ := api.NewPdsAPI(nil, echo.MIMEApplicationForm)
	req, _ := pds.Request("/notification/send_sms_promo", echo.POST, body)
	resp, _ := pds.Do(req, &response)

	fmt.Println(req)
	fmt.Println(resp)
	fmt.Println(response)
}
