package handler

import (
	"go-tech/internal/app/constant"
	"go-tech/internal/app/dto"
	"go-tech/internal/app/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	HandlerOption
}

// func (h UserHandler) GetPermission(c echo.Context) (resp dto.HttpResponse) {
// 	actx, err := util.NewAppContext(c)
// 	if err != nil {
// 		return
// 	}
// 	userID := actx.GetUserID()
// 	res, err := h.Services.User.GetPermissions(actx, userID)
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(err, nil)
// 		return
// 	}

// 	resp = dto.SuccessHttpResponse(http.StatusOK, "", "Berhasil mendapatkan akses pengguna", res)
// 	return
// }

func (h UserHandler) Profile(c echo.Context) (resp dto.HttpResponse) {
	actx, err := util.NewAppContext(c)
	if err != nil {
		return
	}
	userID := actx.GetUserID()
	var res dto.UserProfileResponse
	res, err = h.findByID(actx, userID)
	if err != nil {
		resp = dto.FailedHttpResponse(err, nil)
		return
	}

	resp = dto.SuccessHttpResponse(http.StatusOK, "", "Berhasil mendapatkan profil", res)
	return
}

func (h UserHandler) findByID(ctx echo.Context, userID uint) (resp dto.UserProfileResponse, err error) {
	user, err := h.Services.User.Profile(ctx, userID)
	if err != nil {
		return
	}

	resp = dto.UserProfileResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		NIK:      user.NIK,
		Status:   user.Status,
		Role: dto.RoleEmbed{
			ID:       user.Role.ID,
			Name:     user.Role.Name,
			RoleType: user.Role.RoleType,
		},
		LegalName:   user.LegalName,
		PhoneNumber: user.PhoneNumber,
		// CreatedAt:   user.CreatedAt.Format(constant.FeDatetimeFormat),
		CreatedAt: user.CreatedAt.Format(constant.FeDatetimeFormat),
		UpdatedAt: user.UpdatedAt.Format(constant.FeDatetimeFormat),
	}
	// if user.Role.RoleType == constant.RoleTypeStore && user.Store != nil {
	// 	resp.Store = dto.StoreEmbed{
	// 		ID:        user.Store.ID,
	// 		Name:      user.Store.Name,
	// 		StoreCode: user.Store.StoreCode,
	// 	}
	// }

	return
}

// func (h UserHandler) Create(c echo.Context) (resp dto.HttpResponse) {
// 	var err error
// 	req := new(dto.CreateUserRequest)
// 	if err = c.Bind(req); err != nil {
// 		h.HandlerOption.Options.Logger.Error("Error bind request",
// 			zap.Error(err),
// 		)
// 		resp = dto.FailedHttpResponse(util.ErrBindRequest(), nil)
// 		return
// 	}

// 	err = req.Validate()
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(util.ErrRequestValidation(err.Error()), nil)
// 		return
// 	}

// 	err = h.Services.User.Create(c, req)
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(err, nil)
// 		return
// 	}

// 	resp = dto.SuccessHttpResponse(http.StatusCreated, "", "Berhasil membuat pengguna", nil)
// 	return
// }

// func (h UserHandler) Update(c echo.Context) (resp dto.HttpResponse) {
// 	actx, err := util.NewAppContext(c)
// 	if err != nil {
// 		return
// 	}
// 	req := new(dto.UpdateUserRequest)
// 	if err = actx.Bind(req); err != nil {
// 		h.HandlerOption.Options.Logger.Error("Error bind request",
// 			zap.Error(err),
// 		)
// 		resp = dto.FailedHttpResponse(util.ErrBindRequest(), nil)
// 		return
// 	}

// 	err = req.Validate()
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(util.ErrRequestValidation(err.Error()), nil)
// 		return
// 	}

// 	ID := c.Param("ID")
// 	if ID == "" {
// 		resp = dto.FailedHttpResponse(util.ErrRequestValidation("ID pengguna tidak valid"), nil)
// 		return
// 	}

// 	userID := actx.GetUserID()
// 	user := model.User{
// 		FullName:    req.FullName,
// 		Email:       req.Email,
// 		UpdatedBy:   userID,
// 		Status:      req.Status,
// 		EmployeeID:  req.EmployeeID,
// 		PhoneNumber: req.PhoneNumber,
// 		RoleID:      req.RoleID,
// 		StoreID: sql.NullInt64{
// 			Valid: true,
// 			Int64: int64(req.StoreID),
// 		},
// 	}
// 	err = h.Services.User.Update(c, cast.ToUint(ID), user)
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(err, nil)
// 		return
// 	}

// 	resp = dto.SuccessHttpResponse(http.StatusOK, "", "Berhasil memperbaharui pengguna", nil)
// 	return
// }

// func (h UserHandler) ChangePassword(c echo.Context) (resp dto.HttpResponse) {
// 	actx, err := util.NewAppContext(c)
// 	if err != nil {
// 		return
// 	}
// 	userID := actx.GetUserID()
// 	req := new(dto.ChangePasswordRequest)
// 	if err = actx.Bind(req); err != nil {
// 		h.HandlerOption.Options.Logger.Error("Error bind ChangePasswordRequest request",
// 			zap.Error(err),
// 		)
// 		resp = dto.FailedHttpResponse(util.ErrBindRequest(), nil)
// 		return
// 	}

// 	err = req.Validate()
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(util.ErrRequestValidation(err.Error()), nil)
// 		return
// 	}

// 	if req.NewPassword != req.ConfirmPassword {
// 		resp = dto.FailedHttpResponse(util.ErrRequestValidation("Kata sandi baru tidak cocok dengan konfirmasi kata sandi"), nil)
// 		return
// 	}

// 	err = h.Services.User.ChangePassword(actx, userID, req.OldPassword, req.NewPassword)
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(err, nil)
// 		return
// 	}

// 	resp = dto.SuccessHttpResponse(http.StatusOK, "", "Ubah kata sandi berhasil", nil)
// 	return
// }

// func (h UserHandler) UpdateProfile(c echo.Context) (resp dto.HttpResponse) {
// 	actx, err := util.NewAppContext(c)
// 	if err != nil {
// 		return
// 	}
// 	userID := actx.GetUserID()
// 	req := new(dto.UpdateProfileRequest)
// 	if err = actx.Bind(req); err != nil {
// 		h.HandlerOption.Options.Logger.Error("Error bind ChangePasswordRequest request",
// 			zap.Error(err),
// 		)
// 		resp = dto.FailedHttpResponse(util.ErrBindRequest(), nil)
// 		return
// 	}

// 	err = req.Validate()
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(util.ErrRequestValidation(err.Error()), nil)
// 		return
// 	}

// 	err = h.Services.User.UpdateProfile(actx, userID, req)
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(err, nil)
// 		return
// 	}

// 	res, err := h.findByID(actx, userID)
// 	if err != nil {
// 		resp = dto.FailedHttpResponse(err, nil)
// 		return
// 	}

// 	resp = dto.SuccessHttpResponse(http.StatusOK, "", "Ubah profile pengguna berhasil", res)
// 	return
// }
