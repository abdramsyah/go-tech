package repository

import (
	"go-tech/internal/app/commons"
	"go-tech/internal/app/model"
	"go-tech/internal/app/util"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Count() (count int64, err error)
	Find(pConfig commons.PaginationConfig) (units []model.User, err error)
	FindByNIK(nik string) (user model.User, err error)
	FindByEmail(ctx echo.Context, email string) (user model.User, err error)
	FindByID(ctx echo.Context, userID uint) (user model.User, err error)
	UpdateWithMap(ctx echo.Context, userID uint, user map[string]interface{}) (err error)
	Update(user model.User, conditions *model.User, tx *gorm.DB) (err error)
	Create(user *model.User, tx *gorm.DB) (err error)
}

type userRepository struct {
	opt Option
}

func NewUserRepository(opt Option) IUserRepository {
	return &userRepository{
		opt: opt,
	}
}

func (r *userRepository) Count() (count int64, err error) {
	err = r.opt.DB.Model(&model.User{}).Count(&count).Error
	return
}

func (r *userRepository) Find(pConfig commons.PaginationConfig) (users []model.User, err error) {
	err = r.opt.DB.Scopes(util.Paginate(pConfig)).Preload("Role").Preload("Unit").Preload("JobTitle").Preload("AdminCreate").Preload("AdminUpdate").Order("id DESC").Find(&users).Error
	return
}

func (r *userRepository) FindByNIK(nik string) (user model.User, err error) {
	err = r.opt.DB.Preload("Role").First(&user, "nik = ?", nik).Error
	return
}

func (r *userRepository) FindByEmail(ctx echo.Context, email string) (user model.User, err error) {
	err = r.opt.DB.Joins("Role").First(&user, "email = ?", email).Error
	return
}

func (r *userRepository) FindByID(ctx echo.Context, userID uint) (user model.User, err error) {
	err = r.opt.DB.
		Joins("Role").
		// Joins("UserCreate").
		// Joins("UserUpdate").
		// Joins("Store").
		First(&user, userID).Error
	return
}

func (r *userRepository) UpdateWithMap(ctx echo.Context, userID uint, user map[string]interface{}) (err error) {
	err = r.opt.DB.Model(&model.User{}).Where("id = ?", userID).Updates(user).Error
	return
}

func (r *userRepository) Update(user model.User, conditions *model.User, tx *gorm.DB) (err error) {
	if tx != nil {
		err = tx.Where(conditions).Updates(user).Error
	} else {
		err = r.opt.DB.Where(conditions).Updates(user).Error
	}
	return
}

func (r *userRepository) Create(user *model.User, tx *gorm.DB) (err error) {
	if tx != nil {
		err = tx.Create(user).Error
	} else {
		err = r.opt.DB.Create(user).Error
	}
	return
}
