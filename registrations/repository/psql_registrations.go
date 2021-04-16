package repository

import (
	"database/sql"
	"fmt"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlRegistrationsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlRegistrationsRepository will create an object that represent the registrations.Repository interface
func NewPsqlRegistrationsRepository(Conn *sql.DB, dbpg *pg.DB) registrations.Repository {
	return &psqlRegistrationsRepository{Conn, dbpg}
}

func (regis *psqlRegistrationsRepository) CreateApplication(c echo.Context, app models.Applications,
	acc models.Account, pi models.PersonalInformation) error {
	var nilFilters []string

	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`INSERT INTO applications (application_number, status, created_at, expired_at)
			VALUES ($1, $2, $3, $4) RETURNING id;`,
			[]string{"appID"}, app.ApplicationNumber, app.Status, app.CreatedAt, app.ExpiredAt),
		gcdb.NewPipelineStmt(`INSERT INTO personal_informations (hand_phone_number, created_at)
			VALUES ($1, $2) RETURNING id;`,
			[]string{"piID"}, pi.HandPhoneNumber, time.Now()),
		gcdb.NewPipelineStmt(`INSERT INTO accounts (cif, branch_code, product_request, billing_cycle,
			card_deliver, bank_id, emergency_contact_id, created_at, application_id, personal_information_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, {appID}, {piID});`,
			nilFilters, acc.CIF, acc.BranchCode, acc.ProductRequest, acc.BillingCycle, acc.CardDeliver,
			acc.BankID, acc.EmergencyContactID, time.Now()),
	}

	err := gcdb.WithTransaction(regis.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetEmergencyContactIDByType(c echo.Context, typeDef string) (int64, error) {
	var ecID int64
	query := `SELECT id FROM emergency_contacts WHERE type = $1 LIMIT 1`
	err := regis.Conn.QueryRow(query, typeDef).Scan(&ecID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return ecID, err
	}

	return ecID, nil
}

func (regis *psqlRegistrationsRepository) GetBankIDByCode(c echo.Context, bankCode string) (int64, error) {
	var bankID int64
	query := `SELECT id FROM banks WHERE code = $1`
	err := regis.Conn.QueryRow(query, bankCode).Scan(&bankID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return bankID, err
	}

	return bankID, nil
}

func (regis *psqlRegistrationsRepository) PostAddress(c echo.Context, acc models.Account) error {
	acc.UpdatedAt = models.NowDbpg()

	// update card deliver
	_, err := regis.DBpg.Model(&acc).Column([]string{"card_deliver", "updated_at"}...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) PostSavingAccount(c echo.Context, acc models.Account) error {
	query := `UPDATE applications set saving_account = $1, saving_account_opening_date = $2, updated_at = $3
		WHERE id = $4;`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	_, err = stmt.Exec(acc.Application.SavingAccount, time.Now().Format(models.DateTimeFormatZone), time.Now(), acc.ApplicationID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetAccountByAppNumber(c echo.Context, acc *models.Account) error {
	newAcc := models.Account{}
	docs := []models.Document{}
	excludedStatus := []string{models.AppStatusInactive, models.AppStatusExpired}
	err := regis.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Relation("Occupation").Relation("EmergencyContact").
		Where("application_number = ? and application.status not in (?)",
			acc.Application.ApplicationNumber, pg.In(excludedStatus)).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		return models.ErrAppNumberNotFound
	}

	err = regis.DBpg.Model(&docs).Where("application_id = ?", newAcc.ApplicationID).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	newAcc.Application.Documents = docs

	*acc = newAcc
	return nil
}

func (regis *psqlRegistrationsRepository) GetAllRegData(c echo.Context, appNumber string) (models.PayloadBriRegister, error) {
	var plRegister models.PayloadBriRegister
	var pi models.PersonalInformation

	query := `select acc.product_request, acc.billing_cycle, acc.card_deliver, c.card_name, c.card_limit, 
		pi.first_name, pi.last_name, pi.hand_phone_number, pi.email, pi.npwp, pi.nik, pi.birth_place,
		pi.birth_date, pi.nationality, pi.sex, pi.education, pi.marital_status, pi.mother_name,
		pi.home_phone_area, pi.home_phone_number, pi.home_status, pi.stayed_since, pi.child, o.job_bidang_usaha,
		o.job_sub_bidang_usaha, o.job_category, o.job_status, o.total_employee, o.company, o.job_title,
		o.work_since, o.office_address_1, o.office_address_2, o.office_address_3, o.office_zipcode,
		o.office_city, o.office_phone, o.income, ec.name emergency_name, ec.relation emergency_relation,
		ec.phone_number emergency_phone_number, ec.address_line_1 emergency_address_1,
		ec.address_line_2 emergency_address_2, ec.address_line_3 emergency_address_3,
		ec.address_city emergency_city, ec.zipcode emergency_zipcode,
		pi.address_line_1,
		pi.address_line_2,
		pi.address_line_3,
		pi.address_city,
		pi.zipcode
		from accounts acc
		left join applications app on acc.application_id = app.id
		left join cards c on acc.card_id = c.id
		left join emergency_contacts ec on acc.emergency_contact_id = ec.id
		left join occupations o on acc.occupation_id = o.id
		left join personal_informations pi on acc.personal_information_id = pi.id
		where app.application_number = ?;`

	_, err := regis.DBpg.QueryOne(&plRegister, query, appNumber)

	if err != nil || (plRegister == models.PayloadBriRegister{}) {
		return plRegister, err
	}

	plRegister.Sex = pi.GetSexInt(plRegister.SexString)
	plRegister.SexString = ""

	return plRegister, nil
}

func (regis *psqlRegistrationsRepository) UpdateAllRegistrationData(c echo.Context, acc models.Account) error {
	var nilFilters []string
	pi := acc.PersonalInformation

	stmts := []*gcdb.PipelineStmt{
		// update card
		gcdb.NewPipelineStmt("UPDATE cards SET card_name = $1, updated_at = $2 WHERE id = $3;",
			nilFilters, acc.Card.CardName, time.Now(), acc.CardID),
		// update personal_infomation
		gcdb.NewPipelineStmt(`UPDATE personal_informations set first_name = $1, last_name = $2,
			email = $3, npwp = $4, nik = $5, birth_place = $6, birth_date = $7, nationality = $8,
			sex = $9, education = $10, marital_status = $11, mother_name = $12, home_phone_area = $13,
			home_phone_number = $14, home_status = $15, address_line_1 = $16, address_line_2 = $17,
			address_line_3 = $18, zipcode = $19, address_city = $20, stayed_since = $21, child = $22,
			updated_at = $23, relative_phone_number = $24, address_province = $25, address_subdistrict = $26, 
			address_village = $27 WHERE id = $28`, nilFilters, pi.FirstName, pi.LastName, pi.Email,
			pi.Npwp, pi.Nik, pi.BirthPlace, pi.BirthDate, pi.Nationality, "male", pi.Education,
			pi.MaritalStatus, pi.MotherName, pi.HomePhoneArea, pi.HomePhoneNumber, pi.HomeStatus,
			pi.AddressLine1, pi.AddressLine2, pi.AddressLine3, pi.Zipcode, pi.AddressCity,
			pi.StayedSince, pi.Child, time.Now(), pi.RelativePhoneNumber, pi.AddressProvince,
			pi.AddressSubdistrict, pi.AddressVillage, acc.PersonalInformationID),
		// update account
		gcdb.NewPipelineStmt(`UPDATE accounts set product_request = $1, billing_cycle = $2,
			card_deliver = $3, updated_at = $4 WHERE id = $5`,
			nilFilters, acc.ProductRequest, acc.BillingCycle, acc.CardDeliver, time.Now(),
			acc.ID),
	}

	err := gcdb.WithTransaction(regis.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) PostOccupation(c echo.Context, acc models.Account) error {
	var nilFilters []string
	occ := acc.Occupation

	stmts := []*gcdb.PipelineStmt{
		// insert occupation
		gcdb.NewPipelineStmt(`INSERT INTO occupations (job_bidang_usaha, job_sub_bidang_usaha,
			job_category, job_status, total_employee, company, job_title, work_since,
			office_address_1, office_address_2, office_address_3, office_zipcode, office_province,
			office_city, office_subdistrict, office_village, office_phone, income, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) RETURNING id;`,
			[]string{"occID"}, occ.JobBidangUsaha, occ.JobSubBidangUsaha, occ.JobCategory,
			occ.JobStatus, occ.TotalEmployee, occ.Company, occ.JobTitle, occ.WorkSince,
			occ.OfficeAddress1, occ.OfficeAddress2, occ.OfficeAddress3, occ.OfficeZipcode,
			occ.OfficeProvince, occ.OfficeCity, occ.OfficeSubdistrict, occ.OfficeVillage,
			occ.OfficePhone, occ.Income, time.Now()),
		// update account
		gcdb.NewPipelineStmt(`UPDATE accounts set occupation_id = {occID}, updated_at = $1 WHERE id = $2`,
			nilFilters, time.Now(), acc.ID),
	}

	err := gcdb.WithTransaction(regis.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetZipcode(c echo.Context, addrData models.AddressData) (string, error) {
	var zipcode string
	query := `SELECT postal_code FROM ref_postal_codes pc
		LEFT JOIN ref_provinces p ON pc.province_code = p.province_code
		WHERE p.province_name = $1 AND pc.city = $2 AND pc.sub_district = $3 AND pc.village = $4
		LIMIT 1`

	err := regis.Conn.QueryRow(query, strings.ToUpper(addrData.Province), strings.ToUpper(addrData.City),
		strings.ToUpper(addrData.Subdistrict), strings.ToUpper(addrData.Village)).Scan(&zipcode)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return zipcode, err
	}

	return zipcode, nil
}

func (regis *psqlRegistrationsRepository) GetCityFromZipcode(c echo.Context, zipcode string) (models.AddressData, error) {
	addrData := models.AddressData{}
	query := `SELECT city, province_name as province, sub_district, village FROM ref_postal_codes pc
		left join ref_provinces p on pc.province_code = p.province_code WHERE pc.postal_code = ? LIMIT 1`

	_, err := regis.DBpg.QueryOne(&addrData, query, zipcode)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return addrData, err
	}

	addrData.Zipcode = zipcode
	return addrData, err
}

func (regis *psqlRegistrationsRepository) UpdateCardLimit(c echo.Context, acc models.Account, fnAfter func() error) error {
	var nilFilters []string
	var upsertFilters []string
	var upsertQueryCard string
	var stmts []*gcdb.PipelineStmt

	// query for table cards if account has not any card data yet
	upsertFilters = []string{"cardID"}
	upsertQueryCard = `INSERT INTO cards (card_limit, created_at, gold_limit, stl_limit, balance,
		gold_balance, stl_balance, previous_card_balance, previous_card_balance_date, previous_card_limit,
		previous_card_limit_date) VALUES ($1, $2, $3, $4, $1, $3, $4, $1, $2, $1, $2) RETURNING id;`

	// query for table cards if account has any card data then update
	if acc.CardID != 0 {
		upsertQueryCard = `UPDATE cards set card_limit = $1, updated_at = $2, gold_limit = $3,
			stl_limit = $4, balance = $1, gold_balance = $3, stl_balance = $4 WHERE id = ` +
			strconv.Itoa(int(acc.CardID)) + ` RETURNING id;`
	}

	stmts = []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(upsertQueryCard, upsertFilters, acc.Card.CardLimit, time.Now(),
			acc.Card.GoldLimit, acc.Card.StlLimit),
		gcdb.NewPipelineStmt(`UPDATE accounts set card_id = {cardID}, updated_at = $1 WHERE id = $2`,
			nilFilters, time.Now(), acc.ID),
		gcdb.NewPipelineStmt(`UPDATE applications set card_limit = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, acc.Application.CardLimit, time.Now(), acc.ApplicationID),
	}

	err := gcdb.WithTransaction(regis.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if fnAfter != nil {
		err = fnAfter()

		if err != nil {
			return err
		}
	}

	return nil
}

func (regis *psqlRegistrationsRepository) UpdateBrixkeyID(c echo.Context, acc models.Account) error {
	acc.UpdatedAt = models.NowDbpg()
	_, err := regis.DBpg.Model(&acc).
		Set(`brixkey = ?brixkey, status = ?status, updated_at = ?updated_at`).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetAppStatus(c echo.Context, app models.Applications) (models.AppStatus, error) {
	var appStatus models.AppStatus

	query := `select status, application_processed_date, card_processed_date, card_send_date,
		card_sent_date, rejected_date from applications where id = ?;`

	_, err := regis.DBpg.Query(&appStatus, query, app.ID)

	if err != nil || (appStatus == models.AppStatus{}) {
		logger.Make(c, nil).Debug(err)

		return appStatus, err
	}

	return appStatus, nil
}

func (regis *psqlRegistrationsRepository) UpdateAppStatusTimeout(c echo.Context, app models.Applications) error {
	appCheck := models.Applications{}
	app.UpdatedAt = models.NowDbpg()
	app.Status = models.AppStatusExpired

	err := regis.DBpg.Model(&appCheck).
		Where("application_number = ? and status = ?", app.ApplicationNumber, models.AppStatusOngoing).Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return models.ErrGetParameter
	}

	_, err = regis.DBpg.Model(&app).Set("status = ?status, updated_at = ?updated_at").
		Where("application_number = ? and status = ?", app.ApplicationNumber, models.AppStatusOngoing).
		Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) ForceUpdateAppStatusTimeout() error {
	app := models.Applications{}

	_, err := regis.DBpg.Model(app).Exec(`UPDATE applications SET status = ?, updated_at = ?
		where expired_at <= ? and status = ?`, models.AppStatusExpired, models.NowDbpg(), models.NowDbpg(),
		models.AppStatusOngoing)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) UpdateAppStatus(c echo.Context, app models.Applications) error {
	app.UpdatedAt = models.NowDbpg()
	key := app.GetStatusDateKey()
	dynQuery := "status = ?status, updated_at = ?updated_at"

	if key != "" {
		dynQuery = fmt.Sprintf(dynQuery+", %s = ?%s", key, key)
	}

	_, err := regis.DBpg.Model(&app).Set(dynQuery).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) ResetAppStatusToCardProcessed(appsId int64) error {
	app := models.Applications{}

	_, err := regis.DBpg.Model(app).Exec(`UPDATE applications SET status = ?, card_processed_date = ?, card_sent_date = null, card_send_date = null, 
		rejected_date = null where id = ?`, models.AppStatusCardProcessed, models.NowDbpg(), appsId)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) UpdateAppDocID(c echo.Context, app models.Applications) error {
	app.UpdatedAt = models.NowDbpg()
	_, err := regis.DBpg.Model(&app).
		Set(`ktp_doc_id = ?ktp_doc_id, npwp_doc_id = ?npwp_doc_id, selfie_doc_id = ?selfie_doc_id,
			updated_at = ?updated_at`).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) UpdateApplication(c echo.Context, app models.Applications, col []string) error {
	app.UpdatedAt = models.NowDbpg()
	col = append(col, "updated_at")
	_, err := regis.DBpg.Model(&app).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) UpsertAppDocument(c echo.Context, doc models.Document) error {
	var err error

	if doc.ID == 0 {
		err = regis.insertAppDocument(c, doc)
	} else {
		doc.UpdatedAt = models.NowDbpg()
		err = regis.DBpg.Update(&doc)
	}

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetCoreServiceStatus(c echo.Context) error {
	var param models.Parameter
	query := `select id, key, value, description, created_at, updated_at
		from parameters where key = ? limit 1;`

	_, err := regis.DBpg.Query(&param, query, "CORE_EOD_HEALTH")

	if err != nil || (param == models.Parameter{}) {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	coreStatus, _ := strconv.ParseBool(param.Value)

	if !coreStatus {
		return models.ErrCoreEODStatus
	}

	return nil
}

// UpdateCoreOpen for update field core open on application table
func (regis *psqlRegistrationsRepository) UpdateCoreOpen(c echo.Context, acc *models.Account) error {
	app := models.Applications{
		CoreOpen:  true,
		UpdatedAt: models.NowDbpg(),
	}
	_, err := regis.DBpg.Model(&app).Column("core_open", "updated_at").
		Where("application_number = ?", acc.Application.ApplicationNumber).Update()

	if err != nil {
		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) insertAppDocument(c echo.Context, doc models.Document) error {
	doc.CreatedAt = models.NowDbpg()
	err := regis.DBpg.Insert(&doc)

	if err != nil {
		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetDocumentByApplicationId(appId int64, docType string) ([]models.Document, error) {
	var listDocument []models.Document
	docTypes := models.DocTypes

	if docType != "" {
		docTypes = []string{docType}
	}

	err := regis.DBpg.Model(&listDocument).
		Where("application_id = ? AND type in (?)", appId, pg.In(docTypes)).Select()

	if err != nil || (listDocument == nil) {
		logger.Make(nil, nil).Debug(err)

		return listDocument, err
	}

	return listDocument, nil
}

func (regis *psqlRegistrationsRepository) GetSignatoryNameParam(c echo.Context) (string, error) {
	newPrm := models.Parameter{}
	err := regis.DBpg.Model(&newPrm).
		Where("key = ?", "SIGNATORY_NAME").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", models.ErrGetParameter
	}

	return newPrm.Value, nil
}

func (regis *psqlRegistrationsRepository) GetSignatoryNipParam(c echo.Context) (string, error) {
	newPrm := models.Parameter{}
	err := regis.DBpg.Model(&newPrm).
		Where("key = ?", "SIGNATORY_NIP").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", models.ErrGetParameter
	}

	return newPrm.Value, nil
}

func (regis *psqlRegistrationsRepository) DeactiveAccount(c echo.Context, acc models.Account) error {
	var nilFilters []string
	acc.Status = models.AccStatusInactive

	stmts := []*gcdb.PipelineStmt{
		// update account
		gcdb.NewPipelineStmt("UPDATE accounts SET status = $1, updated_at = $2 WHERE id = $3;",
			nilFilters, acc.Status, time.Now(), acc.ID),
		// update application
		gcdb.NewPipelineStmt(`UPDATE applications set status = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, models.AppStatusInactive, time.Now(), acc.Application.ID),
	}

	err := gcdb.WithTransaction(regis.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetAppOngoing() ([]models.Account, error) {
	accs := []models.Account{}

	err := regis.DBpg.Model(&accs).Relation("Application").Relation("PersonalInformation").
		Where("Application.expired_at > ? and Application.status = ?", models.NowDbpg(),
			models.AppStatusOngoing).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(nil, nil).Debug(err)

		return accs, err
	}

	return accs, nil
}

func (regis *psqlRegistrationsRepository) ForceDeliverAccount(c echo.Context, acc models.Account) error {
	query := `UPDATE applications set status = $1, card_sent_date = $2, updated_at = $2 WHERE id = $3;`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	_, err = stmt.Exec(models.AppStatusForceDeliver, time.Now(), acc.Application.ID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
