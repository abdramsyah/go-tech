package dto

import "time"

type AuditTrailRequest struct {
	RouteName string
	AdminID   uint
	URL       string
}

type AuditTrailFilter struct {
	AdminEmail *string
	Action     *string
	StartDate  *time.Time
	EndDate    *time.Time
}

type BORoutesResponse struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Name   string `json:"name"`
}

type AuditTrailResponse struct {
	ID         int64  `json:"id"`
	AdminID    int64  `json:"admin_id"`
	AdminEmail string `json:"admin_email"`
	AdminName  string `json:"admin_name"`
	AdminRole  string `json:"admin_role"`
	Action     string `json:"action"`
	URL        string `json:"url"`
	CreatedAt  string `json:"created_at"`
	RequestID  string `json:"request_id"`
}
