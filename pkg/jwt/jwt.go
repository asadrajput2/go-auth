package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userId uint8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret")) // TODO: change secret

	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func ValidateToken(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil // TODO: change secret
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("claims: ", claims)
		return claims["user_id"], nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
