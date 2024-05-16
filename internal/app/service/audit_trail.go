package service

import (
	"emoney-backoffice/internal/app/commons"
	"emoney-backoffice/internal/app/constant"
	"emoney-backoffice/internal/app/dto"
	"emoney-backoffice/internal/app/model"
	"emoney-backoffice/internal/app/util"
	"errors"
	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
)

type IAuditTrailService interface {
	Create(ctx echo.Context, req *dto.AuditTrailRequest) (err error)
	FindRoutes(ctx echo.Context) (routes []dto.BORoutesResponse)
	Find(ctx echo.Context, pConfig commons.PaginationConfig, filter *dto.AuditTrailFilter) (list []model.AuditTrails, count int64, httpStatus int, err error)
}

type auditTrailService struct {
	opt Option
}

func NewAuditTrailService(opt Option) IAuditTrailService {
	return &auditTrailService{
		opt: opt,
	}
}

func (s *auditTrailService) Create(ctx echo.Context, req *dto.AuditTrailRequest) (err error) {
	admin, err := s.opt.Repository.Admin.FindById(uint64(req.AdminID))
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error get admin by ID",
			zap.Error(err),
		)
		return
	}

	data := model.AuditTrails{
		AdminID:    int64(req.AdminID),
		AdminEmail: admin.Email,
		AdminName:  admin.FullName,
		AdminRole:  admin.Role.Name,
		Action:     req.RouteName,
		URL:        req.URL,
		CreatedAt:  time.Now(),
		RequestID:  null.NewString(util.GetRequestID(ctx), true),
	}

	err = s.opt.Repository.AuditTrail.Create(ctx, &data)
	if err != nil {
		s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Error save audit trail",
			zap.Error(err),
		)
	}
	return
}

func (s *auditTrailService) FindRoutes(ctx echo.Context) (routes []dto.BORoutesResponse) {
	allRoutes := ctx.Echo().Routes()
	for _, r := range allRoutes {
		if strings.Contains(r.Name, "github.com/") || strings.Contains(r.Path, "/mobile/") || strings.Contains(r.Path, "emoney-backoffice/") {
			continue
		}
		routes = append(routes, dto.BORoutesResponse{
			Name:   r.Name,
			Path:   r.Path,
			Method: r.Method,
		})
	}
	return
}

func (s *auditTrailService) Find(ctx echo.Context, pConfig commons.PaginationConfig, filter *dto.AuditTrailFilter) (list []model.AuditTrails, count int64, httpStatus int, err error) {
	var waitGroup sync.WaitGroup
	c := make(chan error)

	waitGroup.Add(2)

	go func() {
		waitGroup.Wait()
		close(c)
	}()

	go func() {
		defer waitGroup.Done()

		count, err = s.opt.Repository.AuditTrail.Count(ctx, filter)
		if err != nil {
			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Get Audit Trail Count",
				zap.Error(err),
			)
			err = errors.New(constant.ErrFailedGetDataCount)
			c <- err
		}
	}()

	go func() {
		defer waitGroup.Done()

		list, err = s.opt.Repository.AuditTrail.Find(ctx, pConfig, filter)
		if err != nil {
			s.opt.Logger.With(zap.String("RequestID", util.GetRequestID(ctx))).Error("Get Audit Trail",
				zap.Error(err),
			)
			err = errors.New(constant.ErrFailedFetchData)
			c <- err
		}
	}()

	for errChan := range c {
		if errChan != nil {
			err = errChan
			httpStatus = http.StatusInternalServerError
			return
		}
	}

	return
}
