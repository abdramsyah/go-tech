package dto

import "time"

type UserProfileResponse struct {
	ID          uint       `json:"id"`
	NIK         string     `json:"nik"`
	FullName    string     `json:"full_name"`
	LegalName   string     `json:"legal_name"`
	BirthPlace  string     `json:"birth_place"`
	BirthDate   string     `json:"birth_date"`
	Salary      float64    `json:"salary"`
	KTPPhoto    string     `json:"ktp_photo"`
	SelfiePhoto string     `json:"selfie_photo"`
	Status      string     `json:"status"`
	Email       string     `json:"email"`
	RoleID      uint       `json:"role_id"`
	Role        RoleEmbed  `json:"role"`
	PhoneNumber string     `json:"phone_number"`
	CreatedBy   uint       `json:"created_by"`
	UpdatedBy   uint       `json:"updated_by"`
	DeletedBy   *uint      `json:"deleted_by,omitempty"` // Use pointer to omit if null
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"` // Use pointer to omit if null
}
