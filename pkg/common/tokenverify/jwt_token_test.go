package tokenverify

import (
	"github.com/golang-jwt/jwt/v4"
	"orderin-server/pkg/common/constant"
	"testing"
)

var secret = "jwt_secret"

func Test_ParseToken(t *testing.T) {
	uid := int64(1)
	claims1 := BuildClaims(uid, constant.AndroidPadPlatformID, true, 10)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	claim2, err := GetClaimFromToken(tokenString, secretFun())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claim2)
}

func secretFun() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
