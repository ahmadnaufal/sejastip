package usecase

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

// AuthProvider is a wrapper of dependencies used by the implementation of AuthUsecase
type AuthProvider struct {
	UserRepository api.UserRepository
	JWTPrivateKey  string
}

type authUsecase struct {
	*AuthProvider
}

// NewAuthUsecase creates an instance of AuthUsecase
func NewAuthUsecase(pvd *AuthProvider) api.AuthUsecase {
	return &authUsecase{pvd}
}

// AuthenticateUser handles user authentication based on the provided credentials
func (u *authUsecase) AuthenticateUser(ctx context.Context, auth *entity.AuthCredentials) (*entity.AuthResponse, error) {
	// get the user first
	user, err := u.AuthProvider.UserRepository.GetUserByEmail(ctx, auth.Email)
	if err != nil {
		return nil, errors.Wrap(err, "Authentication failed")
	}

	// Check for user password validity
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {
		return nil, api.ErrInvalidCredentials
	}

	// After all credentials are valid, we create a claim to store all user data
	createdAt := time.Now()
	expirationTime := time.Now().Add(96 * time.Hour)
	claims := entity.ResourceClaims{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Phone:        user.Phone,
		RegisteredAt: user.CreatedAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Then we create the token with HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// And sign it with our private key
	tokenString, err := token.SignedString([]byte(u.AuthProvider.JWTPrivateKey))
	if err != nil {
		return nil, err
	}

	authResponse := &entity.AuthResponse{
		Token:     tokenString,
		CreatedAt: createdAt,
		ExpiredAt: expirationTime,
	}
	return authResponse, nil
}
