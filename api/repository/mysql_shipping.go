package repository

import (
	"context"
	"database/sql"
	"time"

	"sejastip.id/api"
	"sejastip.id/api/entity"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlShipping struct {
	db *sqlx.DB
}

// NewMysqlShipping creates a new instance of MySQL shipping repository
func NewMysqlShipping(db *sql.DB) api.ShippingRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlShipping{newDB}
}

func (m *mysqlShipping) InsertShipping(ctx context.Context, shipping *entity.TransactionShipping) error {
	now := time.Now()
	shipping.CreatedAt = now
	shipping.UpdatedAt = now

	query := `INSERT INTO transaction_shippings
		(transaction_id, awb_number, courier, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?)`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error preparing insert shipping query")
	}
	defer prep.Close()

	res, err := prep.ExecContext(ctx,
		shipping.TransactionID, shipping.AWBNumber, shipping.Courier,
		shipping.CreatedAt, shipping.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "error executing insert shipping query")
	}

	shipping.ID, err = res.LastInsertId()
	return err
}

func (m *mysqlShipping) GetShipping(ctx context.Context, transactionID int64) (*entity.TransactionShipping, error) {
	query := `
		SELECT * FROM transaction_shippings
		WHERE transaction_id = ?
		LIMIT 1
	`
	result := &entity.TransactionShipping{}
	err := m.db.GetContext(ctx, result, query, transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}

		return nil, err
	}

	return result, nil
}
