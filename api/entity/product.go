package entity

import "time"

const (
	ProductStatusIdle = iota
	ProductStatusOffered
	ProductStatusOutOfStock
)

var mapProductStatusToString = map[uint]string{
	ProductStatusIdle:       "idle",
	ProductStatusOffered:    "offered",
	ProductStatusOutOfStock: "out of stock",
}

// Product stores database row representations of a product data
type Product struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Price       uint       `db:"price"`
	SellerID    int64      `db:"seller_id"`
	CountryID   int64      `db:"country_id"`
	Image       string     `db:"image"`
	Status      uint       `db:"status"`
	FromDate    time.Time  `db:"from_date"`
	ToDate      time.Time  `db:"to_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (p *Product) ConvertToPublic(c *Country, u *User) ProductPublic {
	return ProductPublic{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Seller:      u.ConvertToPublic(),
		Country:     c,
		Status:      mapProductStatusToString[p.Status],
		FromDate:    p.FromDate.Format("2006-01-02"),
		ToDate:      p.ToDate.Format("2006-01-02"),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type ProductPublic struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Price       uint        `json:"price"`
	Seller      *UserPublic `json:"seller,omitempty"`
	Country     *Country    `json:"country,omitempty"`
	Status      string      `json:"status"`
	FromDate    string      `json:"from_date"`
	ToDate      string      `json:"to_date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
