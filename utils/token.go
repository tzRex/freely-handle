package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtSchema struct {
	UserId   uint
	UserRole uint
	Username string
	IsSuper  bool
}

var jwtCode = []byte("rex-service-go-1689266196285-4e7426cb23fd43ac")

func TokenCreate(tokenData *JwtSchema, ttl time.Duration, tags ...string) string {
	claims := jwt.MapClaims{
		"exp":      jwt.NewNumericDate(time.Now().Add(ttl)),
		"UserId":   tokenData.UserId,
		"UserRole": tokenData.UserRole,
		"Username": tokenData.Username,
		"IsSuper":  tokenData.IsSuper,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwt, err := token.SignedString(jwtCode)
	if err != nil {
		return ""
	}

	return jwt
}

func TokenParse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtCode, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
			return nil, errors.New("token.is.expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid.token")
}
