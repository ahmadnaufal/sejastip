package entity

import "time"

type Transaction struct {
	ID             int64      `db:"id"`
	ProductID      int64      `db:"product_id"`
	BuyerID        int64      `db:"buyer_id"`
	SellerID       int64      `db:"seller_id"`
	BuyerAddressID int64      `db:"buyer_address_id"`
	Quantity       uint       `db:"quantity"`
	Notes          string     `db:"notes"`
	TotalPrice     int64      `db:"total_price"`
	Status         int        `db:"status"`
	PaidAt         *time.Time `db:"paid_at"`
	FinishedAt     *time.Time `db:"finished_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

type TransactionPublic struct {
	ID           int64              `json:"id"`
	Product      *ProductPublic     `json:"product"`
	Buyer        *UserPublic        `json:"buyer"`
	BuyerAddress *UserAddressPublic `json:"buyer_address"`
	Quantity     uint               `json:"quantity"`
	Notes        string             `json:"notes"`
	TotalPrice   int64              `json:"total_price"`
	Status       string             `json:"status"`
	PaidAt       *time.Time         `json:"paid_at"`
	FinishedAt   *time.Time         `json:"finished_at"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type TransactionForm struct {
	ProductID int64  `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	AddressID int64  `json:"address_id"`
	Notes     string `json:"notes"`
}
