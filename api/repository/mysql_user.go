package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"sejastip.id/api"
)

type mysqlUser struct {
	db *sqlx.DB
}

// NewMysqlUser creates a new instance of MySQL User repository
func NewMysqlUser(db *sql.DB) api.UserRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlUser{newDB}
}

// CreateUser inserts a newly registered user data into our mysql repository
func (m *mysqlUser) CreateUser(ctx context.Context, user *api.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	log.Println(user)

	query := `INSERT INTO users
		(email, name, phone, password, bank_name, bank_account,
		created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	// execute query
	res, err := prep.ExecContext(ctx,
		user.Email, user.Name, user.Phone, user.Password, user.BankName,
		user.BankAccount, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	user.ID, err = res.LastInsertId()
	return err
}

// GetUsers fetches registered users' data from our mysql repository
func (m *mysqlUser) GetUsers(ctx context.Context, limit, offset int) ([]api.User, int64, error) {
	var count int64
	err := m.db.GetContext(ctx, &count, `SELECT COUNT(id) FROM users`)
	if err != nil {
		return nil, 0, err
	}

	// prepare query
	query := `
		SELECT * FROM users
		ORDER BY updated_at DESC
		LIMIT ?, ?
	`
	results := []api.User{}
	err = m.db.SelectContext(ctx, &results, query, offset, limit)
	return results, count, err
}

// GetUser fetch a user data by its ID
func (m *mysqlUser) GetUser(ctx context.Context, ID int64) (*api.User, error) {
	query := `
		SELECT * FROM users
		WHERE id = ?
	`
	result := &api.User{}
	err := m.db.GetContext(ctx, result, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}
	}
	return result, err
}

// UpdateUser updates a user's data row selected by its ID with new provided data
func (m *mysqlUser) UpdateUser(ctx context.Context, ID int64, user *api.User) error {
	now := time.Now()
	user.UpdatedAt = now

	query := `
		UPDATE users SET
		email = ?, name = ?, phone = ?, bank_name = ?, bank_account = ?,
		updated_at = ?
		WHERE id = ?
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx,
		user.Email, user.Name, user.Phone,
		user.BankName, user.BankAccount,
		user.UpdatedAt,
		user.ID,
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
			Message: fmt.Sprintf("Unexpected behavior detected when updating user (total rows affected: %d)", affectedRows),
		}
	}

	return nil
}

// GetUserByEmail fetches a user having the provided email
func (m *mysqlUser) GetUserByEmail(ctx context.Context, email string) (*api.User, error) {
	query := `
		SELECT * FROM users
		WHERE email = ?
		LIMIT 1
	`
	var result api.User
	err := m.db.GetContext(ctx, &result, query, email)
	if err == sql.ErrNoRows {
		return nil, api.ErrNotFound
	}
	return &result, err
}
