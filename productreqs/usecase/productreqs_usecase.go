package usecase

import (
	"encoding/json"
	"gade/srv-goldcard/productreqs"
	"os"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

type productreqsUseCase struct {
	prodReqsUseCase productreqs.UseCase
}

// ProductReqsUseCase represent product requirements Use Case
func ProductReqsUseCase() productreqs.UseCase {
	return &productreqsUseCase{}
}

// ProductRequirements represent to get all product requirements
func (prodreqs *productreqsUseCase) ProductRequirements(c echo.Context) (map[string]interface{}, error) {
	val := []byte("")
	viper.AddConfigPath(os.Getenv(`CONFIG_DIR`)) // load all configs
	viper.SetConfigName("product_requirements")
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		return nil, err
	}

	readResponse := viper.Get("requirements")

	json.Unmarshal(val, &readResponse)

	myMap := readResponse.(map[string]interface{})

	return myMap, err

}
