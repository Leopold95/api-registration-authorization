package shared

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	secret     string
	expiration time.Duration
}

func NewJwtService() *JwtService {
	return &JwtService{
		secret:     "secret",
		expiration: time.Minute * 15,
	}
}

func (this *JwtService) CreateToken(model *UserModel) (string, error) {
	claims := jwt.MapClaims{
		"sub":   model.Id.String(),
		"email": model.Email,
		"name":  model.Name,
		"exp":   time.Now().Add(this.expiration).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(this.secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (this *JwtService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(this.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
