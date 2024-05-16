package model

import (
	"database/sql"
	"errors"
	"go-tech/internal/app/constant"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string
	Description sql.NullString
	CreatedBy   uint64
	UpdatedBy   uint64
	DeletedBy   sql.NullInt64
	Admins      []Admin
	AdminCreate Admin `gorm:"foreignKey:ID;references:CreatedBy"`
	AdminUpdate Admin `gorm:"foreignKey:ID;references:UpdatedBy"`
	AdminDelete Admin `gorm:"foreignKey:ID;references:DeletedBy"`
}

func (m *Role) BeforeDelete(tx *gorm.DB) (err error) {
	result := tx.Where("role_id = ?", m.ID).Find(&Admin{})
	if result.RowsAffected > 0 {
		err = errors.New(constant.ErrDataRelatedToOtherData)
	}
	return
}

func (m *Role) AfterDelete(tx *gorm.DB) (err error) {
	err = tx.Model(m).Unscoped().Where("id = ?", m.ID).Update("deleted_by", m.DeletedBy.Int64).Error
	return
}
