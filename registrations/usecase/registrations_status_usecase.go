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

	appStatus, err = reg.regRepo.GetAppStatus(c, acc.Application)

	if err != nil {
		return appStatus, models.ErrGetAppStatus
	}

	if acc.Application.Status == models.AppStatusForceDeliver {
		return appStatus, nil
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
	accs, err := reg.regRepo.GetAppOngoing()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	var diff, delay time.Duration

	go func() {
		for _, acc := range accs {
			diff = acc.Application.ExpiredAt.Sub(models.NowUTC())
			delay = time.Duration(diff.Seconds())

			reg.appTimeoutJob(nil, acc, diff, delay)
		}
	}()

}

func (reg *registrationsUseCase) afterOpenGoldcard(c echo.Context, acc *models.Account,
	briPl models.PayloadBriRegister, accChan chan models.Account) error {
	var notif models.PdsNotification
	accChannel := <-accChan
	// function to apply to bri
	err := reg.briApply(c, acc, briPl)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		// send notif app failed
		_ = reg.appFailedNotification(c, accChannel)
		return err
	}

	// send notif app succeeded
	notif.GcApplication(accChannel, "succeeded")
	_ = reg.rrr.SendNotification(c, notif, "")
	return nil
}

func (reg *registrationsUseCase) appTimeoutJob(c echo.Context, acc models.Account, diff, delay time.Duration) {
	var notif models.PdsNotification

	go func() {
		logger.Make(c, nil).Debug("Store timeout job to background for appNumber: " + acc.Application.ApplicationNumber)
		time.Sleep(delay * time.Second)
		logger.Make(c, nil).Debug("Start to make appNumber: " + acc.Application.ApplicationNumber + " expired!")

		if err := reg.regRepo.UpdateAppStatusTimeout(c, acc.Application); err == nil {
			// send notif
			notif.GcApplicationExpired(acc)
			_ = reg.rrr.SendNotification(c, notif, "mobile")
		}
	}()
}

func (reg *registrationsUseCase) CheckApplicationByCIF(c echo.Context, pl interface{}) models.Applications {
	r := reflect.ValueOf(pl)
	cif := r.FieldByName("CIF")

	if cif.IsZero() {
		return models.Applications{}
	}

	app, err := reg.regRepo.GetAppByCIF(cif.String())

	if err != nil {
		return models.Applications{}
	}

	return app
}

func (reg *registrationsUseCase) appFailedNotification(c echo.Context, acc models.Account) error {
	var notif models.PdsNotification
	appProcessExisted, err := reg.phUC.IsProcessedAppExisted(c, acc)

	if err != nil {
		return err
	}

	if appProcessExisted {
		return nil
	}

	notif.GcApplication(acc, "failed")

	return reg.rrr.SendNotification(c, notif, "")
}
