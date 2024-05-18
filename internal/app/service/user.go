package service

import (
	"errors"
	"go-tech/internal/app/constant"
	"go-tech/internal/app/model"
	"go-tech/internal/app/util"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IUserService interface {
	// Find(ctx echo.Context, pConfig commons.PaginationConfig, filter *dto.UserFilter) (users []model.User, count int64, err error)
	Profile(ctx echo.Context, userID uint) (user model.User, err error)
	ChangePassword(ctx echo.Context, userID uint, oldPassword string, newPassword string) (err error)
	// Update(ctx echo.Context, userID uint, user model.User) (err error)
	// UpdateRole(ctx echo.Context, userID uint, req *dto.SetUserRoleRequest) (err error)
	// Create(ctx echo.Context, req *dto.CreateUserRequest) (err error)
	// GetPermissions(ctx echo.Context, userID uint) (permissions map[string]interface{}, err error)
	FindUserByID(ctx echo.Context, ID uint) (user model.User, err error)
	// FindAll(ctx echo.Context, filter *dto.UserFilter) (users []model.User, err error)
	// UpdateProfile(ctx echo.Context, userID uint, req *dto.UpdateProfileRequest) (err error)
}

type userService struct {
	opt Option
}

func NewUserService(opt Option) IUserService {
	return &userService{
		opt: opt,
	}
}

// func (s *userService) Find(ctx echo.Context, pConfig commons.PaginationConfig, filter *dto.UserFilter) (users []model.User, count int64, err error) {
// 	var waitGroup sync.WaitGroup
// 	c := make(chan error)

// 	waitGroup.Add(2)

// 	go func() {
// 		waitGroup.Wait()
// 		close(c)
// 	}()

// 	go func() {
// 		defer waitGroup.Done()

// 		count, err = s.opt.Repository.User.Count(ctx, filter)
// 		if err != nil {
// 			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Get user count",
// 				zap.Error(err),
// 			)
// 			err = util.ErrFailedGetDataCount()
// 			c <- err
// 		}
// 	}()

// 	go func() {
// 		defer waitGroup.Done()

// 		users, err = s.opt.Repository.User.Find(ctx, pConfig, filter)
// 		if err != nil {
// 			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Get users",
// 				zap.Error(err),
// 			)
// 			err = util.ErrFailedFetchData()
// 			c <- err
// 		}
// 	}()

// 	for errChan := range c {
// 		if errChan != nil {
// 			err = errChan
// 			return
// 		}
// 	}

// 	return
// }

func (s *userService) Profile(ctx echo.Context, userID uint) (user model.User, err error) {
	user, err = s.opt.Repository.User.FindByID(ctx, userID)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to get profile", zap.Error(err),
			zap.Uint("User ID", userID))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = util.ErrDataNotFound()
			return
		}
		err = util.ErrInternalServerError()
		return
	}
	return
}

func (s *userService) ChangePassword(ctx echo.Context, userID uint, oldPassword string, newPassword string) (err error) {
	user, err := s.Profile(ctx, userID)
	if err != nil {
		return
	}

	check := util.CheckPasswordHash(oldPassword, user.PasswordHash)
	if !check {
		err = util.ErrRequestValidation("Password lama tidak valid")
		return
	}

	isNewPasswordValid := util.PasswordValidator2(newPassword, constant.UserMinPasswordLength)
	if !isNewPasswordValid {
		err = util.ErrRequestValidation("Format password baru tidak sesuai")
		return
	}

	hashPassword, err := util.HashPassword(newPassword)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to hash password",
			zap.Error(err),
			zap.Uint("User ID", userID))
		err = util.ErrUnknownError("Ubah password gagal, silahkan coba lagi")
		return
	}

	dataUpdate := map[string]interface{}{
		"password_hash": hashPassword,
	}
	err = s.opt.Repository.User.UpdateWithMap(ctx, userID, dataUpdate)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to update user", zap.Error(err), zap.Uint("User ID", userID))
		err = util.ErrInternalServerError()
		return
	}
	return
}

// func (s *userService) Update(ctx echo.Context, userID uint, user model.User) (err error) {
// 	userExisting, err := s.opt.Repository.User.FindByID(ctx, userID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			err = util.ErrDataNotFound()
// 		} else {
// 			err = util.ErrInternalServerError()
// 		}
// 		return
// 	}

// 	if user.Email != "" {
// 		userCheck, errCheck := s.opt.Repository.User.FindByEmail(ctx, user.Email)
// 		if errCheck != nil {
// 			if !errors.Is(errCheck, gorm.ErrRecordNotFound) {
// 				err = util.ErrInternalServerError()
// 				return
// 			}
// 		} else {
// 			if userCheck.ID != userID {
// 				err = util.ErrRequestValidation("Email sudah digunakan oleh pengguna lain")
// 				return
// 			}
// 		}
// 	}

// 	if user.EmployeeID != "" {
// 		userCheck, errCheck := s.opt.Repository.User.FindByEmployeeID(ctx, user.EmployeeID)
// 		if errCheck != nil {
// 			if !errors.Is(errCheck, gorm.ErrRecordNotFound) {
// 				err = util.ErrInternalServerError()
// 				return
// 			}
// 		} else {
// 			if userCheck.ID != userID {
// 				err = util.ErrRequestValidation("ID pegawai sudah digunakan oleh pengguna lain")
// 				return
// 			}
// 		}
// 	}
// 	role, err := s.opt.Repository.Role.FindRoleByID(ctx, user.RoleID)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error get role by id",
// 			zap.Error(err),
// 		)
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			err = util.ErrRequestValidation("Data role tidak ditemukan")
// 		} else {
// 			err = util.ErrInternalServerError()
// 		}
// 		return
// 	}

// 	if role.RoleType == constant.RoleTypeStore {
// 		_, err = s.opt.Repository.Store.FindStoreByID(ctx, uint(user.StoreID.Int64))
// 		if err != nil {
// 			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error update user",
// 				zap.Error(err),
// 			)
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				err = util.ErrRequestValidation("Toko tidak ditemukan")
// 			} else {
// 				err = util.ErrInternalServerError()
// 			}
// 			return
// 		}
// 	}

// 	if userExisting.RoleID != user.RoleID {
// 		_, err = s.opt.Options.Rbac.DeleteRolesForUser(util.FormatRbacSubject(user.ID))
// 		if err != nil {
// 			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to delele user role", zap.Error(err),
// 				zap.Uint("User ID", userID))
// 			err = util.ErrInternalServerError()
// 			return
// 		}

// 		_, err = s.opt.Options.Rbac.AddRoleForUser(util.FormatRbacSubject(userID), util.FormatRbacRole(user.RoleID))
// 		if err != nil {
// 			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to add user role", zap.Error(err))
// 			err = util.ErrInternalServerError()
// 			return
// 		}

// 		dataUpdate := map[string]interface{}{
// 			"role_id": user.RoleID,
// 		}
// 		err = s.opt.Repository.User.UpdateWithMap(ctx, userID, dataUpdate)
// 		if err != nil {
// 			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to update user", zap.Error(err), zap.Uint("User ID", userID))
// 			err = util.ErrInternalServerError()
// 			return
// 		}
// 	}

// 	//Auto logout user when status change from active to inactive
// 	if userExisting.Status == constant.UserStatusActive && user.Status == constant.UserStatusInactive {
// 		err = s.autoLogout(userExisting.ID)
// 		if err != nil {
// 			return
// 		}
// 	}

// 	dataUpdate := map[string]interface{}{
// 		"full_name":    user.FullName,
// 		"email":        user.Email,
// 		"status":       user.Status,
// 		"updated_by":   user.UpdatedBy,
// 		"phone_number": user.PhoneNumber,
// 		"employee_id":  user.EmployeeID,
// 		"role_id":      user.RoleID,
// 	}
// 	if user.StoreID.Int64 != 0 {
// 		dataUpdate["store_id"] = user.StoreID
// 	} else {
// 		dataUpdate["store_id"] = nil
// 	}
// 	err = s.opt.Repository.User.UpdateWithMap(ctx, userID, dataUpdate)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to update user", zap.Error(err), zap.Uint("User ID", userID))
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	return
// }

// func (s *userService) UpdateRole(ctx echo.Context, userID uint, req *dto.SetUserRoleRequest) (err error) {
// 	actx, err := util.NewAppContext(ctx)
// 	if err != nil {
// 		return
// 	}

// 	user, err := s.opt.Repository.User.FindByID(actx, userID)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error get user by id",
// 			zap.Error(err),
// 		)
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			err = util.ErrRequestValidation("Data pengguna tidak ditemukan")
// 		} else {
// 			err = util.ErrInternalServerError()
// 		}
// 		return
// 	}

// 	_, err = s.opt.Repository.Role.FindRoleByID(actx, req.RoleID)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error get role by id",
// 			zap.Error(err),
// 		)
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			err = util.ErrRequestValidation("Data role tidak ditemukan")
// 		} else {
// 			err = util.ErrInternalServerError()
// 		}
// 		return
// 	}

// 	_, err = s.opt.Options.Rbac.DeleteRolesForUser(util.FormatRbacSubject(user.ID))
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Failed to delele user role", zap.Error(err),
// 			zap.Uint("User ID", userID))
// 		err = util.ErrInternalServerError()
// 		return
// 	}

// 	_, err = s.opt.Options.Rbac.AddRoleForUser(util.FormatRbacSubject(userID), util.FormatRbacRole(req.RoleID))
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Failed to add user role", zap.Error(err))
// 		err = util.ErrInternalServerError()
// 		return
// 	}

// 	dataUpdate := map[string]interface{}{
// 		"role_id": req.RoleID,
// 	}
// 	err = s.opt.Repository.User.UpdateWithMap(ctx, userID, dataUpdate)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to update user", zap.Error(err), zap.Uint("User ID", userID))
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	return
// }

// func (s *userService) Create(ctx echo.Context, req *dto.CreateUserRequest) (err error) {
// 	actx, err := util.NewAppContext(ctx)
// 	if err != nil {
// 		return
// 	}
// 	_, err = s.opt.Repository.User.FindByEmail(actx, req.Email)
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Warn("Error get user",
// 			zap.String("Email", req.Email),
// 			zap.Error(err))
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	if err == nil {
// 		err = util.ErrRequestValidation("Email sudah digunakan oleh pengguna lain")
// 		return
// 	}

// 	_, err = s.opt.Repository.User.FindByEmployeeID(actx, req.EmployeeID)
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Warn("Error get user",
// 			zap.String("Employee ID", req.EmployeeID),
// 			zap.Error(err))
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	if err == nil {
// 		err = util.ErrRequestValidation("ID pegawai sudah digunakan oleh pengguna lain")
// 		return
// 	}

// 	_, err = s.opt.Repository.User.FindByUsername(actx, req.Username)
// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Warn("Error get user",
// 			zap.String("Username", req.Username),
// 			zap.Error(err))
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	if err == nil {
// 		err = util.ErrRequestValidation("Username sudah digunakan oleh pengguna lain")
// 		return
// 	}

// 	role, err := s.opt.Repository.Role.FindRoleByID(actx, req.RoleID)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error create user",
// 			zap.Error(err),
// 		)
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			err = util.ErrRequestValidation("Role tidak ditemukan")
// 		} else {
// 			err = util.ErrInternalServerError()
// 		}
// 		return
// 	}

// 	if role.RoleType == constant.RoleTypeStore {
// 		if req.StoreID != 0 {
// 			_, err = s.opt.Repository.Store.FindStoreByID(actx, uint(req.StoreID))
// 			if err != nil {
// 				s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error create user",
// 					zap.Error(err),
// 				)
// 				if errors.Is(err, gorm.ErrRecordNotFound) {
// 					err = util.ErrRequestValidation("Toko tidak ditemukan")
// 				} else {
// 					err = util.ErrInternalServerError()
// 				}
// 				return
// 			}
// 		} else {
// 			err = util.ErrRequestValidation("ID toko tidak boleh kosong")
// 			return
// 		}
// 	}

// 	userID := actx.GetUserID()
// 	user := &model.User{
// 		FullName:     req.FullName,
// 		Email:        req.Email,
// 		RoleID:       req.RoleID,
// 		Status:       constant.UserStatusActive,
// 		CreatedBy:    userID,
// 		UpdatedBy:    userID,
// 		Username:     req.Username,
// 		PhoneNumber:  req.PhoneNumber,
// 		PasswordHash: req.Password,
// 		EmployeeID:   req.EmployeeID,
// 	}
// 	if req.StoreID != 0 {
// 		user.StoreID = sql.NullInt64{
// 			Int64: int64(req.StoreID),
// 			Valid: true,
// 		}
// 	}
// 	tx := s.opt.DB.Begin()
// 	err = s.createUser(actx, user, tx)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error create user",
// 			zap.Error(err),
// 		)
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	tx.Commit()
// 	return
// }

// func (s *userService) createUser(actx echo.Context, user *model.User, tx *gorm.DB) (err error) {
// 	password := user.PasswordHash
// 	passwordHash, err := util.HashPassword(password)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error create password",
// 			zap.Error(err),
// 		)
// 		err = util.ErrUnknownError("Gagal untuk melakukan encrypt password")
// 		return
// 	}

// 	user.PasswordHash = passwordHash
// 	err = s.opt.Repository.User.Create(actx, user, tx)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error create user",
// 			zap.Error(err),
// 		)
// 		err = util.ErrUnknownError("Gagal untuk menambahkan pengguna")
// 		return
// 	}

// 	_, err = s.opt.Options.Rbac.AddRoleForUser(util.FormatRbacSubject(user.ID), util.FormatRbacRole(user.RoleID))
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Failed to set role",
// 			zap.Error(err),
// 		)
// 		err = util.ErrUnknownError("Gagal untuk set role")
// 		tx.Rollback()
// 		return
// 	}

// 	var (
// 		replacer        []string
// 		placeholderName = map[string]interface{}{
// 			"fullname":  user.FullName,
// 			"username":  user.Username,
// 			"password":  password,
// 			"login_url": s.opt.Config.LoginURL,
// 		}
// 	)
// 	for k, p := range placeholderName {
// 		replacer = append(replacer, fmt.Sprintf("[%s]", k))
// 		replacer = append(replacer, cast.ToString(p))
// 	}
// 	r := strings.NewReplacer(replacer...)
// 	emailMessage := constant.EmailUserCreated
// 	emailMessage = r.Replace(emailMessage)
// 	err = s.opt.EMailService.Send(actx.Request().Context(), email.EmailRequest{
// 		EmailTo: []string{user.Email},
// 		Subject: "User Credential",
// 		Message: emailMessage,
// 	})
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Failed to send user credential",
// 			zap.Error(err),
// 		)
// 		err = util.ErrUnknownError("Gagal mengirimkan kredensial pengguna")
// 		tx.Rollback()
// 		return
// 	}

// 	return
// }

// func (s *userService) GetPermissions(ctx echo.Context, userID uint) (permissions map[string]interface{}, err error) {
// 	subject := util.FormatRbacSubject(userID)
// 	permissionsB, err := casbin.CasbinJsGetPermissionForUserOld(s.opt.Rbac, subject)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error get user's permissions",
// 			zap.Error(err),
// 		)
// 		err = util.ErrUnknownError("Gagal mendapatkan hak akses pengguna")
// 		return
// 	}
// 	err = json.Unmarshal(permissionsB, &permissions)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error unmarshal permissions data",
// 			zap.Error(err),
// 		)
// 		err = util.ErrUnknownError("Gagal mendapatkan hak akses pengguna")
// 	}
// 	return
// }

func (s *userService) FindUserByID(ctx echo.Context, ID uint) (user model.User, err error) {
	user, err = s.opt.Repository.User.FindByID(ctx, ID)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error users by ID",
			zap.Uint("User ID", ID),
			zap.Error(err),
		)
		err = util.ErrUnknownError("Gagal menemukan pengguna melalui ID")
	}
	return
}

// func (s *userService) autoLogout(userID uint) (err error) {
// 	uuidAccessKey := fmt.Sprintf("%s:%s:%d", util.CacheKeyFormatter("uuid"), constant.TokenAccessType, userID)
// 	uuidAccessCache, err := s.opt.Cache.ReadCache(uuidAccessKey)
// 	if err != nil {
// 		s.opt.Logger.Error("error get uuid access token from cache",
// 			zap.Uint("user id", userID),
// 			zap.Error(err),
// 		)
// 		if !strings.Contains(err.Error(), "Cache key didn't exists") {
// 			err = util.ErrUnknownError("Gagal logout otomatis")
// 			return
// 		}
// 		//if err is cache key didn't exists, ignore error
// 		err = nil
// 	} else {
// 		uuidAccessCacheValue := new(commons.AuthUUIDCacheValue)
// 		err = json.Unmarshal(uuidAccessCache, uuidAccessCacheValue)
// 		if err != nil {
// 			s.opt.Logger.Error("error unmarshal uuid access token from cache",
// 				zap.Uint("user id", userID),
// 				zap.Error(err),
// 			)
// 			err = util.ErrUnknownError("Gagal logout otomatis")
// 			return
// 		}
// 		// delete access token
// 		err = s.opt.Cache.DeleteCache(fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), uuidAccessCacheValue.UUID))
// 		if err != nil {
// 			s.opt.Logger.Error(util.ErrLogoutDefault().Error(),
// 				zap.Uint("user id", userID),
// 				zap.Error(err))
// 			err = util.ErrUnknownError("Gagal logout otomatis")
// 			return
// 		}
// 	}

// 	uuidRefreshKey := fmt.Sprintf("%s:%s:%d", util.CacheKeyFormatter("uuid"), constant.TokenRefreshType, userID)
// 	uuidRefreshCache, err := s.opt.Cache.ReadCache(uuidRefreshKey)
// 	if err != nil {
// 		s.opt.Logger.Error("error get uuid refresh token from cache",
// 			zap.Error(err),
// 		)
// 		if !strings.Contains(err.Error(), "Cache key didn't exists") {
// 			err = util.ErrUnknownError("Gagal logout otomatis")
// 			return
// 		}
// 		//if err is cache key didn't exists, ignore error
// 		err = nil
// 	} else {
// 		uuidRefreshCacheValue := new(commons.AuthUUIDCacheValue)
// 		err = json.Unmarshal(uuidRefreshCache, uuidRefreshCacheValue)
// 		if err != nil {
// 			s.opt.Logger.Error("error unmarshal uuid refresh token from cache",
// 				zap.Uint("user id", userID),
// 				zap.Error(err),
// 			)
// 			err = util.ErrUnknownError("Gagal logout otomatis")
// 			return
// 		}

// 		// delete refresh token
// 		err = s.opt.Cache.DeleteCache(fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), uuidRefreshCacheValue.UUID))
// 		if err != nil {
// 			s.opt.Logger.Error(util.ErrLogoutDefault().Error(),
// 				zap.Uint("user id", userID),
// 				zap.Error(err))
// 			err = util.ErrUnknownError("Gagal logout otomatis")
// 			return
// 		}
// 	}
// 	return
// }

// func (s *userService) UpdateProfile(ctx echo.Context, userID uint, req *dto.UpdateProfileRequest) (err error) {
// 	_, err = s.opt.Repository.User.FindByID(ctx, userID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			err = util.ErrDataNotFound()
// 		} else {
// 			err = util.ErrInternalServerError()
// 		}
// 		return
// 	}

// 	if req.Email != "" {
// 		userCheck, errCheck := s.opt.Repository.User.FindByEmail(ctx, req.Email)
// 		if errCheck != nil {
// 			if !errors.Is(errCheck, gorm.ErrRecordNotFound) {
// 				err = util.ErrInternalServerError()
// 				return
// 			}
// 		} else {
// 			if userCheck.ID != userID {
// 				err = util.ErrRequestValidation("Email sudah digunakan oleh pengguna lain")
// 				return
// 			}
// 		}
// 	}

// 	dataUpdate := map[string]interface{}{
// 		"full_name":    req.FullName,
// 		"email":        req.Email,
// 		"phone_number": req.PhoneNumber,
// 	}
// 	err = s.opt.Repository.User.UpdateWithMap(ctx, userID, dataUpdate)
// 	if err != nil {
// 		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to update user", zap.Error(err), zap.Uint("User ID", userID))
// 		err = util.ErrInternalServerError()
// 		return
// 	}
// 	return
// }
