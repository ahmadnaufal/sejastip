package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"sejastip.id/api"
)

type mysqlCountry struct {
	db *sqlx.DB
}

// NewMysqlCountry creates a new instance of Mysql country repository
func NewMysqlCountry(db *sql.DB) api.CountryRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlCountry{newDB}
}

// CreateCountry inserts a country data to repository
func (m *mysqlCountry) CreateCountry(ctx context.Context, country *api.Country) error {
	now := time.Now()
	country.CreatedAt = now
	country.UpdatedAt = now

	query := `INSERT INTO countries
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
		country.Name, country.Image, country.CreatedAt, country.UpdatedAt,
	)
	if err != nil {
		return err
	}

	country.ID, err = res.LastInsertId()
	return err
}

// GetCountries get all registered countries
func (m *mysqlCountry) GetCountries(ctx context.Context, limit, offset int) ([]api.Country, int64, error) {
	var count int64
	err := m.db.GetContext(ctx, &count, `SELECT COUNT(id) FROM countries`)
	if err != nil {
		return nil, 0, err
	}

	// prepare query, default is ordered by name
	query := `
		SELECT * FROM countries
		ORDER BY name ASC
		LIMIT ?, ?
	`
	results := []api.Country{}
	err = m.db.SelectContext(ctx, &results, query, offset, limit)
	return results, count, err
}

// GetCountry returns a country in mysql, by the ID
func (m *mysqlCountry) GetCountry(ctx context.Context, ID int64) (*api.Country, error) {
	query := `
		SELECT * FROM countries
		WHERE id = ?
		LIMIT 1
	`
	var result api.Country
	err := m.db.GetContext(ctx, &result, query, ID)
	if err == sql.ErrNoRows {
		return nil, api.ErrNotFound
	}
	return &result, err
}
