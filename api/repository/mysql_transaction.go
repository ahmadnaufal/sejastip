package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/shuoli84/sqlm"

	"sejastip.id/api/entity"

	"sejastip.id/api"

	"github.com/jmoiron/sqlx"
)

type mysqlTransaction struct {
	db *sqlx.DB
}

var transactionAllowedFilters = map[string]struct{}{
	"product_id": struct{}{},
}

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
			notes, total_price, invoice_id, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error preparing insert transaction query")
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx,
		transaction.ProductID, transaction.BuyerID, transaction.SellerID,
		transaction.BuyerAddressID, transaction.Quantity, transaction.Notes,
		transaction.TotalPrice, transaction.InvoiceID, transaction.CreatedAt, transaction.UpdatedAt,
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
	filteredQueries := buildTransactionDynamicQuery(filter)
	countQuery, countArgs := sqlm.Build(
		"SELECT COUNT(id) FROM transactions",
		"WHERE", sqlm.And(filteredQueries),
	)

	var count int64
	err := m.db.GetContext(ctx, &count, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	query, args := sqlm.Build(
		"SELECT * FROM transactions",
		"WHERE", sqlm.And(filteredQueries),
		"ORDER BY updated_at DESC",
		sqlm.Exp("LIMIT", sqlm.P(offset), ",", sqlm.P(limit)),
	)
	results := []entity.Transaction{}
	err = m.db.SelectContext(ctx, &results, query, args...)
	return results, count, err
}

func buildTransactionDynamicQuery(filter entity.DynamicFilter) []sqlm.Expression {
	var filters []sqlm.Expression

	// we define default cases here
	var defaultExpression sqlm.Expression
	switch role, _ := filter["role"]; role {
	case "buyer":
		defaultExpression = sqlm.Exp("buyer_id", "=", sqlm.P(filter["buyer_id"]))
	case "seller":
		defaultExpression = sqlm.Exp("seller_id", "=", sqlm.P(filter["seller_id"]))
	default:
		defaultExpression = sqlm.Or(
			sqlm.Exp("buyer_id", "=", sqlm.P(filter["buyer_id"])),
			sqlm.Exp("seller_id", "=", sqlm.P(filter["seller_id"])),
		)
	}
	filters = append(filters, defaultExpression)

	for key, val := range filter {
		if _, ok := transactionAllowedFilters[key]; ok {
			filters = append(filters, sqlm.Exp(key, "=", sqlm.P(val)))
		}
	}

	return filters
}

func (m *mysqlTransaction) UpdateTransactionState(ctx context.Context, transactionID int64, transaction *entity.Transaction) error {
	now := time.Now()
	transaction.UpdatedAt = now

	query := `UPDATE transactions SET
		status = ?, invoice_id = ?, paid_at = ?, finished_at = ?, updated_at = ?
		WHERE id = ?`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error preparing update transaction query")
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx,
		transaction.Status, transaction.InvoiceID, transaction.PaidAt,
		transaction.FinishedAt, transaction.UpdatedAt, transactionID,
	)
	if err != nil {
		return errors.Wrap(err, "error executing update transaction query")
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return errors.New(fmt.Sprintf("Unexpected behavior detected when updating product (total rows affected: %d)", affectedRows))
	}

	return nil
}
