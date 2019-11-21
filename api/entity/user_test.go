package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"sejastip.id/api/entity"
)

func TestUserNormalize(t *testing.T) {
	user := entity.User{
		Email:    " test@gmail.com",
		Password: "test123",
		Name:     " Racoon Test    ",
		Phone:    "08923213213",
	}

	expected := entity.User{
		Email:    "test@gmail.com",
		Password: "test123",
		Name:     "Racoon Test",
		Phone:    "628923213213",
	}

	user.Normalize()
	assert.Equal(t, user, expected)
}

func TestUserValidate(t *testing.T) {
	user := entity.User{Name: "John 123"}
	userInvalidName := entity.User{Name: "Alfa @@@@ 123"}

	err := user.Validate()
	assert.NoError(t, err)

	err = userInvalidName.Validate()
	assert.Error(t, err)
}

func TestUserConvertToPublic(t *testing.T) {
	now := time.Now()
	user := entity.User{
		ID:          1,
		Email:       "test@gmail.com",
		Password:    "test123",
		Name:        "Racoon Test",
		Phone:       "628923213213",
		BankName:    "BANK TEST",
		LastLoginAt: &now,
	}

	expected := &entity.UserPublic{
		ID:       1,
		Email:    "test@gmail.com",
		Name:     "Racoon Test",
		Phone:    "628923213213",
		BankName: "BANK TEST",
	}

	userPublic := user.ConvertToPublic()
	assert.Equal(t, userPublic, expected)
}
