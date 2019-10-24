package entity

import "time"

// Country represents country data in database
type Country struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Image     string    `json:"image" db:"image"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// CountryForm
type CountryForm struct {
	Name      string `json:"name"`
	ImageFile string `json:"image_file"`
}
