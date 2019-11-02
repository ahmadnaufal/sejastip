package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

type mysqlBank struct {
	db *sqlx.DB
}

// NewMysqlBank creates a new instance of Mysql bank repository
func NewMysqlBank(db *sql.DB) api.BankRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlBank{newDB}
}

// CreateBank inserts a bank data to repository
func (m *mysqlBank) CreateBank(ctx context.Context, bank *entity.Bank) error {
	now := time.Now()
	bank.CreatedAt = now
	bank.UpdatedAt = now

	query := `INSERT INTO banks
		(name, image, created_at, updated_at)
		VALUES
		(?, ?, ?, ?)
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	// execute query
	res, err := prep.ExecContext(ctx,
		bank.Name, bank.Image, bank.CreatedAt, bank.UpdatedAt,
	)
	if err != nil {
		return err
	}

	bank.ID, err = res.LastInsertId()
	return err
}

// GetBanks get all registered banks
func (m *mysqlBank) GetBanks(ctx context.Context, limit, offset int) ([]entity.Bank, int64, error) {
	var count int64
	err := m.db.GetContext(ctx, &count, `SELECT COUNT(id) FROM banks`)
	if err != nil {
		return nil, 0, err
	}

	// prepare query, default is ordered by name
	query := `
		SELECT * FROM banks
		ORDER BY name ASC
		LIMIT ?, ?
	`
	results := []entity.Bank{}
	err = m.db.SelectContext(ctx, &results, query, offset, limit)
	return results, count, err
}

// GetBankByName returns a bank in mysql, by the bank name
func (m *mysqlBank) GetBankByName(ctx context.Context, name string) (*entity.Bank, error) {
	query := `
		SELECT * FROM banks
		WHERE name = ?
		LIMIT 1
	`
	var result entity.Bank
	err := m.db.GetContext(ctx, &result, query, name)
	if err == sql.ErrNoRows {
		return nil, api.ErrNotFound
	}
	return &result, err
}
