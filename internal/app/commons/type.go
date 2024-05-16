package commons

type PaginationConfig struct {
	Offset   int
	PageSize int
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AuthCacheValue struct {
	AdminID uint64 `json:"admin_id"`
}
