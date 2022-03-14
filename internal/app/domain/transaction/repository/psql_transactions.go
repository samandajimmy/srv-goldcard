package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	gcdb "srv-goldcard/internal/pkg/database"
	"srv-goldcard/internal/pkg/logger"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlTransactionsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlTransactionsRepository will create an object that represent the transaction.Repository interface
func NewPsqlTransactionsRepository(Conn *sql.DB, DBpg *pg.DB) transaction.Repository {
	return &psqlTransactionsRepository{Conn, DBpg}
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByBrixKey(c echo.Context, brixkey string) (model.Account, error) {
	acc := model.Account{}
	err := PSQLTrx.DBpg.Model(&acc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Where("Account.brixkey = ?", brixkey).Where("Account.status = ?", model.AccStatusActive).Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return acc, err
	}

	if err == pg.ErrNoRows {
		return acc, model.ErrGetAccByBrixkey
	}

	return acc, nil
}

func (PSQLTrx *psqlTransactionsRepository) PostTransactions(c echo.Context, trx model.Transaction) error {
	trx.CreatedAt = model.NowDbpg()
	err := PSQLTrx.DBpg.Insert(&trx)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByAccountNumber(c echo.Context, acc *model.Account) error {
	newAcc := model.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Relation("Occupation").Relation("EmergencyContact").
		Where("account_number = ? AND account.status = ?", acc.AccountNumber, model.AccStatusActive).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		return model.ErrAppNumberNotFound
	}

	*acc = newAcc
	return nil
}

func (PSQLTrx *psqlTransactionsRepository) UpdateCardBalance(c echo.Context, card model.Card) error {
	card.UpdatedAt = model.NowDbpg()
	col := []string{"balance", "gold_balance", "stl_balance", "card_limit", "previous_card_balance",
		"previous_card_balance_date", "previous_card_limit", "previous_card_limit_date", "updated_at"}

	_, err := PSQLTrx.DBpg.Model(&card).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) UpdatePayInquiryStatusPaid(c echo.Context, pay model.PaymentInquiry) error {
	pay.UpdatedAt = model.NowDbpg()
	pay.Status = model.BillTrxPaid
	col := []string{"status", "updated_at"}
	_, err := PSQLTrx.DBpg.Model(&pay).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetPayInquiryByRefTrx(c echo.Context, acc model.Account, refTrx string) (model.PaymentInquiry, error) {
	payment := model.PaymentInquiry{}
	err := PSQLTrx.DBpg.Model(&payment).Relation("Billing").
		Where("payment_inquiry.account_id = ? AND ref_trx = ?", acc.ID, refTrx).
		Where("payment_inquiry.status = ?", model.BillTrxUnpaid).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return payment, err
	}

	if err == pg.ErrNoRows {
		return payment, model.ErrGetPaymentTransaction
	}

	return payment, nil
}

func (PSQLTrx *psqlTransactionsRepository) PostPaymentInquiry(c echo.Context, paymentInq model.PaymentInquiry) error {
	paymentInq.InquiryDate = model.NowDbpg()
	paymentInq.CreatedAt = model.NowDbpg()
	err := PSQLTrx.DBpg.Insert(&paymentInq)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) PostPayment(c echo.Context, trx model.Transaction, bill model.Billing) error {
	var nilFilters []string
	billPayment := trx.BillingPayments[0]

	stmts := []*gcdb.PipelineStmt{
		// insert payment trx
		gcdb.NewPipelineStmt(`INSERT INTO transactions (account_id, ref_trx_pgdn, ref_trx, nominal,
			gold_nominal, type, status, balance, gold_balance, methods, trx_date, description,
			compare_id, transaction_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id;`,
			[]string{"trxID"}, trx.AccountId, trx.RefTrxPgdn, trx.RefTrx, trx.Nominal, trx.GoldBalance,
			trx.Type, trx.Status, trx.Balance, trx.GoldBalance, trx.Methods, trx.TrxDate, trx.Description,
			trx.CompareID, trx.TransactionID, time.Now()),
		// insert billing payment
		gcdb.NewPipelineStmt(`INSERT INTO billing_payments (trx_id, bill_id, source, created_at)
			VALUES ({trxID}, $1, $2, $3) RETURNING id;`,
			nilFilters, bill.ID, billPayment.Source, time.Now()),
		// update billings record
		gcdb.NewPipelineStmt(`UPDATE billings set debt_amount = $1,
			debt_gold = $2, debt_stl = $3, updated_at = $4 WHERE id = $5;`,
			nilFilters, bill.DebtAmount, bill.DebtGold, bill.DebtSTL, time.Now(), bill.ID),
	}

	err := gcdb.WithTransaction(PSQLTrx.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetAllActiveAccount(c echo.Context) ([]model.Account, error) {
	var listActiveAccount []model.Account
	err := PSQLTrx.DBpg.Model(&listActiveAccount).Relation("Card").Relation("PersonalInformation").
		Where("account.status = ?", model.AccStatusActive).Select()

	if err != nil || (listActiveAccount == nil) {
		logger.Make(nil, nil).Debug(err)

		return listActiveAccount, err
	}

	return listActiveAccount, nil
}

func (PSQLTrx *psqlTransactionsRepository) GetPaymentInquiryNotificationData(c echo.Context, pi model.PaymentInquiry) (model.PaymentInquiryNotificationData, error) {
	var pind model.PaymentInquiryNotificationData

	nilPind := model.PaymentInquiryNotificationData{}
	query := `SELECT id, core_response->>'reffSwitching' as reff_switching, core_response->>'administrasi' as administration
	FROM payment_inquiries where id = ? limit 1;`

	_, err := PSQLTrx.DBpg.Query(&pind, query, pi.ID)

	if err != nil || (pind == nilPind) {
		logger.Make(nil, nil).Debug(err)

		return nilPind, err
	}

	return pind, nil
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByCIF(c echo.Context, acc *model.Account) error {
	newAcc := model.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Relation("Occupation").Relation("EmergencyContact").
		Where("cif = ? AND account.status = ?", acc.CIF, model.AccStatusActive).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		return model.ErrGetAccByCIF
	}

	*acc = newAcc
	return nil
}
