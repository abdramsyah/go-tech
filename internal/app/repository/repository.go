package repository

import (
	"emoney-backoffice/internal/app/commons"
)

// Option anything any repo object needed
type Option struct {
	commons.Options
}

// Repository all repo object injected here
type Repository struct {
	Admin      IAdminRepository
	AuditTrail IAuditTrailRepository
}
