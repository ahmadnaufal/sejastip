package entity

import (
	"time"

	"github.com/pkg/errors"
)

const (
	TransactionStatusInit = iota
	TransactionStatusPaid
	TransactionStatusInProgress
	TransactionStatusDelivered
	TransactionStatusFinished
	TransactionStatusRejected
	TransactionStatusExpired
)

var mapStatusToString = map[int]string{
	TransactionStatusInit:       "placed",
	TransactionStatusPaid:       "paid",
	TransactionStatusInProgress: "in_progress",
	TransactionStatusDelivered:  "delivered",
	TransactionStatusFinished:   "finished",
	TransactionStatusRejected:   "rejected",
	TransactionStatusExpired:    "expired",
}

var MapStatusToStringReverse = map[string]int{
	"placed":      TransactionStatusInit,
	"paid":        TransactionStatusPaid,
	"in_progress": TransactionStatusInProgress,
	"delivered":   TransactionStatusDelivered,
	"finished":    TransactionStatusFinished,
	"rejected":    TransactionStatusRejected,
	"expired":     TransactionStatusExpired,
}

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

func (t *Transaction) GetStatusString() string {
	return mapStatusToString[t.Status]
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

func (f *TransactionForm) Validate() error {
	if f.ProductID < 1 {
		return errors.New("Product is not selected yet")
	}

	if f.Quantity < 1 {
		return errors.New("Ordered quantity must be more than 0")
	}

	if f.AddressID < 1 {
		return errors.New("Order address is not selected yet")
	}

	return nil
}

type UpdateTransactionForm struct {
	Status    string `json:"status"`
	AWBNumber string `json:"awb_number"`
	Courier   string `json:"courier"`
}

func (f *UpdateTransactionForm) Validate() error {
	if f.Status == "" {
		return errors.New("Status transaksi tujuan diperlukan")
	}

	if f.Status == "delivered" {
		if f.AWBNumber == "" {
			return errors.New("Nomor resi pengiriman wajib diisi")
		}

		if f.Courier == "" {
			errors.New("Kurir pengiriman perlu diisi")
		}
	}

	return nil
}
