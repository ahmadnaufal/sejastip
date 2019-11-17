package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"sejastip.id/api"
	"sejastip.id/api/entity"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlInvoice struct {
	db *sqlx.DB
}

// NewMysqlInvoice creates a new instance of MySQL invoice repository
func NewMysqlInvoice(db *sql.DB) api.InvoiceRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlInvoice{newDB}
}

func (m *mysqlInvoice) InsertInvoice(ctx context.Context, invoice *entity.Invoice) error {
	now := time.Now()
	invoice.CreatedAt = now
	invoice.UpdatedAt = now

	query := `INSERT INTO invoices
		(transaction_id, invoice_code, coded_price, payment_method, status,
			receipt_proof, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error preparing insert invoice query")
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx,
		invoice.TransactionID, invoice.InvoiceCode, invoice.CodedPrice, invoice.PaymentMethod,
		invoice.Status, invoice.ReceiptProof, invoice.CreatedAt, invoice.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "error executing insert invoice query")
	}

	invoice.ID, err = res.LastInsertId()
	return err
}

func (m *mysqlInvoice) GetInvoice(ctx context.Context, invoiceID int64) (*entity.Invoice, error) {
	query := `
		SELECT * FROM invoices
		WHERE id = ?
	`
	result := &entity.Invoice{}
	err := m.db.GetContext(ctx, result, query, invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}

		return nil, err
	}

	return result, nil
}

func (m *mysqlInvoice) GetInvoiceFromTransaction(ctx context.Context, transactionID int64) (*entity.Invoice, error) {
	query := `
		SELECT * FROM invoices
		WHERE transaction_id = ?
	`
	result := &entity.Invoice{}
	err := m.db.GetContext(ctx, result, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}

		return nil, err
	}

	return result, nil
}

func (m *mysqlInvoice) UpdateInvoice(ctx context.Context, invoiceID int64, invoice *entity.Invoice) error {
	now := time.Now()
	invoice.UpdatedAt = now

	query := `UPDATE invoices SET
		invoice_code = ?, coded_price = ?, payment_method = ?, status = ?,
		paid_at = ?, receipt_proof = ?, updated_at = ?
		WHERE id = ?`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error preparing update invoice")
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx,
		invoice.InvoiceCode, invoice.CodedPrice, invoice.PaymentMethod, invoice.Status,
		invoice.PaidAt, invoice.ReceiptProof, invoice.UpdatedAt, invoiceID,
	)
	if err != nil {
		return errors.Wrap(err, "error executing update invoice")
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return errors.New(fmt.Sprintf("Unexpected behavior detected when updating invoice (total rows affected: %d)", affectedRows))
	}

	return nil
}
