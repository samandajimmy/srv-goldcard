package usecase

import (
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/api"
	"srv-goldcard/internal/pkg/logger"

	"github.com/labstack/echo"
)

func (reg *registrationsUseCase) briApply(c echo.Context, acc *model.Account, pl model.PayloadBriRegister) error {
	err := reg.briRegister(c, acc, pl)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		// insert error to process handler
		// ubah status error jadi true di table process_statuses
		go func() { _ = reg.phUC.UpsertAppProcess(c, acc, err.Error()) }()

		return err
	}

	// upload document to BRI API
	go func() {
		err := reg.uploadAppDocs(c, acc)

		if err != nil {
			return
		}

		_ = reg.appNotification(c, *acc, "succeeded", false)
	}()

	return nil
}

func (reg *registrationsUseCase) briRegister(c echo.Context, acc *model.Account, pl model.PayloadBriRegister) error {
	if acc.BrixKey != "" {
		return nil
	}

	// mapping and validate bri register specification
	_ = pl.ValidateBRIRegisterSpecification()

	if err := c.Validate(pl); err != nil {
		logger.Make(c, nil).Debug(model.ErrValidateBRIRegSpec)

		return model.ErrValidateBRIRegSpec
	}

	resp := api.BriResponse{}
	reqBody := api.BriRequest{RequestData: pl}
	err := api.RetryableBriPost(c, "/register", reqBody, &resp)

	if err != nil {
		return err
	}

	// update brixkey id
	if _, ok := resp.DataOne["briXkey"].(string); !ok {
		logger.Make(c, nil).Debug(model.ErrSetVar)

		return model.ErrSetVar
	}

	acc.BrixKey = resp.DataOne["briXkey"].(string)
	// concurrently update brixkey from BRI API
	go func() {
		_ = reg.regRepo.UpdateBrixkeyID(c, *acc)
	}()

	return nil
}

func (reg *registrationsUseCase) uploadAppDocs(c echo.Context, acc *model.Account) error {
	// concurrently upload application documents to BRI
	var errors []error

	for _, doc := range acc.Application.Documents {
		err := reg.UploadAppDoc(c, acc.BrixKey, doc)

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		// insert error to process handler
		// ubah status error jadi true di table proess_statuses
		go func() { _ = reg.phUC.UpsertAppProcess(c, acc, errors[0].Error()) }()

		return errors[0]
	}

	return nil
}

func (reg *registrationsUseCase) UploadAppDoc(c echo.Context, brixkey string, doc model.Document) error {
	if doc.DocID != "" {
		return nil
	}

	briReq := model.AppDocument{
		BriXkey:    brixkey,
		DocType:    model.MapBRIDocType[doc.Type],
		FileName:   doc.FileName,
		FileExt:    doc.FileExtension,
		Base64file: model.MapBRIExtBase64File[doc.FileExtension] + doc.FileBase64,
	}

	resp := api.BriResponse{}
	reqBody := api.BriRequest{RequestData: briReq}
	err := api.RetryableBriPost(c, "/document", reqBody, &resp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if _, ok := resp.DataOne["documentId"].(string); !ok {
		return model.ErrDocIDNotFound
	}

	doc.DocID = resp.DataOne["documentId"].(string)
	// concurrently insert or update application document
	go func() {
		_ = reg.regRepo.UpsertAppDocument(c, doc)
	}()

	return nil
}

func (reg *registrationsUseCase) upsertDocument(c echo.Context, app model.Applications) error {
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

func (reg *registrationsUseCase) generateOtherDocs(c echo.Context, acc *model.Account) error {
	// Get Document (ktp, npwp, selfie, slip_te, and app_form)
	docs, err := reg.regRepo.GetDocumentByApplicationId(acc.ApplicationID, "")

	if err != nil {
		return model.ErrGetDocument
	}

	acc.Application.Documents = docs
	// Generate Application Form BRI Document
	err = reg.GenerateApplicationFormDocument(c, acc)

	if err != nil {
		return err
	}

	// Generate Slip TE Document
	err = reg.GenerateSlipTEDocument(c, acc)

	if err != nil {
		return err
	}

	return nil
}
