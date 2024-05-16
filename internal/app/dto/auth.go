package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

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
