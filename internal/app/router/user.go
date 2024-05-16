package router

import (
	"go-tech/internal/app/handler"
	"go-tech/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

type user struct {
	server     *echo.Echo
	handlers   handler.Handlers
	middleware *middleware.CustomMiddleware
}

func newUser(server *echo.Echo, handlers handler.Handlers, middleware *middleware.CustomMiddleware) *user {
	return &user{
		server:     server,
		handlers:   handlers,
		middleware: middleware,
	}
}

func (h *user) initialize() {
	g := h.server.Group("/user")
	// g.POST("/login", middleware.HandlerWrapperJson(h.handlers.Auth.Login))
	// g.POST("/refresh", middleware.HandlerWrapperJson(h.handlers.Auth.RefreshToken))
	// g.Use(h.middleware.JWTMiddleware)
	// g.Use(h.middleware.AuditTrailMiddleware)
	g.GET("/logout", middleware.HandlerWrapperJson(h.handlers.Auth.Logout)).Name = "Logout"
}
