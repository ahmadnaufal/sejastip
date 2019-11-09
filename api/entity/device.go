package entity

import (
	"time"
)

type Device struct {
	ID        int64     `db:"id"`
	DeviceID  string    `db:"device_id"`
	Platform  string    `db:"platform"`
	UserAgent string    `db:"user_agent"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type DeviceForm struct {
	DeviceID string `json:"device_id"`
}
