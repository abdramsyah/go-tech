package appcontext

import (
	"emoney-backoffice/config"
	"emoney-backoffice/internal/app/driver"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// AppContext the app context struct
type AppContext struct {
	config config.ConfigObject
}

// NewAppContext initiate appcontext object
func NewAppContext(config config.ConfigObject) *AppContext {
	return &AppContext{
		config: config,
	}
}

func (a *AppContext) GetDBInstance() (db *gorm.DB, err error) {
	dbOption := a.getPostgreOption()
	db, err = driver.NewPostgreDatabase(dbOption)

	return
}

func (a *AppContext) getPostgreOption() driver.DBPostgreOption {
	return driver.DBPostgreOption{
		Host:        a.config.DBHost,
		Port:        a.config.DBPort,
		Username:    a.config.DBUsername,
		Password:    a.config.DBPassword,
		DBName:      a.config.DBName,
		MaxPoolSize: a.config.DBMaxPoolSize,
		BatchSize:   a.config.DBBatchSize,
	}
}

// GetCachePool get cache pool connection
func (a *AppContext) GetCachePool() *redis.Pool {
	return driver.NewCache(a.getCacheOption())
}

func (a *AppContext) getCacheOption() driver.CacheOption {
	return driver.CacheOption{
		Host:               a.config.RedisHost,
		Port:               a.config.RedisPort,
		Namespace:          a.config.RedisNamespace,
		Password:           a.config.RedisPassword,
		DialConnectTimeout: cast.ToDuration(a.config.RedisDialConnectTimeout),
		ReadTimeout:        cast.ToDuration(a.config.RedisReadTimeout),
		WriteTimeout:       cast.ToDuration(a.config.RedisWriteTimeout),
		IdleTimeout:        cast.ToDuration(a.config.RedisIdleTimeout),
		MaxConnLifetime:    cast.ToDuration(a.config.RedisConnLifetimeMax),
		MaxIdle:            a.config.RedisConnIdleMax,
		MaxActive:          a.config.RedisConnActiveMax,
		Wait:               a.config.RedisIsWait,
	}
}

func (a *AppContext) GetRbacOption(db *gorm.DB) (enforcer *casbin.SyncedEnforcer, err error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return
	}
	enforcer, err = casbin.NewSyncedEnforcer(a.config.CasbinModelPath, adapter)

	return
}
