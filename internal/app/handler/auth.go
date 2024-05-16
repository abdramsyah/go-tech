package handler

import (
	"emoney-backoffice/internal/app/dto"
	"emoney-backoffice/internal/app/util"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	HandlerOption
}

func (h AuthHandler) Login(c echo.Context) (status int, resp dto.HttpResponse) {
	var err error
	req := new(dto.LoginRequest)
	if err = c.Bind(req); err != nil {
		h.HandlerOption.Options.Logger.Error("Error bind request",
			zap.Error(err),
		)
		status = http.StatusBadRequest
		resp = dto.FailedHttpResponse("", "Error bind request", nil)
		return
	}

	err = req.Validate()
	if err != nil {
		status = http.StatusBadRequest
		resp = dto.FailedHttpResponse("", err.Error(), nil)
		return
	}

	status, token, err := h.Services.Auth.Login(c, req)
	if err != nil {
		resp = dto.FailedHttpResponse("", err.Error(), nil)
		return
	}

	resp = dto.SuccessHttpResponse("", "Login succeed", token)
	return
}

func (h AuthHandler) Logout(c echo.Context) (status int, resp dto.HttpResponse) {
	var err error
	actx, err := util.NewAppContext(c)
	if err != nil {
		resp = dto.FailedHttpResponse("", err.Error(), nil)
		return
	}

	adminID := actx.GetAdminID()
	accessUUID := actx.GetAccessUUID()
	status, err = h.Services.Auth.Logout(c, adminID, accessUUID)
	if err != nil {
		resp = dto.FailedHttpResponse("", err.Error(), nil)
		return
	}

	resp = dto.SuccessHttpResponse("", "Logout succeed", nil)
	return
}

func (h AuthHandler) RefreshToken(c echo.Context) (status int, resp dto.HttpResponse) {
	var err error
	req := new(dto.RefreshTokenRequest)
	if err = c.Bind(req); err != nil {
		h.HandlerOption.Options.Logger.Error("Error bind request",
			zap.Error(err),
		)
		status = http.StatusBadRequest
		resp = dto.FailedHttpResponse("", "Error bind request", nil)
		return
	}

	err = req.Validate()
	if err != nil {
		status = http.StatusBadRequest
		resp = dto.FailedHttpResponse("", err.Error(), nil)
		return
	}

	status, token, err := h.Services.Auth.RefreshToken(c, req.RefreshToken)
	if err != nil {
		resp = dto.FailedHttpResponse("", err.Error(), nil)
		return
	}

	resp = dto.SuccessHttpResponse("", "Token refreshed", token)
	return
}
