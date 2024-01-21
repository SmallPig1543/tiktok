package util

import (
	"github.com/golang-jwt/jwt"
	"tiktok/pkg/constants"
	"time"
)

type Claims struct {
	ID       uint
	UserName string
	jwt.StandardClaims
}

// GenerateAccessToken 签发access_token
func GenerateAccessToken(uid uint, username string) (accessToken string, err error) {
	claims := &Claims{
		ID:       uid,
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (time.Now().Add(3 * time.Hour)).Unix(),
			Issuer:    "tiktok",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(constants.JwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// GenerateRefreshToken 签发refresh_token
func GenerateRefreshToken(uid uint, username string) (accessToken string, err error) {
	claims := &Claims{
		ID:       uid,
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: (time.Now().Add(24 * time.Hour)).Unix(),
			Issuer:    "tiktok",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(constants.JwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
func ParseToken(token string) (*Claims, bool, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(constants.JwtSecret), nil
	})

	if err != nil {
		return nil, false, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, time.Now().Unix() < claims.ExpiresAt, nil
	}
	return nil, false, err
}
