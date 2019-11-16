package entity

import "time"

type TransactionShipping struct {
	ID            int64     `db:"id"`
	TransactionID int64     `db:"transaction_id"`
	AWBNumber     string    `db:"awb_number"`
	Courier       string    `db:"courier"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type TransactionShippingPublic struct {
	AWBNumber string    `json:"awb_number"`
	Courier   string    `json:"courier"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
