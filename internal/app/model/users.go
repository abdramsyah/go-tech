package model

import (
	"database/sql"

	"gorm.io/gorm"
)

// Users represents a consumer in the system
type User struct {
	gorm.Model
	NIK          string
	FullName     string
	LegalName    string
	BirthPlace   string
	BirthDate    string
	Salary       float64
	KTPPhoto     string
	SelfiePhoto  string
	Status       string
	Email        string
	RoleID       uint
	PasswordHash string
	CreatedBy    uint64
	UpdatedBy    uint64
	DeletedBy    sql.NullInt64
	UserCreate   *Admin `gorm:"foreignKey:CreatedBy"`
	UserUpdate   *Admin `gorm:"foreignKey:UpdatedBy"`
	UserDelete   *Admin `gorm:"foreignKey:DeletedBy"`
}

// TableName sets the insert table name for this struct type
func (c *User) TableName() string {
	return "users"
}
