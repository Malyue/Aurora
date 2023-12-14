package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtConfig struct {
	Secret          string `yaml:"secret"`
	TokenExpireTime int    `yaml:"token_expire_time"`
}

type MyClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var DefaultSecret = []byte("dapr412ei1283")

// 生成jwt
func GenToken(id string, config JwtConfig) (string, error) {
	c := MyClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.TokenExpireTime) * time.Second).Unix(),
			Issuer:    "prompting",
		},
	}

	// 使用指定的签名方法创建对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString([]byte(config.Secret))
}

// 解析JWT
func ParseToken(tokenString string, secret string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}
