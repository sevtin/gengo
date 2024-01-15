package constant

import "time"

const (
	CONST_DURATION_JWT_ACCESS_TOKEN_EXPIRE_IN_SECOND  = 60 * 60 * 24 * 7 * time.Second  //ACCESS_TOKEN有效期
	CONST_DURATION_JWT_REFRESH_TOKEN_EXPIRE_IN_SECOND = 60 * 60 * 24 * 30 * time.Second //REFRESH_TOKEN有效期
)
