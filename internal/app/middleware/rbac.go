package middleware

import (
	"go-tech/internal/app/constant"
	"go-tech/internal/app/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (m *CustomMiddleware) RBACMiddleware(object string, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isPermitted, err := m.Service.Auth.PermissionCheck(c, object, action)
			if err != nil {
				data := dto.FailedHttpResponse("", "Failed to do permission check", nil)
				return c.JSON(http.StatusInternalServerError, data)
			}
			if !isPermitted {
				data := dto.FailedHttpResponse("", constant.ErrUnauthorized, nil)
				return c.JSON(http.StatusForbidden, data)
			}
			return next(c)
		}
	}
}

// Sample request [][]interface{}{{"13", "37"}, {"13", "39"}}
func (m *CustomMiddleware) RBACBatchMiddleware(request [][]interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			isPermitted, err := m.Service.Auth.BatchPermissionCheck(c, request)
			if err != nil {
				data := dto.FailedHttpResponse("", "Failed to do permission check", nil)
				return c.JSON(http.StatusInternalServerError, data)
			}
			if !isPermitted {
				data := dto.FailedHttpResponse("", constant.ErrUnauthorized, nil)
				return c.JSON(http.StatusForbidden, data)
			}
			return next(c)
		}
	}
}
