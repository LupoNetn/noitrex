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

	accessToken, err := access.SignedString([]byte(JWTAccessSecret))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := refresh.SignedString([]byte(JWTRefreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func VerifyJwt(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
