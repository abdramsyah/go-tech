package middleware

import (
	"emoney-backoffice/internal/app/constant"
	"emoney-backoffice/internal/app/dto"
	"emoney-backoffice/internal/app/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (m *CustomMiddleware) JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		mapClaims, err := m.Service.Auth.ValidateToken(c, c.Request())
		if err != nil {
			resp := dto.FailedHttpResponse("", constant.ErrUnauthorized, nil)
			return c.JSON(http.StatusUnauthorized, resp)
		}
		adminID, ok := mapClaims[constant.AdminIDContextKey].(float64)
		if !ok {
			resp := dto.FailedHttpResponse("", "Invalid token claim", nil)
			return c.JSON(http.StatusUnauthorized, resp)
		}

		actx := util.NewEmptyAppContext(c)
		actx.SetAdminID(uint64(adminID))
		accessUUID := mapClaims[constant.AccessUUIDContextKey].(string)
		actx.SetAccessUUID(accessUUID)
		return next(actx)
	}
}
