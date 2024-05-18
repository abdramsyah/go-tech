package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// RegisterRequest adalah struct yang digunakan untuk menerima data registrasi konsumen
type RegisterRequest struct {
	// NIK          string  `json:"nik" tag:"nik"`
	FullName  string `json:"full_name" tag:"full_name"`
	LegalName string `json:"legal_name" tag:"legal_name"`
	// BirthPlace   string  `json:"birth_place" tag:"birth_place"`
	// BirthDate    string  `json:"birth_date" tag:"birth_date"`
	// Salary       float64 `json:"salary" tag:"salary"`
	// KTPPhoto     string  `json:"ktp_photo" tag:"ktp_photo"`
	// SelfiePhoto  string  `json:"selfie_photo" tag:"selfie_photo"`
	RoleID   uint   `json:"role_id" tag:"role_id"`
	Email    string `json:"email" tag:"email"`
	Password string `json:"password" tag:"password"`
}

// Validate adalah metode untuk memvalidasi RegisterRequest
func (r RegisterRequest) Validate() error {
	validation.ErrorTag = "tag"
	return validation.ValidateStruct(&r,
		// validation.Field(&r.NIK,
		// 	validation.Required.Error("NIK is required"),
		// 	validation.Length(16, 16).Error("NIK must be exactly 16 characters")),
		validation.Field(&r.FullName,
			validation.Required.Error("Full name is required"),
			validation.Length(1, 100).Error("Full name must be between 1 and 100 characters")),
		validation.Field(&r.LegalName,
			validation.Required.Error("Legal name is required"),
			validation.Length(1, 100).Error("Legal name must be between 1 and 100 characters")),
		// validation.Field(&r.BirthPlace,
		// 	validation.Required.Error("Birth place is required"),
		// 	validation.Length(1, 50).Error("Birth place must be between 1 and 50 characters")),
		// validation.Field(&r.BirthDate,
		// 	validation.Required.Error("Birth date is required"),
		// 	isDate("2006-01-02").Error("Birth date must be in YYYY-MM-DD format")),
		// validation.Field(&r.Salary,
		// 	validation.Required.Error("Salary is required"),
		// 	validation.Min(0.0).Error("Salary must be a positive number")),
		// validation.Field(&r.KTPPhoto,
		// 	validation.Required.Error("KTP photo is required")),
		// validation.Field(&r.SelfiePhoto,
		// 	validation.Required.Error("Selfie photo is required")),
		validation.Field(&r.RoleID,
			validation.Required.Error("Role ID is required")),
		validation.Field(&r.Email,
			validation.Required.Error("Email is required")),
	)
}

type RegisterResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" tag:"Email"`
	Password string `json:"password" tag:"Password"`
}

func (r LoginRequest) Validate() error {
	validation.ErrorTag = "tag"
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email,
			validation.Required,
			is.Email),
		validation.Field(&r.Password,
			validation.Required),
	)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" tag:"Refresh Token"`
}

func (r RefreshTokenRequest) Validate() error {
	validation.ErrorTag = "tag"
	return validation.ValidateStruct(&r,
		validation.Field(&r.RefreshToken,
			validation.Required),
	)
}

type JwtToken struct {
	AccessToken         string `json:"accessToken"`
	AccessTokenExpires  int64  `json:"accessTokenExpires"`
	RefreshToken        string `json:"refreshToken"`
	RefreshTokenExpires int64  `json:"refreshTokenExpires"`
}

type TokenValidationResult struct {
	UserID     uint
	AccessUUID string
	RoleType   string
}
