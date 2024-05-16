package util

import (
	"emoney-backoffice/internal/app/constant"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
)

type AppContext struct {
	echo.Context
	adminID    *uint64
	accessUUID *string
}

func (c *AppContext) GetRequestID() string {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		return uuid.NewV4().String()
	}
	return requestID
}

func (c *AppContext) SetAdminID(adminID uint64) {
	c.adminID = &adminID
}

func (c *AppContext) SetAccessUUID(accessUUID string) {
	c.accessUUID = &accessUUID
}

func (c *AppContext) GetAdminID() uint64 {
	return *c.adminID
}

func (c *AppContext) GetAccessUUID() string {
	return *c.accessUUID
}

func NewEmptyAppContext(parent echo.Context) *AppContext {
	return &AppContext{parent, nil, nil}
}

func NewAppContext(parent echo.Context) (*AppContext, error) {
	pctx, ok := parent.(*AppContext)
	if !ok {
		return nil, errors.New(constant.ErrUnauthorized)
	}
	if pctx.adminID == nil || pctx.accessUUID == nil {
		return nil, errors.New(constant.ErrUnauthorized)
	}
	return pctx, nil
}
