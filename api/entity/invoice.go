package entity

import (
	"time"

	"github.com/pkg/errors"
)

type InvoiceNumber string

type InvoiceStatus int

const (
	InvoiceStatusPending = iota
	InvoiceStatusPaid
	InvoiceStatusExpired
)

var mapInvoiceStatusToString = map[InvoiceStatus]string{
	InvoiceStatusPending: "pending",
	InvoiceStatusPaid:    "paid",
	InvoiceStatusExpired: "expired",
}

type Invoice struct {
	ID            int64         `db:"id"`
	TransactionID int64         `db:"transaction_id"`
	InvoiceCode   InvoiceNumber `db:"invoice_code"`
	CodedPrice    int64         `db:"coded_price"`
	PaymentMethod string        `db:"payment_method"`
	Status        InvoiceStatus `db:"status"`
	PaidAt        *time.Time    `db:"paid_at"`
	ReceiptProof  string        `db:"receipt_proof"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
}

func (i *Invoice) ConvertToPublic() InvoicePublic {
	return InvoicePublic{
		ID:            i.ID,
		TransactionID: i.TransactionID,
		InvoiceCode:   i.InvoiceCode,
		CodedPrice:    i.CodedPrice,
		PaymentMethod: i.PaymentMethod,
		Status:        mapInvoiceStatusToString[i.Status],
		PaidAt:        i.PaidAt,
		ReceiptProof:  i.ReceiptProof,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

type InvoicePublic struct {
	ID            int64         `json:"id"`
	TransactionID int64         `json:"transaction_id"`
	InvoiceCode   InvoiceNumber `json:"invoice_code"`
	CodedPrice    int64         `json:"coded_price"`
	PaymentMethod string        `json:"payment_method"`
	Status        string        `json:"status"`
	PaidAt        *time.Time    `json:"paid_at"`
	ReceiptProof  string        `json:"receipt_proof"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type InvoiceCreateForm struct {
	TransactionID int64  `json:"transaction_id"`
	PaymentMethod string `json:"payment_method"`
}

func (f *InvoiceCreateForm) Validate() error {
	if f.TransactionID < 1 {
		return errors.New("Transaction ID tidak valid")
	}

	if len(f.PaymentMethod) < 2 {
		return errors.New("Metode pembayaran tidak valid")
	}

	return nil
}

type InvoiceUpdateForm struct {
	Status       string `json:"status"`
	ReceiptProof string `json:"receipt_proof"`
}
