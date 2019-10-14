package fixture

import (
	"time"

	"sejastip.id/api"
)

// StubbedUser create a stubbed user
func StubbedUser() api.User {
	now := time.Now()
	return api.User{
		Email:       "rockybalboa@gmail.com",
		Name:        "Rocky Balboa",
		Phone:       "628961234321",
		Password:    "rockybalboa",
		BankName:    "BCA",
		BankAccount: "012341234",
		Avatar:      "https://sejastip.id/img/rockybalboa.jpg",
		LastLoginAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// StubbedBank create a stubbed bank row
func StubbedBank() api.Bank {
	now := time.Now()
	return api.Bank{
		Name:      "Bank Krud",
		Image:     "https://sejastip.id/img/rockybalboa.jpg",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
