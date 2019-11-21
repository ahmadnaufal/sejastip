package entity_test

import (
	"testing"

	"sejastip.id/api/entity"

	"github.com/stretchr/testify/assert"
)

func TestAuthNormalize(t *testing.T) {
	creds := entity.AuthCredentials{
		Email:    "  testcred@gmail.com   ",
		Password: "  testcred  ",
	}

	creds.Normalize()
	assert.Equal(t, creds.Email, "testcred@gmail.com")
	assert.Equal(t, creds.Password, "testcred")
}
