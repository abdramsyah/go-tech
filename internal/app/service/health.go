package service

import (
	"emoney-backoffice/internal/app/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

type IHealthService interface {
	CheckHealth(ctx echo.Context) (status int, resp dto.HealthResponse)
}

type healthService struct {
	opt Option
}

func NewHealthService(opt Option) IHealthService {
	return &healthService{
		opt: opt,
	}
}

func (s *healthService) CheckHealth(ctx echo.Context) (status int, resp dto.HealthResponse) {
	status = http.StatusOK
	resp = dto.HealthResponse{Message: "OK"}
	return
}
