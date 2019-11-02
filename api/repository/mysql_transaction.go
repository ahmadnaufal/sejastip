package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"sejastip.id/api/entity"

	"sejastip.id/api"

	"github.com/jmoiron/sqlx"
)

type mysqlTransaction struct {
	db *sqlx.DB
}

var transactionAllowedFilters = map[string]struct{}{}

func NewMysqlTransaction(db *sql.DB) api.TransactionRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlTransaction{newDB}
}

func (m *mysqlTransaction) CreateTransaction(ctx context.Context, transaction *entity.Transaction) error {
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	query := `INSERT INTO transactions
		(product_id, buyer_id, seller_id, buyer_address_id, quantity,
			notes, total_price, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error preparing insert transaction query")
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx,
		transaction.ProductID, transaction.BuyerID, transaction.SellerID,
		transaction.BuyerAddressID, transaction.Quantity, transaction.Notes,
		transaction.TotalPrice, transaction.CreatedAt, transaction.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "error executing insert transaction query")
	}

	transaction.ID, err = res.LastInsertId()
	return err
}

func (m *mysqlTransaction) GetTransaction(ctx context.Context, transactionID int64) (*entity.Transaction, error) {
	query := `
		SELECT * FROM transactions
		WHERE id = ?
	`
	result := &entity.Transaction{}
	err := m.db.GetContext(ctx, result, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		} else {
			return nil, err
		}
	}

	return result, nil
}

func (m *mysqlTransaction) GetTransactions(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.Transaction, int64, error) {
	return nil, 0, nil
}
