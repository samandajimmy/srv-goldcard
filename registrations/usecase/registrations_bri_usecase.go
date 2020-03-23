package usecase

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

func (reg *registrationsUseCase) briApply(c echo.Context, acc *models.Account, pl models.PayloadBriRegister) error {
	err := reg.briRegister(c, acc, pl)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	// upload document to BRI API
	err = reg.uploadAppDocs(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}

func (reg *registrationsUseCase) briRegister(c echo.Context, acc *models.Account, pl models.PayloadBriRegister) error {
	if acc.BrixKey != "" {
		return nil
	}

	// validate bri register specification
	err := pl.ValidateBRIRegisterSpecification()
	if err != nil {
		logger.Make(c, nil).Debug(models.ErrValidateBRIRegSpec)

		return models.ErrValidateBRIRegSpec
	}

	resp := api.BriResponse{}
	reqBody := api.BriRequest{RequestData: pl}
	err = api.RetryableBriPost(c, "/v1/cobranding/register", reqBody, &resp)

	if err != nil {
		return err
	}

	// update brixkey id
	if _, ok := resp.DataOne["briXkey"].(string); !ok {
		logger.Make(c, nil).Debug(models.ErrSetVar)

		return models.ErrSetVar
	}

	acc.BrixKey = resp.DataOne["briXkey"].(string)
	// concurrently update brixkey from BRI API
	go func() {
		_ = reg.regRepo.UpdateBrixkeyID(c, *acc)
	}()

	return nil
}

func (reg *registrationsUseCase) uploadAppDocs(c echo.Context, acc *models.Account) error {
	for _, doc := range acc.Application.Documents {
		// concurrently upload application documents to BRI
		go func(doc models.Document) {
			_ = reg.uploadAppDoc(c, acc.BrixKey, doc)
		}(doc)
	}

	return nil
}

func (reg *registrationsUseCase) uploadAppDoc(c echo.Context, brixkey string, doc models.Document) error {
	if doc.DocID != "" {
		return nil
	}

	briReq := models.AppDocument{
		BriXkey:    brixkey,
		DocType:    models.MapBRIDocType[doc.Type],
		FileName:   doc.FileName,
		FileExt:    doc.FileExtension,
		Base64file: models.MapBRIExtBase64File[doc.FileExtension] + doc.FileBase64,
	}

	resp := api.BriResponse{}
	reqBody := api.BriRequest{RequestData: briReq}
	err := api.RetryableBriPost(c, "/v1/cobranding/document", reqBody, &resp)

	if err != nil {
		return err
	}

	if _, ok := resp.DataOne["documentId"].(string); !ok {
		return models.ErrDocIDNotFound
	}

	doc.DocID = resp.DataOne["documentId"].(string)
	// concurrently insert or update application document
	go func() {
		_ = reg.regRepo.UpsertAppDocument(c, doc)
	}()

	return nil
}

func (reg *registrationsUseCase) upsertDocument(c echo.Context, app models.Applications) error {
	if len(app.Documents) == 0 {
		return nil
	}

	for _, doc := range app.Documents {
		err := reg.regRepo.UpsertAppDocument(c, doc)

		if err != nil {
			return err
		}
	}

	return nil
}
