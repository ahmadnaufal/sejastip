package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	return nil, 0, nil
}

func (m *mysqlProduct) GetProductsByFilter(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.Product, int64, error) {
	filteredQueries := buildDynamicQuery(filter)
	log.Println(filteredQueries, filter)
	countQuery, countArgs := sqlm.Build(
		"SELECT COUNT(id) FROM products",
		"WHERE", sqlm.And(filteredQueries),
	)
	log.Println(countQuery, countArgs)
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
	log.Println(query)
	results := []entity.Product{}
	err = m.db.SelectContext(ctx, &results, query, args...)
	return results, count, err
}

func (m *mysqlProduct) GetProduct(ctx context.Context, ID int64) (*entity.Product, error) {
	return nil, nil
}

func (m *mysqlProduct) UpdateProduct(ctx context.Context, ID int64, newProduct *entity.Product) error {
	return nil
}

func (m *mysqlProduct) DeleteProduct(ctx context.Context, ID int64) error {
	return nil
}

func buildDynamicQuery(filter entity.DynamicFilter) []interface{} {
	var filters []interface{}
	// to handle no filter
	filters = append(filters, sqlm.Exp("1=1"))
	for key, val := range filter {
		lowerKey := strings.ToLower(key)
		if _, ok := allowedFilters[lowerKey]; ok {
			if lowerKey == "q" {
				filters = append(filters, sqlm.Exp("title", "LIKE", sqlm.P(fmt.Sprintf("%%%s%%", val))))
			} else {
				filters = append(filters, sqlm.Exp(lowerKey, "=", sqlm.P(val)))
			}
		}
	}

	return filters
}
