package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shuoli84/sqlm"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

type mysqlProduct struct {
	db *sqlx.DB
}

var allowedFilters = map[string]struct{}{
	"q":          struct{}{},
	"seller_id":  struct{}{},
	"country_id": struct{}{},
}

func NewMysqlProduct(db *sql.DB) api.ProductRepository {
	newDB := sqlx.NewDb(db, "mysql")
	return &mysqlProduct{newDB}
}

func (m *mysqlProduct) CreateProduct(ctx context.Context, product *entity.Product) error {
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	query := `INSERT INTO products
		(title, description, price, seller_id, country_id, image, status,
		from_date, to_date, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	// execute query
	res, err := prep.ExecContext(ctx,
		product.Title, product.Description, product.Price, product.SellerID,
		product.CountryID, product.Image, product.Status, product.FromDate,
		product.ToDate, product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return err
	}

	product.ID, err = res.LastInsertId()
	return err
}

func (m *mysqlProduct) GetProductsByUser(ctx context.Context, userID int64, limit, offset int) ([]entity.Product, int64, error) {
	var count int64
	err := m.db.GetContext(ctx, &count, `SELECT COUNT(id) FROM products WHERE seller_id=?`, userID)
	if err != nil {
		return nil, 0, err
	}

	// prepare query
	query := `
		SELECT * FROM products
		WHERE seller_id=? AND deleted_at IS NULL
		ORDER BY updated_at DESC
		LIMIT ?, ?
	`
	results := []entity.Product{}
	err = m.db.SelectContext(ctx, &results, query, userID, offset, limit)
	return results, count, err
}

func (m *mysqlProduct) GetProductsByFilter(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.Product, int64, error) {
	filteredQueries := buildDynamicQuery(filter)
	countQuery, countArgs := sqlm.Build(
		"SELECT COUNT(id) FROM products",
		"WHERE", sqlm.And(filteredQueries),
	)

	var count int64
	err := m.db.GetContext(ctx, &count, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	query, args := sqlm.Build(
		"SELECT * FROM products",
		"WHERE", sqlm.And(filteredQueries),
		"ORDER BY updated_at DESC",
		sqlm.Exp("LIMIT", sqlm.P(offset), ",", sqlm.P(limit)),
	)
	results := []entity.Product{}
	err = m.db.SelectContext(ctx, &results, query, args...)
	return results, count, err
}

func (m *mysqlProduct) GetProduct(ctx context.Context, ID int64) (*entity.Product, error) {
	query := `
		SELECT * FROM products
		WHERE id = ?
	`
	result := &entity.Product{}
	err := m.db.GetContext(ctx, result, query, ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, api.ErrNotFound
		}
	}
	return result, err
}

func (m *mysqlProduct) UpdateProduct(ctx context.Context, ID int64, newProduct *entity.Product) error {
	now := time.Now()
	newProduct.UpdatedAt = now

	query := `
		UPDATE products SET
		title = ?, description = ?, price = ?, country_id = ?, status = ?,
		from_date = ?, to_date = ?, updated_at = ?
		WHERE id = ?
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx,
		newProduct.Title, newProduct.Description, newProduct.Price,
		newProduct.CountryID, newProduct.Status, newProduct.FromDate,
		newProduct.ToDate, newProduct.UpdatedAt,
		ID,
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
			Message: fmt.Sprintf("Unexpected behavior detected when updating product (total rows affected: %d)", affectedRows),
		}
	}

	return nil
}

func (m *mysqlProduct) DeleteProduct(ctx context.Context, ID int64) error {
	now := time.Now()
	query := `
		UPDATE products SET
		deleted_at = ?
		WHERE id = ?
	`
	prep, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := prep.ExecContext(ctx, now, ID)
	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		return api.SejastipError{
			Message: fmt.Sprintf("Unexpected behavior detected when deleting product (total rows affected: %d)", affectedRows),
		}
	}

	return nil
}

func buildDynamicQuery(filter entity.DynamicFilter) []interface{} {
	var filters []interface{}
	// to handle no filter
	filters = append(filters, sqlm.Exp("deleted_at IS NULL"))
	for key, val := range filter {
		if _, ok := allowedFilters[key]; ok {
			if key == "q" {
				filters = append(filters, sqlm.Exp("title", "LIKE", sqlm.P(fmt.Sprintf("%%%s%%", val))))
			} else {
				filters = append(filters, sqlm.Exp(key, "=", sqlm.P(val)))
			}
		}
	}

	return filters
}
