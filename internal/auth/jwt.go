package auth

import (
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// todo: change this to your secret key
var JwtKey = []byte("your_secret_key") // Change this to a strong secret key

// Claims defines the structure for the JWT claims
type Claims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	jwt.RegisteredClaims
}

func Login(user entities.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		ID:    user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
