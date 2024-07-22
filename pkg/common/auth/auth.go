package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"orderin-server/pkg/common/cache"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/tokenverify"
	"orderin-server/pkg/common/utils"
)

type AuthService interface {
	// 结果为空 不返回错误
	GetTokensWithoutError(ctx context.Context, userID int64, platformID int) (map[string]int, error)
	// 创建token
	CreateToken(ctx context.Context, userID int64, superAdmin bool, platformID int) (string, error)
	//删除Token
	DeleteToken(ctx context.Context, userID int64, platformID int) error
}

type authService struct {
	cache        cache.TokenCache
	accessSecret string
	accessExpire int64
}

func NewAuthService(cache cache.TokenCache, accessSecret string, accessExpire int64) AuthService {
	return &authService{cache: cache, accessSecret: accessSecret, accessExpire: accessExpire}
}

// 结果为空 不返回错误.
func (a *authService) GetTokensWithoutError(
	ctx context.Context,
	userID int64,
	platformID int,
) (map[string]int, error) {
	return a.cache.GetTokensWithoutError(ctx, userID, platformID)
}

// 创建token.
func (a *authService) CreateToken(ctx context.Context, userID int64, superAdmin bool, platformID int) (string, error) {
	tokens, err := a.cache.GetTokensWithoutError(ctx, userID, platformID)
	if err != nil {
		return "", err
	}
	var deleteTokenKey []string
	for k, v := range tokens {
		_, err = tokenverify.GetClaimFromToken(k, Secret())
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err := a.cache.DeleteTokenByUidPid(ctx, userID, platformID, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}
	claims := tokenverify.BuildClaims(userID, platformID, superAdmin, a.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return tokenString, a.cache.AddTokenFlag(ctx, userID, platformID, tokenString, constant.NormalToken)
}

func (a *authService) DeleteToken(ctx context.Context, userID int64, platformID int) error {
	tokens, err := a.cache.GetTokensWithoutError(ctx, userID, platformID)
	if err != nil {
		return err
	}
	var deleteTokenKey []string
	for k, v := range tokens {
		_, err = tokenverify.GetClaimFromToken(k, Secret())
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err := a.cache.DeleteTokenByUidPid(ctx, userID, platformID, deleteTokenKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Secret), nil
	}
}
