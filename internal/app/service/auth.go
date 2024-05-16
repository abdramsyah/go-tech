package service

import (
	"emoney-backoffice/internal/app/commons"
	"emoney-backoffice/internal/app/constant"
	"emoney-backoffice/internal/app/dto"
	"emoney-backoffice/internal/app/model"
	"emoney-backoffice/internal/app/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type IAuthService interface {
	Login(ctx echo.Context, req *dto.LoginRequest) (httpStatus int, jwtToken dto.JwtToken, err error)
	Logout(ctx echo.Context, adminID uint64, accessUUID string) (httpStatus int, err error)
	RefreshToken(ctx echo.Context, refreshToken string) (httpStatus int, jwtToken dto.JwtToken, err error)
	ValidateToken(ctx echo.Context, r *http.Request) (claims jwt.MapClaims, err error)
	PermissionCheck(ctx echo.Context, object string, action string) (isPermitted bool, err error)
	BatchPermissionCheck(ctx echo.Context, request [][]interface{}) (isPermitted bool, err error)
}

type authService struct {
	opt Option
}

func NewAuthService(opt Option) IAuthService {
	return &authService{
		opt: opt,
	}
}

func (s *authService) createJwtToken(adminID uint64) (td *commons.TokenDetails, err error) {
	td = &commons.TokenDetails{}
	accessExpired := cast.ToDuration(s.opt.Config.JwtAccessTtl)
	refreshExpired := cast.ToDuration(s.opt.Config.JwtRefreshTtl)
	td.AtExpires = time.Now().Add(accessExpired).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(refreshExpired).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + strconv.Itoa(int(adminID))

	accessSecret := s.opt.Config.JwtAccessSecret
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["admin_id"] = adminID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return
	}
	err = s.storeToRedis(td.AccessUuid, adminID, accessExpired)
	if err != nil {
		return
	}

	refreshSecret := s.opt.Config.JwtRefreshSecret
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["admin_id"] = adminID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return
	}

	err = s.storeToRedis(td.RefreshUuid, adminID, refreshExpired)

	return
}

func (s *authService) storeToRedis(uuid string, adminID uint64, duration time.Duration) (err error) {
	val := commons.AuthCacheValue{
		AdminID: adminID,
	}
	b, err := json.Marshal(val)
	if err != nil {
		s.opt.Logger.Error("Failed to marshal data", zap.Error(err))
		return
	}
	err = s.opt.Cache.WriteCache(fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), uuid), b, duration)
	return
}

func (s *authService) Login(ctx echo.Context, req *dto.LoginRequest) (httpStatus int, jwtToken dto.JwtToken, err error) {
	admin, err := s.opt.Repository.Admin.FindByEmail(req.Email)
	if err != nil {
		s.opt.Logger.Warn(constant.ErrLoginDefault,
			zap.String("Email", req.Email),
			zap.Error(err))
		httpStatus = http.StatusNotFound
		err = errors.New(constant.ErrLoginDefault)
		return
	}

	if admin.Status != constant.AdminStatusActive {
		httpStatus = http.StatusUnauthorized
		err = errors.New("Admin status inactive")
		return
	}

	check := util.CheckPasswordHash(req.Password, admin.PasswordHash)
	if !check {
		httpStatus = http.StatusUnauthorized
		err = errors.New(constant.ErrLoginDefault)
		return
	}

	tokenDetail, err := s.createJwtToken(uint64(admin.ID))
	if err != nil {
		s.opt.Logger.Error("Failed to generate access token",
			zap.String("Email", req.Email),
			zap.Error(err))
		httpStatus = http.StatusInternalServerError
		err = errors.New("Failed to generate access token")
		return
	}

	auditTrail := model.AuditTrails{
		AdminID:    int64(admin.ID),
		AdminEmail: admin.Email,
		AdminName:  admin.FullName,
		AdminRole:  admin.Role.Name,
		Action:     "Login",
		URL:        "[POST] /login",
		CreatedAt:  time.Now(),
		RequestID:  null.NewString(util.GetRequestID(ctx), true),
	}

	err = s.opt.Repository.AuditTrail.Create(ctx, &auditTrail)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Failed to save audit trail",
			zap.String("Email", req.Email),
			zap.Error(err))
		httpStatus = http.StatusInternalServerError
		err = errors.New("Failed to generate access token")
		return
	}

	httpStatus = http.StatusOK
	jwtToken = dto.JwtToken{
		AccessToken:         tokenDetail.AccessToken,
		AccessTokenExpires:  tokenDetail.AtExpires,
		RefreshToken:        tokenDetail.RefreshToken,
		RefreshTokenExpires: tokenDetail.RtExpires,
	}

	return
}

func (s *authService) Logout(ctx echo.Context, adminID uint64, accessUUID string) (httpStatus int, err error) {
	refreshUuid := fmt.Sprintf("%s++%d", accessUUID, adminID)
	// delete access token
	err = s.opt.Cache.DeleteCache(fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), accessUUID))
	if err != nil {
		s.opt.Logger.Warn(constant.ErrLogoutDefault,
			zap.Uint64("Admin ID", adminID),
			zap.Error(err))
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New(constant.ErrLogoutDefault)
		return
	}
	// delete refresh token
	err = s.opt.Cache.DeleteCache(fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), refreshUuid))
	if err != nil {
		s.opt.Logger.Warn(constant.ErrLogoutDefault,
			zap.Uint64("Admin ID", adminID),
			zap.Error(err))
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New(constant.ErrLogoutDefault)
		return
	}

	httpStatus = http.StatusOK
	return
}

func (s *authService) verifyToken(tokenString string, tokenType string) (tokenStruct *jwt.Token, err error) {
	var secret string
	switch tokenType {
	case constant.TokenAccessType:
		secret = s.opt.Config.JwtAccessSecret
	case constant.TokenRefreshType:
		secret = s.opt.Config.JwtRefreshSecret
	default:
		err = errors.New("Unknown token type")
		return
	}

	tokenStruct, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method : %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	return
}

func (s *authService) RefreshToken(ctx echo.Context, refreshToken string) (httpStatus int, jwtToken dto.JwtToken, err error) {
	tokenStruct, err := s.verifyToken(refreshToken, constant.TokenRefreshType)
	if err != nil {
		s.opt.Logger.Error("Error Verify Token",
			zap.Error(err),
		)
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New("Invalid refresh token")
		return
	}

	claims, ok := tokenStruct.Claims.(jwt.MapClaims)
	if !ok || !tokenStruct.Valid {
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New("Invalid refresh token")
	}

	refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
	if !ok {
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New("Invalid refresh token claim")
		return
	}

	adminID, errConv := strconv.ParseUint(fmt.Sprintf("%.f", claims["admin_id"]), 10, 64)
	if err != nil {
		s.opt.Logger.Error("Invalid Admin ID",
			zap.Error(errConv),
		)
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New("Invalid refresh token")
		return
	}

	cacheKey := fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), refreshUuid)
	isRefreshTokenExist := s.opt.Cache.CheckCacheExists(cacheKey)
	if !isRefreshTokenExist {
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New("You already logged out, please re login")
		return
	}

	err = s.opt.Cache.DeleteCache(cacheKey)
	if err != nil {
		s.opt.Logger.Error("Failed to delete refresh token",
			zap.Error(err),
		)
		httpStatus = http.StatusUnprocessableEntity
		err = errors.New("Failed to delete refresh token")
		return
	}

	tokenDetail, errCrt := s.createJwtToken(adminID)
	if errCrt != nil {
		s.opt.Logger.Error("Failed to generate refresh token",
			zap.Error(errCrt))
		httpStatus = http.StatusInternalServerError
		err = errors.New("Failed to generate refresh token")
		return
	}

	httpStatus = http.StatusOK
	jwtToken = dto.JwtToken{
		AccessToken:         tokenDetail.AccessToken,
		AccessTokenExpires:  tokenDetail.AtExpires,
		RefreshToken:        tokenDetail.RefreshToken,
		RefreshTokenExpires: tokenDetail.RtExpires,
	}
	return
}

func (s *authService) ValidateToken(ctx echo.Context, r *http.Request) (claims jwt.MapClaims, err error) {
	tokenString, err := util.ExtractBearerToken(r.Header)
	if err != nil {
		s.opt.Logger.Error("Error Extract Token",
			zap.Error(err),
		)
		return
	}
	tokenStruct, err := s.verifyToken(tokenString, constant.TokenAccessType)
	if err != nil {
		s.opt.Logger.Error("Error Verify Token",
			zap.Error(err),
		)
		return
	}
	claims, ok := tokenStruct.Claims.(jwt.MapClaims)
	if !ok || !tokenStruct.Valid {
		s.opt.Logger.Error("Token claims is not valid")
		err = errors.New("Token claims is not valid")
		return
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		s.opt.Logger.Error("Access UUID is not valid")
		err = errors.New("Access UUID is not valid")
		return
	}
	adminID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["admin_id"]), 10, 64)
	if err != nil {
		s.opt.Logger.Error("Error Parse Admin ID From Token",
			zap.Error(err),
		)
		return
	}
	cacheKey := fmt.Sprintf("%s:%s", util.CacheKeyFormatter("jwt"), accessUuid)
	authCache, err := s.opt.Cache.ReadCache(cacheKey)
	if err != nil {
		s.opt.Logger.Error("Error Get Admin ID From Cache",
			zap.Error(err),
		)
		return
	}
	authCacheValue := new(commons.AuthCacheValue)
	err = json.Unmarshal(authCache, authCacheValue)
	if adminID != authCacheValue.AdminID {
		s.opt.Logger.Error("Admin ID from token different with Admin ID from cache")
		err = errors.New("Admin ID from token different with Admin ID from cache")
		return
	}
	return
}

func (s *authService) PermissionCheck(ctx echo.Context, object string, action string) (isPermitted bool, err error) {
	actx, err := util.NewAppContext(ctx)
	if err != nil {
		return
	}
	adminID := actx.GetAdminID()
	subject := util.FormatRbacSubject(adminID)
	isPermitted, err = s.opt.Options.Rbac.Enforce(subject, object, action)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error Checking Policy",
			zap.Error(err),
		)
	}
	return
}

func (s *authService) BatchPermissionCheck(ctx echo.Context, request [][]interface{}) (isPermitted bool, err error) {
	actx, err := util.NewAppContext(ctx)
	if err != nil {
		return
	}
	adminID := actx.GetAdminID()
	subject := util.FormatRbacSubject(adminID)
	permissions := [][]interface{}{}
	for _, req := range request {
		permissions = append(permissions, util.PrependArray(req, subject))
	}

	results, err := s.opt.Options.Rbac.BatchEnforce(permissions)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(actx))).Error("Error Checking Policy",
			zap.Error(err),
		)
		return
	}
	for _, r := range results {
		if r {
			isPermitted = true
			return
		}
	}
	err = errors.New("User don't has permission")
	return
}
