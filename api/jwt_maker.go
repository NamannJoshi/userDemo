package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

// ongoin(proirity level - major)
func (maker *JWTMaker) createToken(id int, email string, duration time.Duration) (string,  error) {
	claims := &jwt.MapClaims{
		"id":id,
		"email": email,
		"duration": duration,
	}
	jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return "", nil
}
