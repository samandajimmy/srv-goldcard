package repository

import (
	"database/sql"

	"srv-goldcard/internal/app/domain/datasync"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlSyncdataRepo struct {
	Conn *sql.DB
	DBpg *pg.DB
}

func NewDatasyncRepository(sql *sql.DB, dbpg *pg.DB) datasync.IDataSyncRp {
	return &psqlSyncdataRepo{sql, dbpg}
}

func (sdRp *psqlSyncdataRepo) RGetAllAccount(c echo.Context) ([]model.SyncActivation, error) {
	var accounts []model.SyncActivation

	query := `select a.brixkey, concat(pi2.first_name, ' ', pi2.last_name) as name, pi2.email, pi2.nik,
		a.cif, pi2.hand_phone_number, c.activated_date, c.encrypted_card_number, coalesce(histlu.card_limit, app.card_limit) credit_limit,
		app.status application_status, a.status account_status, c.status card_status
		from accounts a
		left join cards c on a.card_id = c.id
		left join personal_informations pi2 on a.personal_information_id = pi2.id
		left join applications app on a.application_id = app.id
		left join (
			select lu.account_id, lu.card_limit from limit_updates lu
			inner join (
				select lu2.account_id, max(lu2.created_at) latest_created_date
				from limit_updates lu2
				where lu2.status = 'approved'
				group by lu2.account_id
			) sublu on
				lu.account_id = sublu.account_id
				and lu.created_at = sublu.latest_created_date
		) histlu on a.id = histlu.account_id
		where a.id is not null and brixkey is not null 
		and app.status not in ('inactive', 'application_ongoing', 'rejected', 'card_suspended', 'expired')
		--and c.status != 'active'
		--and a.status = 'active'
		--and a.status != 'active' and c.status != 'active'
		order by app.application_processed_date desc
		--limit 10
	`

	_, err := sdRp.DBpg.Query(&accounts, query)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return accounts, err
	}

	return accounts, nil
}
