package entity

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthCredentials stores user authentication credentials
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Normalize is a method to normalize all field values
func (a *AuthCredentials) Normalize() {
	a.Email = strings.TrimSpace(a.Email)
	a.Password = strings.TrimSpace(a.Password)
}

// ResourceClaims for claims
type ResourceClaims struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
	jwt.StandardClaims
}

// AuthResponse is our structure for token response after user authentication
type AuthResponse struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
