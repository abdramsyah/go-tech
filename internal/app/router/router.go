package router

import (
	"emoney-backoffice/internal/app/handler"
	"emoney-backoffice/internal/app/middleware"
	"github.com/labstack/echo/v4"
)

type Router struct {
	health *health
	auth   *auth
}

func NewRouter(server *echo.Echo, handlers handler.Handlers, cmiddleware *middleware.CustomMiddleware) (router *Router) {
	health := newHealth(server, handlers, cmiddleware)
	auth := newAuth(server, handlers, cmiddleware)
	return &Router{
		health: health,
		auth:   auth,
	}
}

func (r *Router) Initialize() {
	r.health.initialize()
	r.auth.initialize()
}
