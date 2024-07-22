package tokenverify

import (
	"github.com/golang-jwt/jwt/v4"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/utils"
	"time"
)

type Claims struct {
	UserID     int64
	PlatformID int // login platform
	SuperAdmin bool
	jwt.RegisteredClaims
}

func BuildClaims(uid int64, platformID int, superAdmin bool, ttl int64) Claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return Claims{
		UserID:     uid,
		PlatformID: platformID,
		SuperAdmin: superAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        // Issuing time
			NotBefore: jwt.NewNumericDate(before),                                     // Begin Effective time
		},
	}
}

func GetClaimFromToken(tokensString string, secretFunc jwt.Keyfunc) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secretFunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, utils.Wrap(errs.ErrTokenMalformed, "")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.Wrap(errs.ErrTokenExpired, "")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, utils.Wrap(errs.ErrTokenNotValidYet, "")
			} else {
				return nil, utils.Wrap(errs.ErrTokenUnknown, "")
			}
		} else {
			return nil, utils.Wrap(errs.ErrTokenUnknown, "")
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.Wrap(errs.ErrTokenUnknown, "")
	}
}
