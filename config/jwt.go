package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const SECRET_KEY = "y0ur53cr3tk3y"

type CustomClaims struct {
	TeacherID string `json:"teacherId"`
	Role      string `json:"role"`
	*jwt.RegisteredClaims
}

func GenerateToken(teacherId, role string) (string, error) {
	claims := &CustomClaims{
		teacherId,
		role,
		&jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "SIGAP SYSTEM",
			Subject:   teacherId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func ValidateToken(tokenString string) (bool, error) {
	_, err := ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRole(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
