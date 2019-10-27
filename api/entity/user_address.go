package entity

import "time"

// UserAddress stores database row representations of a user address data
type UserAddress struct {
	ID          int64     `db:"id"`
	Address     string    `db:"address"`
	Phone       string    `db:"phone"`
	AddressName string    `db:"address_name"`
	UserID      int64     `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// UserAddressPublic is the public representation of UserAddress
type UserAddressPublic struct {
	ID          int64     `json:"id"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	AddressName string    `json:"address_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *UserAddress) ConvertToPublic() UserAddressPublic {
	return UserAddressPublic{
		ID:          u.ID,
		Address:     u.Address,
		Phone:       u.Phone,
		AddressName: u.AddressName,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type UserAddressForm struct {
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	AddressName string `json:"address_name"`
}
