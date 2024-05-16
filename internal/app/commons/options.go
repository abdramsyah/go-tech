package commons

import (
	"go-tech/config"
	"go-tech/internal/app/appcontext"
	"github.com/casbin/casbin/v2"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Options common option for all object that needed
type Options struct {
	Config    config.ConfigObject
	DB        *gorm.DB
	Logger    *zap.Logger
	CachePool *redis.Pool
	Rbac      *casbin.SyncedEnforcer
	Errors    []error
}

func InitCommonOptions(options ...func(*Options)) *Options {
	opt := &Options{}
	for _, o := range options {
		o(opt)
		if opt.Errors != nil {
			return opt
		}
	}
	return opt
}

func WithConfig(cfg config.ConfigObject) func(*Options) {
	return func(opt *Options) {
		opt.Config = cfg
	}
}

func WithDB(appCtx *appcontext.AppContext) func(*Options) {
	return func(opt *Options) {
		db, err := appCtx.GetDBInstance()
		if err != nil {
			opt.Errors = append(opt.Errors, err)
			return
		}
		opt.DB = db
	}
}

func WithLogger(logger *zap.Logger) func(*Options) {
	return func(opt *Options) {
		opt.Logger = logger
	}
}

func WithCache(appCtx *appcontext.AppContext) func(*Options) {
	return func(opt *Options) {
		cache := appCtx.GetCachePool()
		opt.CachePool = cache
	}
}

//Must call after WithDB to prevent nil pointer exception
func WithRBAC(appCtx *appcontext.AppContext) func(*Options) {
	return func(opt *Options) {
		rbac, err := appCtx.GetRbacOption(opt.DB)
		if err != nil {
			opt.Errors = append(opt.Errors, err)
			return
		}
		opt.Rbac = rbac
	}
}
