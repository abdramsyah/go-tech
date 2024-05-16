package service

import (
	"emoney-backoffice/internal/app/commons"
	"emoney-backoffice/internal/app/repository"
	"emoney-backoffice/pkg/cache"
)

// Option anything any service object needed
type Option struct {
	commons.Options
	*repository.Repository
	Cache cache.ICache
}

// Services all service object injected here
type Services struct {
	Health     IHealthService
	Auth       IAuthService
	AuditTrail IAuditTrailService
}
