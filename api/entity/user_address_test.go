package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"sejastip.id/api/entity"
)

func TestUserAddressConvertToPublic(t *testing.T) {
	now := time.Now()
	ua := entity.UserAddress{
		ID:          1,
		Address:     "Jalan test",
		Phone:       "08921321321",
		AddressName: "test",
		UserID:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	expected := entity.UserAddressPublic{
		ID:          1,
		Address:     "Jalan test",
		Phone:       "08921321321",
		AddressName: "test",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	uap := ua.ConvertToPublic()

	assert.Equal(t, uap, expected)
}
