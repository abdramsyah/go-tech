package router

import (
	"emoney-backoffice/internal/app/handler"
	"emoney-backoffice/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

type auth struct {
	server     *echo.Echo
	handlers   handler.Handlers
	middleware *middleware.CustomMiddleware
}

func newAuth(server *echo.Echo, handlers handler.Handlers, middleware *middleware.CustomMiddleware) *auth {
	return &auth{
		server:     server,
		handlers:   handlers,
		middleware: middleware,
	}
}

func (h *auth) initialize() {
	g := h.server.Group("/auth")
	g.POST("/login", middleware.HandlerWrapperJson(h.handlers.Auth.Login))
	g.POST("/refresh", middleware.HandlerWrapperJson(h.handlers.Auth.RefreshToken))
	g.Use(h.middleware.JWTMiddleware)
	g.Use(h.middleware.AuditTrailMiddleware)
	g.GET("/logout", middleware.HandlerWrapperJson(h.handlers.Auth.Logout)).Name = "Logout"
}
