package api_test

import (
	"fmt"
	"srv-goldcard/internal/pkg/api"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestNewPdsAPI(t *testing.T) {
	response := map[string]interface{}{}
	body := map[string]string{
		"email":    "082141217929",
		"password": "gadai123",
		"agen":     "android",
		"version":  "3",
	}

	pds, _ := api.NewPdsAPI(nil, echo.MIMEApplicationForm)
	req, _ := pds.Request("/auth/login/new", echo.POST, body)
	resp, _ := pds.Do(req, &response)

	fmt.Println(req)
	fmt.Println(resp)
	fmt.Println(response)

	assert.Equal(t, "SAMANDA RASU", response)
}
