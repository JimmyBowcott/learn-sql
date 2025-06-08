package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Username string `json:"username,omitempty"`
	Level    int    `json:"level"`
	jwt.RegisteredClaims
}

func GetAuthToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func GetClaims(r *http.Request) Claims {
	token, err := DecodeToken(GetAuthToken(r))
	if err != nil {
		return Claims{Username: "", Level: 1}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{Username: "", Level: 1}
	}

	username, _ := claims["username"].(string)
	levelFloat, ok := claims["level"].(float64)
	level := 1
	if ok {
		level = int(levelFloat)
	}

	return Claims{
		Username: username,
		Level:    level,
	}
}

func DecodeToken(token string) (*jwt.Token, error) {
	res, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !res.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return res, nil
}

func GenerateToken(username string, level int) (string, error) {
	claims := Claims{
		Username: username,
		Level:    level,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
