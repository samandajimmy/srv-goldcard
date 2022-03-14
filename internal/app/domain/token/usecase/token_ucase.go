package usecase

import (
	"srv-goldcard/internal/app/domain/token"
	"srv-goldcard/internal/app/model"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type tokenUseCase struct {
	tokenRepo      token.Repository
	contextTimeout time.Duration
}

// NewTokenUseCase will create new an TokenUseCase object representation of Tokens.UseCase interface
func NewTokenUseCase(tkn token.Repository, timeout time.Duration) token.UseCase {
	return &tokenUseCase{
		tokenRepo:      tkn,
		contextTimeout: timeout,
	}
}

func (tkn *tokenUseCase) CreateToken(c echo.Context, accToken *model.AccountToken) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(accToken.Password), bcrypt.DefaultCost)
	accToken.Password = string(hashedPassword)
	err := tkn.tokenRepo.Create(c, accToken)

	if err != nil {
		return model.ErrCreateToken
	}

	return nil
}

func (tkn *tokenUseCase) GetToken(c echo.Context, username string, password string) (*model.AccountToken, error) {
	accToken := &model.AccountToken{}
	accToken.Username = username

	// get account
	err := tkn.tokenRepo.GetByUsername(c, accToken)

	if err != nil {
		return nil, model.ErrUsername
	}

	if err = verifyToken(accToken, password, false); err != nil {
		return nil, err
	}

	// rearrange accountToken
	accToken.ID = 0
	accToken.Password = ""
	accToken.Status = nil

	return accToken, nil
}

func (tkn *tokenUseCase) RefreshToken(c echo.Context, username string, password string) (*model.AccountToken, error) {
	accToken := &model.AccountToken{}
	accToken.Username = username

	// get account
	err := tkn.tokenRepo.GetByUsername(c, accToken)

	if err != nil {
		return nil, model.ErrUsername
	}

	if err = verifyToken(accToken, password, true); err != nil {
		return nil, err
	}

	// refresh JWT
	err = tkn.tokenRepo.UpdateToken(c, accToken)

	if err != nil {
		return nil, model.ErrCreateToken
	}

	_ = tkn.tokenRepo.GetByUsername(c, accToken)

	// rearrange accountToken
	accToken.ID = 0
	accToken.Password = ""
	accToken.Status = nil

	return accToken, nil
}

func (tkn *tokenUseCase) RefreshAllToken() error {
	// update all account token data
	err := tkn.tokenRepo.UpdateAllAccountTokenExpiry()

	if err != nil {
		return err
	}

	return nil
}

func verifyToken(accToken *model.AccountToken, password string, isUpdate bool) error {
	now := model.NowUTC()
	// validate account
	// check password
	err := bcrypt.CompareHashAndPassword([]byte(accToken.Password), []byte(password))

	if err != nil {
		return model.ErrPassword
	}

	// token availabilty
	if accToken.ExpireAt.Before(now) && !isUpdate {
		return model.ErrTokenExpired
	}

	return nil
}
