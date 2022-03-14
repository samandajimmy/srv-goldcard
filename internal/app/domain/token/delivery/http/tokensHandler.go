package http

import (
	"srv-goldcard/internal/app/domain/token"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// TokensHandler represent the httphandler for tokens
type TokensHandler struct {
	response     model.Response
	respErrors   model.ResponseErrors
	TokenUseCase token.UseCase
}

// NewTokensHandler represent to register tokens endpoint
func NewTokensHandler(echoGroup model.EchoGroup, tknUseCase token.UseCase) {
	handler := &TokensHandler{
		TokenUseCase: tknUseCase,
	}

	echoGroup.Token.POST("/create", handler.createToken)
	echoGroup.Token.GET("/get", handler.getToken)
	echoGroup.Token.GET("/refresh", handler.refreshToken)
}

func (tkn *TokensHandler) createToken(c echo.Context) error {
	var accountToken model.AccountToken
	tkn.response, tkn.respErrors = model.NewResponse()
	err := c.Bind(&accountToken)

	if err != nil {
		tkn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		tkn.response.SetResponse("", &tkn.respErrors)

		return tkn.response.Body(c, err)
	}

	err = tkn.TokenUseCase.CreateToken(c, &accountToken)

	if err != nil {
		tkn.respErrors.SetTitle(err.Error())
		tkn.response.SetResponse("", &tkn.respErrors)

		return tkn.response.Body(c, err)
	}

	tkn.response.SetResponse(accountToken, &tkn.respErrors)

	return tkn.response.Body(c, err)
}

func (tkn *TokensHandler) getToken(c echo.Context) error {
	tkn.response, tkn.respErrors = model.NewResponse()
	var getToken model.PayloadToken

	if err := c.Bind(&getToken); err != nil {
		tkn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		tkn.response.SetResponse("", &tkn.respErrors)

		return tkn.response.Body(c, err)
	}

	accToken, err := tkn.TokenUseCase.GetToken(c, getToken.UserName, getToken.Password)

	if err != nil {
		tkn.respErrors.SetTitle(err.Error())
		tkn.response.SetResponse("", &tkn.respErrors)

		return tkn.response.Body(c, err)
	}

	tkn.response.SetResponse(accToken, &tkn.respErrors)
	return tkn.response.Body(c, err)
}

func (tkn *TokensHandler) refreshToken(c echo.Context) error {
	tkn.response, tkn.respErrors = model.NewResponse()
	var refToken model.PayloadToken
	if err := c.Bind(&refToken); err != nil {
		tkn.respErrors.SetTitle(model.MessageUnprocessableEntity)
		tkn.response.SetResponse("", &tkn.respErrors)

		return tkn.response.Body(c, err)
	}

	accToken, err := tkn.TokenUseCase.RefreshToken(c, refToken.UserName, refToken.Password)

	if err != nil {
		tkn.respErrors.SetTitle(err.Error())
		tkn.response.SetResponse("", &tkn.respErrors)

		return tkn.response.Body(c, err)
	}

	tkn.response.SetResponse(accToken, &tkn.respErrors)

	return tkn.response.Body(c, err)
}
