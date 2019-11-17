package usecase

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

// UserProvider is a wrapper of dependencies used by the implementation of UserUsecase
type UserProvider struct {
	UserRepository api.UserRepository
}

type userUsecase struct {
	*UserProvider
}

// NewUserUsecase creates an instance of our UserUsecase
func NewUserUsecase(pvd *UserProvider) api.UserUsecase {
	return &userUsecase{pvd}
}

func (u *userUsecase) Register(ctx context.Context, user *entity.User) (*entity.UserPublic, error) {
	user.Normalize()
	if err := user.Validate(); err != nil {
		return nil, api.ValidationError(err)
	}

	// Check if there is an existing user with the specified email
	existingUser, err := u.UserProvider.UserRepository.GetUserByEmail(ctx, user.Email)
	if err != nil && err != api.ErrNotFound {
		return nil, errors.Wrap(err, "Error fetching user by email")
	}
	if existingUser != nil {
		return nil, api.CustomValidationError("Email %s sudah digunakan", user.Email)
	}

	// validation complete, register the user
	// encrypt the password first
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 16)
	if err != nil {
		return nil, errors.Wrap(err, "Error encrypting password")
	}
	user.Password = string(hashedPassword)

	// save the user to repository
	err = u.UserProvider.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	publicUser := user.ConvertToPublic()
	return publicUser, nil
}

// GetUser get a single user by ID
func (u *userUsecase) GetUser(ctx context.Context, ID int64) (*entity.UserPublic, error) {
	user, err := u.UserProvider.UserRepository.GetUser(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "Error fetching user")
	}

	publicUser := user.ConvertToPublic()
	return publicUser, nil
}
