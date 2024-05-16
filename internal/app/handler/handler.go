package handler

import (
	"emoney-backoffice/internal/app/commons"
	"emoney-backoffice/internal/app/service"
)

// HandlerOption option for handler, including all service
type HandlerOption struct {
	commons.Options
	*service.Services
}

type Handlers struct {
	Health HealthHandler
	Auth   AuthHandler
}
