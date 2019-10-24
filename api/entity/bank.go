package entity

import "time"

// Bank stores database row representations of a bank data
type Bank struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Image     string    `json:"image" db:"image"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// BankForm
type BankForm struct {
	Name      string `json:"name"`
	ImageFile string `json:"image_file"`
}
