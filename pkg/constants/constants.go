package constants

import "time"

type CtxKey string

var (
	AccessTokenHeader = "Access-Token"

	CtxUserIDKey CtxKey = "user_id"

	DefaultClientTimeout = time.Second * 5
)
