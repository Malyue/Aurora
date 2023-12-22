package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Config struct {
	Secret                     string `yaml:"secret"`
	AccessTokenExpireDuration  int    `yaml:"accessTokenExpireDuration"`
	RefreshTokenExpireDuration int    `yaml:"refreshTokenExpireDuration"`
}

type Claims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var jwtConfig *Config

func InitJWTConfig(cfg *Config) {
	if cfg == nil {
		jwtConfig = cfg
	}
	return
}

// GenerateToken generate accessToken,refreshToken
func GenerateToken(id string) (accessToken, refreshToken string, err error) {
	if jwtConfig == nil {
		return "", "", errors.New("empty config")
	}
	now := time.Now()
	expireTime := now.Add(time.Duration(jwtConfig.AccessTokenExpireDuration) * time.Second)
	refreshTokenExpireTime := now.Add(time.Duration(jwtConfig.RefreshTokenExpireDuration) * time.Second)
	claims := Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	// generate AccessToken
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtConfig.Secret)
	if err != nil {
		return "", "", err
	}

	// generate refresh token
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: refreshTokenExpireTime.Unix(),
	}).SignedString(jwtConfig.Secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

// ParseTokenAndValidExpire parse token and valid expire
func ParseTokenAndValidExpire(token string) (claims *Claims, expire bool, err error) {
	claims, err = ParseToken(token)
	if err != nil {
		return
	}
	return claims, claims.ExpiresAt < time.Now().Unix(), nil
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtConfig.Secret, nil
	})
	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func ParseRefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	refreshClaim, err := ParseToken(refreshToken)
	if err != nil {
		return
	}

	// if accessToken or refreshToken is not expired, refresh the token
	if refreshClaim.ExpiresAt > time.Now().Unix() {
		return GenerateToken(refreshClaim.Id)
	}

	return "", "", errors.New("the token is expire")
}
