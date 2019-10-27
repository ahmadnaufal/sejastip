package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

type mysqlUserAddress struct {
	db *sqlx.DB
}

// NewMysqlUserAddress creates a new instance of MySQL UserAddress repository
func NewMysqlUserAddress(db *sql.DB) api.UserAddressRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlUserAddress{newDB}
}

// CreateAddress inserts a newly registered address data into our mysql repository
func (m *mysqlUserAddress) CreateAddress(ctx context.Context, address *entity.UserAddress) error {
	now := time.Now()
	address.CreatedAt = now
	address.UpdatedAt = now

	query := `INSERT INTO user_addresses
		(address, phone, address_name, user_id,
		created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?)
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	// execute query
	res, err := prep.ExecContext(ctx,
		address.Address, address.Phone, address.AddressName, address.UserID,
		address.CreatedAt, address.UpdatedAt,
	)
	if err != nil {
		return err
	}

	address.ID, err = res.LastInsertId()
	return err
}

// GetUserAddresses fetches registered user address data from our mysql repository
func (m *mysqlUserAddress) GetUserAddresses(ctx context.Context, userID int64, limit, offset int) ([]entity.UserAddress, int64, error) {
	var count int64
	err := m.db.GetContext(ctx, &count, `SELECT COUNT(id) FROM user_addresses WHERE user_id=?`, userID)
	if err != nil {
		return nil, 0, err
	}

	// prepare query
	query := `
		SELECT * FROM user_addresses
		WHERE user_id = ?
		ORDER BY updated_at DESC
		LIMIT ?, ?
	`
	results := []entity.UserAddress{}
	err = m.db.SelectContext(ctx, &results, query, userID, offset, limit)
	return results, count, err
}

// GetUserAddress fetch a user address by its ID
func (m *mysqlUserAddress) GetUserAddress(ctx context.Context, ID int64) (*entity.UserAddress, error) {
	query := `
		SELECT * FROM user_addresses
		WHERE id = ?
	`
	result := &entity.UserAddress{}
	err := m.db.GetContext(ctx, result, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}
	}
	return result, err
}

// UpdateAddress updates a user address row data, selected by its ID with new provided data
func (m *mysqlUserAddress) UpdateAddress(ctx context.Context, ID int64, address *entity.UserAddress) error {
	now := time.Now()
	address.UpdatedAt = now

	query := `
		UPDATE user_addresses SET
		address = ?, phone = ?, address_name = ?, updated_at = ?
		WHERE id = ?
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx,
		address.Address, address.Phone, address.AddressName, address.UpdatedAt, ID,
	)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return api.SejastipError{
			Message: fmt.Sprintf("Unexpected behavior detected when updating address data (total rows affected: %d)", affectedRows),
		}
	}

	return nil
}
