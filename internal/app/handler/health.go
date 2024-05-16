package handler

import (
	"emoney-backoffice/internal/app/constant"
	"emoney-backoffice/internal/app/dto"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	HandlerOption
}

func (h HealthHandler) Check(ctx echo.Context) (status int, resp dto.HttpResponse) {
	status, healthResp := h.Services.Health.CheckHealth(ctx)
	resp = dto.HttpResponse{
		Status:  constant.RespSuccessStatus,
		Code:    "-",
		Message: healthResp.Message,
	}
	return
}
