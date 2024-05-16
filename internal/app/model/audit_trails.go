package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


Table: audit_trails
[ 0] id                                             INT8                 null: false  primary: true   isArray: false  auto: true   col: INT8            len: -1      default: []
[ 1] admin_id                                       INT8                 null: false  primary: false  isArray: false  auto: false  col: INT8            len: -1      default: []
[ 2] admin_email                                    VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
[ 3] admin_name                                     VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
[ 4] admin_role                                     VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
[ 5] action                                         VARCHAR(150)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 150     default: []
[ 6] url                                            VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 7] created_at                                     TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [now()]


JSON Sample
-------------------------------------
{    "id": 61,    "admin_id": 51,    "admin_email": "MdMMhCDSNPfnPbwqRoyFVpqQl",    "admin_name": "CmbYLsyHsUFgRIpaIMiqxCGIU",    "admin_role": "ykqalxqyUQrQLnOBYjrjmeEkS",    "action": "qPCgohPVhXABHSvxrEDFLBnZB",    "url": "qUZZvYwPwSaIhBVHFrHKNsgDC",    "created_at": "2264-02-27T01:33:24.530226326+07:00"}



*/

// AuditTrails struct is a row record of the audit_trails table in the inacash_bo_db database
type AuditTrails struct {
	//[ 0] id                                             INT8                 null: false  primary: true   isArray: false  auto: true   col: INT8            len: -1      default: []
	ID int64 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:INT8;" json:"id"`
	//[ 1] admin_id                                       INT8                 null: false  primary: false  isArray: false  auto: false  col: INT8            len: -1      default: []
	AdminID int64 `gorm:"column:admin_id;type:INT8;" json:"admin_id"`
	//[ 2] admin_email                                    VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	AdminEmail string `gorm:"column:admin_email;type:VARCHAR;size:100;" json:"admin_email"`
	//[ 3] admin_name                                     VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	AdminName string `gorm:"column:admin_name;type:VARCHAR;size:100;" json:"admin_name"`
	//[ 4] admin_role                                     VARCHAR(100)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	AdminRole string `gorm:"column:admin_role;type:VARCHAR;size:100;" json:"admin_role"`
	//[ 5] action                                         VARCHAR(150)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 150     default: []
	Action string `gorm:"column:action;type:VARCHAR;size:150;" json:"action"`
	//[ 6] url                                            VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	URL string `gorm:"column:url;type:VARCHAR;size:255;" json:"url"`
	//[ 7] created_at                                     TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [now()]
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;" json:"created_at"`
	//[ 8] request_id                                     VARCHAR(100)         null: true   primary: false  isArray: false  auto: false  col: VARCHAR         len: 100     default: []
	RequestID null.String `gorm:"column:request_id;type:VARCHAR;size:100;" json:"request_id"`
}

// TableName sets the insert table name for this struct type
func (a *AuditTrails) TableName() string {
	return "audit_trails"
}
