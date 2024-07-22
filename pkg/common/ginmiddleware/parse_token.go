package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/auth"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/tokenverify"
	"orderin-server/pkg/common/utils"
)

func ParseToken(authService auth.AuthService, rdb redis.UniversalClient, excludeUris []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsContain(c.Request.RequestURI, excludeUris) {
			c.Next()
		} else {
			switch c.Request.Method {
			case http.MethodPost:
				token := c.Request.Header.Get(constant.Token)
				if token == "" {
					log.ZWarn(c, "header get token error", errs.ErrArgs.Wrap("header must have token"))
					api.GinError(c, errs.ErrArgs.Wrap("header must have token"))
					c.Abort()
					return
				}
				claims, err := tokenverify.GetClaimFromToken(token, Secret())
				if err != nil {
					log.ZWarn(c, "jwt get token error", errs.ErrTokenUnknown.Wrap())
					api.GinError(c, errs.ErrTokenUnknown.Wrap())
					c.Abort()
					return
				}
				m, err := authService.GetTokensWithoutError(c, claims.UserID, claims.PlatformID)
				if err != nil {
					log.ZWarn(c, "cache get token error", errs.ErrTokenNotExist.Wrap())
					api.GinError(c, errs.ErrTokenNotExist.Wrap())
					c.Abort()
					return
				}
				if len(m) == 0 {
					log.ZWarn(c, "cache do not exist token error", errs.ErrTokenNotExist.Wrap())
					api.GinError(c, errs.ErrTokenNotExist.Wrap())
					c.Abort()
					return
				}
				if v, ok := m[token]; ok {
					switch v {
					case constant.NormalToken:
					case constant.KickedToken:
						log.ZWarn(c, "cache kicked token error", errs.ErrTokenKicked.Wrap())
						api.GinError(c, errs.ErrTokenKicked.Wrap())
						c.Abort()
						return
					default:
						log.ZWarn(c, "cache unknown token error", errs.ErrTokenUnknown.Wrap())
						api.GinError(c, errs.ErrTokenUnknown.Wrap())
						c.Abort()
						return
					}
				} else {
					api.GinError(c, errs.ErrTokenNotExist.Wrap())
					c.Abort()
					return
				}
				c.Set(constant.OpUserPlatform, constant.PlatformIDToName(claims.PlatformID))
				c.Set(constant.OpUserID, claims.UserID)
				c.Set(constant.SuperAdmin, claims.SuperAdmin)

				c.Next()
			}
		}

	}
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Secret), nil
	}
}
