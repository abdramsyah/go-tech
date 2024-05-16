package middleware

import (
	"emoney-backoffice/internal/app/dto"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

type HandlerJson func(ctx echo.Context) (status int, resp dto.HttpResponse)

func HandlerWrapperJson(next HandlerJson) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		start := time.Now()
		status, resp := next(ctx)
		resp.ProcessTime = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
		return ctx.JSON(status, resp)
	}
}
