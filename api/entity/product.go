package entity

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

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

func (p *Product) NormalizeCreate() {
	p.Title = strings.TrimSpace(p.Title)
	p.Description = strings.TrimSpace(p.Description)
	p.Image = strings.TrimSpace(p.Image)
}

func (p *Product) ValidateCreate() error {
	if len(p.Title) < 3 {
		return errors.New("Judul produk harus lebih dari 3 karakter")
	}

	if p.Price < 1 {
		return errors.New("Harga produk tidak boleh kosong atau negatif")
	}

	if p.SellerID < 1 {
		return errors.New("Penjual harus terdaftar")
	}

	if p.ToDate.Before(time.Now()) {
		return errors.New("Waktu akhir penawaran barang tidak boleh di waktu yang lalu")
	}

	if p.Image == "" {
		return errors.New("Harus upload foto produk terlebih dahulu")
	}

	if p.CountryID < 1 {
		return errors.New("Harus memilih negara lokasi penjualan barang")
	}

	return nil
}

type ProductForm struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	CountryID   int64  `json:"country_id"`
	ImageFile   string `json:"image_file"`
	FromDate    string `json:"from_date"`
	ToDate      string `json:"to_date"`
}

func (p *Product) ConvertToPublic(c *Country, u *User) ProductPublic {
	return ProductPublic{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
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
	Image       string      `json:"image"`
	Seller      *UserPublic `json:"seller,omitempty"`
	Country     *Country    `json:"country,omitempty"`
	Status      string      `json:"status"`
	FromDate    string      `json:"from_date"`
	ToDate      string      `json:"to_date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
