package common

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

// jwt加密密钥
var jwtKey = []byte(viper.GetString("jwt.secret"))

// Claims token的claim
type Claims struct {
	UserId uint
	jwt.RegisteredClaims
}

// ReleaseToken 发放token
func ReleaseToken(userId uint) (string, error) {

	//token的有效期
	expirationTime := time.Now().Add(8 * time.Hour)

	claims := &Claims{

		//自定义字段
		UserId: userId,
		//标准字段
		RegisteredClaims: jwt.RegisteredClaims{

			//过期时间
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			//发放的时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//发放者
			Issuer: "127.0.0.1",
			//主题
			Subject: "user token",
		},
	}

	//使用jwt密钥生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	//返回token
	return tokenString, nil
}

// ParseToken 从tokenString中解析出claims并返回
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
