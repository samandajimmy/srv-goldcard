package repository

import (
	"database/sql"
	"fmt"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
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
		gcdb.NewPipelineStmt(`INSERT INTO applications (application_number, status, created_at)
			VALUES ($1, $2, $3) RETURNING id;`,
			[]string{"appID"}, app.ApplicationNumber, app.Status, time.Now()),
		gcdb.NewPipelineStmt(`INSERT INTO personal_informations (hand_phone_number, created_at)
			VALUES ($1, $2) RETURNING id;`,
			[]string{"piID"}, pi.HandPhoneNumber, time.Now()),
		gcdb.NewPipelineStmt(`INSERT INTO accounts (cif, branch_code, bank_id, emergency_contact_id,
			created_at, application_id, personal_information_id) VALUES ($1, $2, $3, $4, $5, {appID}, {piID});`,
			nilFilters, acc.CIF, acc.BranchCode, acc.BankID, acc.EmergencyContactID, time.Now()),
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
	var nilFilters []string
	corr := acc.Correspondence

	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`INSERT INTO correspondences (address_line_1, address_line_2,
			address_line_3, address_city, zipcode, created_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
			[]string{"corrID"}, corr.AddressLine1, corr.AddressLine2, corr.AddressLine3,
			corr.AddressCity, corr.Zipcode, time.Now()),
		gcdb.NewPipelineStmt(`UPDATE accounts set correspondence_id = {corrID}, updated_at = $1
			WHERE id = $2`,
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

func (regis *psqlRegistrationsRepository) PostSavingAccount(c echo.Context, acc models.Account) error {
	query := `UPDATE applications set saving_account = $1, updated_at = $2
		WHERE id = $3;`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	_, err = stmt.Exec(acc.Application.SavingAccount, time.Now(), acc.ApplicationID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetAccountByAppNumber(c echo.Context, acc *models.Account) error {
	newAcc := models.Account{}
	err := regis.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Where("application_number = ?", acc.Application.ApplicationNumber).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		return models.ErrAppNumberNotFound
	}

	*acc = newAcc
	return nil
}

func (regis *psqlRegistrationsRepository) GetAllRegData(c echo.Context, appNumber string) (models.PayloadPersonalInformation, error) {
	var plRegister models.PayloadPersonalInformation
	var pi models.PersonalInformation

	query := `select acc.product_request, acc.billing_cycle, acc.card_deliver, c.card_name,
		pi.first_name, pi.last_name, pi.hand_phone_number, pi.email, pi.npwp, pi.nik, pi.birth_place,
		pi.birth_date, pi.nationality, pi.sex, pi.education, pi.marital_status, pi.mother_name,
		pi.home_phone_area, pi.home_phone_number, pi.home_status, pi.address_line_1, pi.address_line_2,
		pi.address_line_3, pi.zipcode, pi.address_city, pi.stayed_since, pi.child, o.job_bidang_usaha,
		o.job_sub_bidang_usaha, o.job_category, o.job_status, o.total_employee, o.company, o.job_title,
		o.work_since, o.office_address_1, o.office_address_2, o.office_address_3, o.office_zipcode,
		o.office_city, o.office_phone, o.income, ec.name emergency_name, ec.relation emergency_relation,
		ec.phone_number emergency_phone_number, ec.address_line_1 emergency_address_1,
		ec.address_line_2 emergency_address_2, ec.address_line_3 emergency_address_3,
		ec.address_city emergency_city, ec.zipcode emergency_zipcode, corr.address_line_1,
		corr.address_line_2, corr.address_line_3, corr.address_city, corr.zipcode
		from accounts acc
		left join applications app on acc.application_id = app.id
		left join cards c on acc.card_id = c.id
		left join correspondences corr on acc.correspondence_id = corr.id
		left join emergency_contacts ec on acc.emergency_contact_id = ec.id
		left join occupations o on acc.occupation_id = o.id
		left join personal_informations pi on acc.personal_information_id = pi.id
		where app.status = ? and app.application_number = ?;`

	_, err := regis.DBpg.QueryOne(&plRegister, query, models.AppStatusOngoing, appNumber)

	if err != nil || (plRegister == models.PayloadPersonalInformation{}) {
		return plRegister, err
	}

	plRegister.Sex = pi.GetBriSex(plRegister.SexString)
	plRegister.SexString = ""

	return plRegister, nil
}

func (regis *psqlRegistrationsRepository) UpdateAllRegistrationData(c echo.Context, acc models.Account) error {
	var nilFilters []string
	// occ := acc.Occupation
	app := acc.Application
	pi := acc.PersonalInformation

	stmts := []*gcdb.PipelineStmt{
		// update card
		gcdb.NewPipelineStmt("UPDATE cards SET card_name = $1, updated_at = $2 WHERE id = $3;",
			nilFilters, acc.Card.CardName, time.Now(), acc.CardID),
		// insert occupation
		// gcdb.NewPipelineStmt(`INSERT INTO occupations (job_bidang_usaha, job_sub_bidang_usaha,
		// 	job_category, job_status, total_employee, company, job_title, work_since,
		// 	office_address_1, office_address_2, office_address_3, office_zipcode, office_city,
		// 	office_phone, income, created_at)
		// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id;`,
		// 	[]string{"occID"}, occ.JobBidangUsaha, occ.JobSubBidangUsaha, occ.JobCategory,
		// 	occ.JobStatus, occ.TotalEmployee, occ.Company, occ.JobTitle, occ.WorkSince,
		// 	occ.OfficeAddress1, occ.OfficeAddress2, occ.OfficeAddress3, occ.OfficeZipcode,
		// 	occ.OfficeCity, occ.OfficePhone, occ.Income, time.Now()),
		// update application
		gcdb.NewPipelineStmt(`UPDATE applications set ktp_image_base64 = $1, npwp_image_base64 = $2,
			selfie_image_base64 = $3, updated_at = $4 WHERE id = $5`,
			nilFilters, app.KtpImageBase64, app.NpwpImageBase64, app.SelfieImageBase64, time.Now(),
			acc.ApplicationID),
		// update personal_infomation
		gcdb.NewPipelineStmt(`UPDATE personal_informations set first_name = $1, last_name = $2,
			email = $3, npwp = $4, nik = $5, birth_place = $6, birth_date = $7, nationality = $8,
			sex = $9, education = $10, marital_status = $11, mother_name = $12, home_phone_area = $13,
			home_phone_number = $14, home_status = $15, address_line_1 = $16, address_line_2 = $17,
			address_line_3 = $18, zipcode = $19, address_city = $20, stayed_since = $21, child = $22,
			updated_at = $23 WHERE id = $24`, nilFilters, pi.FirstName, pi.LastName, pi.Email,
			pi.Npwp, pi.Nik, pi.BirthPlace, pi.BirthDate, pi.Nationality, "male", pi.Education,
			pi.MaritalStatus, pi.MotherName, pi.HomePhoneArea, pi.HandPhoneNumber, pi.HomeStatus,
			pi.AddressLine1, pi.AddressLine2, pi.AddressLine3, pi.Zipcode, pi.AddressCity,
			pi.StayedSince, pi.Child, time.Now(), acc.PersonalInformationID),
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
			office_address_1, office_address_2, office_address_3, office_zipcode, office_city,
			office_phone, income, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id;`,
			[]string{"occID"}, occ.JobBidangUsaha, occ.JobSubBidangUsaha, occ.JobCategory,
			occ.JobStatus, occ.TotalEmployee, occ.Company, occ.JobTitle, occ.WorkSince,
			occ.OfficeAddress1, occ.OfficeAddress2, occ.OfficeAddress3, occ.OfficeZipcode,
			occ.OfficeCity, occ.OfficePhone, occ.Income, time.Now()),
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

func (regis *psqlRegistrationsRepository) GetCityFromZipcode(c echo.Context, acc models.Account) (string, string, error) {
	var city string
	zipcode := acc.Occupation.OfficeZipcode

	query := `SELECT city FROM ref_postal_codes pc
		WHERE pc.postal_code = $1 LIMIT 1`

	err := regis.Conn.QueryRow(query, acc.Occupation.OfficeZipcode).Scan(&city)

	if err != nil && err != sql.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", "", err
	}

	if err == sql.ErrNoRows || acc.Occupation.OfficeZipcode == "" {
		city = acc.PersonalInformation.AddressCity
		zipcode = acc.PersonalInformation.Zipcode
	}

	return city, zipcode, nil
}

func (regis *psqlRegistrationsRepository) UpdateCardLimit(c echo.Context, acc models.Account) error {
	var nilFilters []string

	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`INSERT INTO cards (card_limit, created_at)
			VALUES ($1, $2) RETURNING id;`,
			[]string{"cardID"}, acc.Card.CardLimit, time.Now()),
		gcdb.NewPipelineStmt(`UPDATE accounts set card_id = {cardID}, updated_at = $1
			WHERE id = $2`,
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

func (regis *psqlRegistrationsRepository) UpdateBrixkeyID(c echo.Context, acc models.Account) error {
	acc.UpdatedAt = time.Now()
	_, err := regis.DBpg.Model(&acc).
		Set(`brixkey = ?brixkey, status = ?status, updated_at = ?updated_at`).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) UpdateGetAppStatus(c echo.Context, app models.Applications) (models.AppStatus, error) {
	var appStatus models.AppStatus
	app.UpdatedAt = time.Now()
	key := app.GetStatusDateKey()

	_, err := regis.DBpg.Model(&app).
		Set(fmt.Sprintf(`status = ?status, %s = ?%s, updated_at = ?updated_at`, key, key)).
		WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return appStatus, err
	}

	query := `select status, application_processed_date, card_processed_date, card_send_date,
		card_sent_date, failed_date from applications where id = ?;`

	_, err = regis.DBpg.Query(&appStatus, query, app.ID)

	if err != nil || (appStatus == models.AppStatus{}) {
		logger.Make(c, nil).Debug(err)

		return appStatus, err
	}

	return appStatus, nil
}

func (regis *psqlRegistrationsRepository) UpdateAppDocID(c echo.Context, app models.Applications) error {
	app.UpdatedAt = time.Now()
	_, err := regis.DBpg.Model(&app).
		Set(`ktp_doc_id = ?ktp_doc_id, npwp_doc_id = ?npwp_doc_id, selfie_doc_id = ?selfie_doc_id,
			updated_at = ?updated_at`).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetAppByID(c echo.Context, appID int64) (models.Applications, error) {
	var app models.Applications
	query := `select id, application_number, status, ktp_image_base64, npwp_image_base64,
		selfie_image_base64, saving_account, created_at, updated_at from applications where id = ?;`

	_, err := regis.DBpg.Query(&app, query, appID)

	if err != nil || (app == models.Applications{}) {
		logger.Make(c, nil).Debug(err)

		return app, err
	}

	return app, nil
}

func (regis *psqlRegistrationsRepository) UpdateApplication(c echo.Context, app models.Applications, col []string) error {
	app.UpdatedAt = time.Now()
	col = append(col, "updated_at")
	_, err := regis.DBpg.Model(&app).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
