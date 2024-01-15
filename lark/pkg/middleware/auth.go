package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"lark/pkg/common/xjwt"
	"lark/pkg/common/xlog"
	"lark/pkg/common/xredis"
	"lark/pkg/constant"
	"lark/pkg/utils"
	"lark/pkg/xhttp"
	"strings"
)

func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err          error
			token        *jwt.Token
			ok           bool
			uid          interface{}
			uidVal       int64
			platform     interface{}
			sessionId    interface{}
			sessionIdVal string
			sessionIdKey string
		)
		token, err = xjwt.ParseFromCookie(ctx)
		if err != nil {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_JWT_TOKEN_ERR, xhttp.ERROR_HTTP_JWT_TOKEN_ERR)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if uid, ok = claims[constant.USER_UID]; ok == false {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_USER_ID_DOESNOT_EXIST, xhttp.ERROR_HTTP_USER_ID_DOESNOT_EXIST)
			return
		}
		if platform, ok = claims[constant.USER_PLATFORM]; ok == false || utils.TryToInt(claims[constant.USER_PLATFORM]) == 0 {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_PLATFORM_DOESNOT_EXIST, xhttp.ERROR_HTTP_PLATFORM_DOESNOT_EXIST)
			return
		}
		ctx.Set(constant.USER_UID, uid)
		ctx.Set(constant.USER_PLATFORM, platform)
		if strings.HasPrefix(ctx.FullPath(), constant.API_PUBLIC) && ctx.Request.Method == constant.HTTP_REQUEST_METHOD_GET {
			ctx.Next()
			return
		}
		if sessionId, ok = claims[constant.USER_JWT_SESSION_ID]; ok == false {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_SESSION_ID_DOESNOT_EXIST, xhttp.ERROR_HTTP_SESSION_ID_DOESNOT_EXIST)
			return
		}
		uidVal, err = utils.ToInt64(uid)
		if err != nil {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_USER_ID_DOESNOT_EXIST, xhttp.ERROR_HTTP_USER_ID_DOESNOT_EXIST)
			return
		}
		sessionIdKey = constant.RK_SYNC_USER_ACCESS_TOKEN_SESSION_ID + utils.GetHashTagKey(uidVal) + ":" + utils.ToString(platform)
		if sessionIdVal, err = xredis.Get(sessionIdKey); err != nil {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_READ_SESSION_ID_FAILED, xhttp.ERROR_HTTP_READ_SESSION_ID_FAILED)
			xlog.Warn(xhttp.ERROR_CODE_HTTP_READ_SESSION_ID_FAILED, xhttp.ERROR_HTTP_READ_SESSION_ID_FAILED, err.Error())
			return
		}
		if sessionIdVal == "" {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_TOKEN_EXPIRES, xhttp.ERROR_HTTP_TOKEN_EXPIRES)
		}
		if sessionIdVal != utils.ToString(sessionId) {
			ctx.Abort()
			xhttp.Error(ctx, xhttp.ERROR_CODE_HTTP_ACCOUNT_KICKED_OUT, xhttp.ERROR_HTTP_ACCOUNT_KICKED_OUT)
			return
		}
		ctx.Next()
	}
}
