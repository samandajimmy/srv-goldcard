package usecase

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"reflect"
	"time"

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

		err := api.RetryableBriPost(c, "/card/appstatus", reqBody, &resp)

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

func (reg *registrationsUseCase) RefreshAppTimeoutJob() {
	// update app that should be timeout
	err := reg.regRepo.ForceUpdateAppStatusTimeout()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	// get apps that need to be timeout later
	apps, err := reg.regRepo.GetAppOngoing()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	for _, app := range apps {
		go reg.appTimeoutJob(nil, app, models.NowUTC())
	}

}

func (reg *registrationsUseCase) afterOpenGoldcard(c echo.Context, acc *models.Account,
	briPl models.PayloadBriRegister, accChan chan models.Account, errAppBri chan error) error {
	var notif models.PdsNotification
	accChannel := <-accChan
	// function to apply to bri
	go func() {
		err := reg.briApply(c, acc, briPl)
		if err != nil {
			logger.Make(c, nil).Debug(err)
			errAppBri <- err
			return
		}
		errAppBri <- nil
	}()

	err := <-errAppBri

	if err != nil {
		// send notif app failed
		notif.GcApplication(accChannel, "failed")
		_ = reg.rrr.SendNotification(c, notif, "")
		return err
	}

	// send notif app succeeded
	notif.GcApplication(accChannel, "succeeded")
	_ = reg.rrr.SendNotification(c, notif, "")
	return nil
}

func (reg *registrationsUseCase) appTimeoutJob(c echo.Context, app models.Applications, now time.Time) {
	diff := app.ExpiredAt.Sub(now)
	delay := time.Duration(diff.Seconds())

	go func() {
		logger.Make(c, nil).Debug("Store timeout job to background for appNumber: " + app.ApplicationNumber)
		time.Sleep(delay * time.Second)
		logger.Make(c, nil).Debug("Start to make appNumber: " + app.ApplicationNumber + " expired!")
		_ = reg.regRepo.UpdateAppStatusTimeout(c, app)
	}()
}
