package usecase

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"reflect"

	"github.com/labstack/echo"
)

type activationsUseCase struct {
	aRepo activations.Repository
	rRepo registrations.Repository
}

// ActivationsUseCase represent Activations Use Case
func ActivationsUseCase(
	aRepo activations.Repository, rRepo registrations.Repository) activations.UseCase {
	return &activationsUseCase{
		aRepo: aRepo,
		rRepo: rRepo,
	}
}

func (act *activationsUseCase) PostActivations(c echo.Context, pa models.PayloadActivations) error {
	acc, err := act.checkApplication(c, pa)

	if err != nil {
		return err
	}

	if acc.Status == models.ActivationsStatus {
		return models.ErrAlreadyActivated
	}

	err = acc.MappingCardActivationsData(c, pa)

	if err != nil {
		return models.ErrMappingData
	}

	// Activations to BRI
	act.activationsToBRI(c, acc, pa)

	// Activations to core
	act.activationsToCore(c, acc)

	errUpdateAct := act.aRepo.PostActivations(c, acc)

	if errUpdateAct != nil {
		return models.ErrPostActivationsFailed
	}

	return nil
}

func (act *activationsUseCase) checkApplication(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	appNumber := r.FieldByName("ApplicationNumber")

	if appNumber.IsZero() {
		return models.Account{}, nil
	}

	acc := models.Account{Application: models.Applications{ApplicationNumber: appNumber.String()}}
	err := act.rRepo.GetAccountByAppNumber(c, &acc)

	if err != nil {
		return models.Account{}, models.ErrAppNumberNotFound
	}

	if acc.BrixKey == "" {
		return models.Account{}, models.ErrEmptyBrixkey
	}

	return acc, nil
}

func (act *activationsUseCase) activationsToCore(c echo.Context, acc models.Account) error {
	respSwitching := api.SwitchingResponse{}
	requestDataSwitching := map[string]interface{}{
		"cif":        acc.CIF,
		"noRek":      acc.Application.SavingAccount,
		"branchCode": acc.BranchCode,
	}
	req := api.MappingRequestSwitching(requestDataSwitching)
	errSwitching := api.RetryableSwitchingPost(c, req, "/goldcard/aktivasi", &respSwitching)

	if errSwitching != nil {
		return errSwitching
	}

	if respSwitching.ResponseCode != api.ApiRCSuccess {
		logger.Make(c, nil).Debug(models.DynamicErr(models.ErrSwitchingAPIRequest, []interface{}{respSwitching.ResponseCode, respSwitching.ResponseDesc}))
		return errSwitching
	}

	return nil
}

func (act *activationsUseCase) activationsToBRI(c echo.Context, acc models.Account, pa models.PayloadActivations) error {
	respBRI := api.BriResponse{}
	requestDataBRI := map[string]interface{}{
		"briXkey":       acc.BrixKey,
		"expDate":       pa.ExpDate,
		"lastSixDigits": pa.LastSixDigits,
	}
	reqBRIBody := api.BriRequest{RequestData: requestDataBRI}
	errBRI := api.RetryableBriPost(c, "/v1/cobranding/card/activation", reqBRIBody, &respBRI)

	if errBRI != nil {
		return errBRI
	}

	return nil
}
