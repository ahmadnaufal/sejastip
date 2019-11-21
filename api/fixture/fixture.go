package fixture

import (
	"time"

	"sejastip.id/api/entity"
)

// StubbedUser create a stubbed user
func StubbedUser() entity.User {
	now := time.Now()
	return entity.User{
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
func StubbedBank() entity.Bank {
	now := time.Now()
	return entity.Bank{
		Name:      "Bank Krud",
		Image:     "https://sejastip.id/img/rockybalboa.jpg",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// StubbedInvoice create a stubbed invoice row
func StubbedInvoice() entity.Invoice {
	now := time.Now()
	return entity.Invoice{
		ID:            1,
		TransactionID: 1,
		Status:        entity.InvoiceStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// StubbedProduct returns a stubbed product row
func StubbedProduct() entity.Product {
	now := time.Now()
	return entity.Product{
		ID:          1,
		Title:       "barang test",
		Description: "barang mudah",
		Price:       12000,
		Status:      entity.ProductStatusOffered,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// StubbedCountry returns a stubbed country row
func StubbedCountry() entity.Country {
	return entity.Country{
		ID:   1,
		Name: "Indonesia",
	}
}
