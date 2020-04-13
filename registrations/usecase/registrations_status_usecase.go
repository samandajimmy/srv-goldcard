package usecase

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"reflect"

	"github.com/labstack/echo"
)

func (reg *registrationsUseCase) CheckApplication(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	appNumber := r.FieldByName("ApplicationNumber")

	if appNumber.IsZero() {
		return models.Account{}, nil
	}

	acc := models.Account{Application: models.Applications{ApplicationNumber: appNumber.String()}}
	err := reg.regRepo.GetAccountByAppNumber(c, &acc)

	if err != nil {
		return models.Account{}, models.ErrAppNumberNotFound
	}

	return acc, nil
}

func (reg *registrationsUseCase) GetAppStatus(c echo.Context, pl models.PayloadAppNumber) (models.AppStatus, error) {
	var appStatus models.AppStatus
	// Get account by app number
	acc, err := reg.CheckApplication(c, pl)

	if err != nil {
		return appStatus, err
	}

	// concurrently get app status from BRI API then update to our DB
	go func() {
		resp := api.BriResponse{}
		reqBody := map[string]interface{}{
			"briXkey": acc.BrixKey,
		}

		err := api.RetryableBriPost(c, "/v1/cobranding/card/appstatus", reqBody, &resp)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			return
		}

		// update application status
		data := resp.Data[0]

		if _, ok := data["appStatus"].(string); !ok {
			logger.Make(c, nil).Debug(err)
			return
		}

		acc.Application.ID = acc.ApplicationID
		acc.Application.SetStatus(data["appStatus"].(string))
		err = reg.regRepo.UpdateAppStatus(c, acc.Application)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			return
		}
	}()

	appStatus, err = reg.regRepo.GetAppStatus(c, acc.Application)

	if err != nil {
		return appStatus, models.ErrGetAppStatus
	}

	return appStatus, nil
}

func (reg *registrationsUseCase) updateSTLPrice(c echo.Context, acc models.Account) {
	hargeEmas, err := reg.rrr.GetCurrentGoldSTL(c)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}

	acc.Card.StlLimit = hargeEmas
	_, err = reg.regRepo.UpdateCardLimit(c, acc, false)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return
	}
}

func (reg *registrationsUseCase) afterOpenGoldcard(c echo.Context, acc *models.Account,
	briPl models.PayloadBriRegister, accChan chan models.Account, errAppBri, errAppCore chan error) error {
	var notif models.PdsNotification
	accChannel := <-accChan
	// function to apply to bri
	applyBri := func() {
		err := reg.briApply(c, acc, briPl)
		if err != nil {
			logger.Make(c, nil).Debug(err)
			errAppBri <- err
			return
		}
		errAppBri <- nil
	}

	for {
		select {
		case err := <-errAppCore:
			if err == nil {
				go applyBri()
			}
			if err != nil {
				// send notif app failed
				notif.GcApplication(accChannel, "failed")
				_ = reg.rrr.SendNotification(c, notif, "")
				return err
			}
		case err := <-errAppBri:
			if err != nil {
				// send notif app failed
				notif.GcApplication(accChannel, "failed")
				_ = reg.rrr.SendNotification(c, notif, "")
				return err
			}
			if err == nil {
				// send notif app succeeded
				notif.GcApplication(accChannel, "succeeded")
				_ = reg.rrr.SendNotification(c, notif, "")
				return err
			}
		}
	}
}
