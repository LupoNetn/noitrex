package utils

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	OperatorID string `json:"sub"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwtPair(JWTAccessSecret string, JWTRefreshSecret string, name string, operatorID string, role string) (string, string, error) {
	accessClaims := Claims{
		OperatorID: operatorID,
		Name:       name,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "noitrex",
			Subject:   operatorID,
		},
	}
	refreshClaims := Claims{
		OperatorID: operatorID,
		Name:       name,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "noitrex",
			Subject:   operatorID,
		},
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err := access.SignedString(JWTAccessSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := refresh.SignedString(JWTRefreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func VerifyJwt() {}
