package repository

import (
	"database/sql"
	"os"
	"srv-goldcard/internal/app/domain/token"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

type psqlTokenRepository struct {
	Conn *sql.DB
}

// NewPsqlTokenRepository will create an object that represent the token.Repository interface
func NewPsqlTokenRepository(Conn *sql.DB) token.Repository {
	return &psqlTokenRepository{Conn}
}

func (m *psqlTokenRepository) Create(c echo.Context, accToken *model.AccountToken) error {
	var lastID int64
	now := time.Now()
	tokenExp := now.Add(stringToDuration(os.Getenv(`JWT_TOKEN_EXP`)) * time.Hour)
	defStatus := int64(1)

	token, err := createJWTToken(accToken, now, tokenExp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	accToken.Token = token
	accToken.ExpireAt = &tokenExp
	accToken.CreatedAt = &now
	accToken.Status = &defStatus

	query := `INSERT INTO account_tokens (username, password, token, expire_at, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	err = stmt.QueryRow(accToken.Username, accToken.Password, accToken.Token,
		accToken.ExpireAt, accToken.Status, accToken.CreatedAt).Scan(&lastID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	accToken.ID = lastID
	accToken.ExpiresAt = time.Duration(time.Until(*accToken.ExpireAt).Seconds())
	return nil
}

func (m *psqlTokenRepository) GetByUsername(c echo.Context, accToken *model.AccountToken) error {
	var expireAt, updatedAt, createdAt pq.NullTime
	query := `SELECT id, username, password, token, expire_at, status, updated_at, created_at
		FROM account_tokens
		WHERE status = 1 AND username = $1`

	err := m.Conn.QueryRow(query, accToken.Username).Scan(
		&accToken.ID, &accToken.Username, &accToken.Password,
		&accToken.Token, &expireAt, &accToken.Status,
		&updatedAt, &createdAt,
	)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	accToken.ExpireAt = &expireAt.Time
	accToken.CreatedAt = &createdAt.Time
	accToken.UpdatedAt = &updatedAt.Time
	accToken.ExpiresAt = time.Duration(accToken.ExpireAt.Sub(model.NowUTC()).Seconds())

	return nil
}

func (m *psqlTokenRepository) UpdateToken(c echo.Context, accToken *model.AccountToken) error {
	var ID int64
	now := time.Now()
	tokenExp := now.Add(stringToDuration(os.Getenv(`JWT_TOKEN_EXP`)) * time.Hour)
	token, err := createJWTToken(accToken, now, tokenExp)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	query := `UPDATE account_tokens SET token = $1, expire_at = $2, updated_at = $3 WHERE username = $4 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	err = stmt.QueryRow(token, tokenExp, now, accToken.Username).Scan(&ID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (m *psqlTokenRepository) UpdateAllAccountTokenExpiry() error {
	query := `UPDATE account_tokens SET expire_at = $1, updated_at = $2 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	_, err = stmt.Query(nil, time.Now())

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return err
	}

	return nil
}

func stringToDuration(str string) time.Duration {
	hours, _ := strconv.Atoi(str)

	return time.Duration(hours)
}

func createJWTToken(accountToken *model.AccountToken, now time.Time, tokenExp time.Time) (string, error) {
	token := model.Token{
		Name:   accountToken.Username,
		Claims: jwt.StandardClaims{Id: accountToken.Username, ExpiresAt: tokenExp.Unix()},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, token.Claims)
	return rawToken.SignedString([]byte(os.Getenv(`JWT_SECRET`)))
}
