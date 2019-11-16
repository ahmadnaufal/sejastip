package entity

import "time"

type InvoiceNumber string

type Invoice struct {
	ID            int64         `db:"id"`
	TransactionID int64         `db:"transaction_id"`
	InvoiceCode   InvoiceNumber `db:"invoice_code"`
	CodedPrice    int64         `db:"coded_price"`
	Status        int           `db:"status"`
	ReceiptProof  string        `db:"receipt_proof"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
}
