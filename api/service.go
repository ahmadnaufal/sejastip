package api

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// User stores database row representations of a user data
type User struct {
	ID          int64      `db:"id"`
	Email       string     `json:"email" db:"email"`
	Name        string     `json:"name" db:"name"`
	Phone       string     `json:"phone" db:"phone"`
	Password    string     `json:"password" db:"password"`
	BankName    string     `json:"bank_name" db:"bank_name"`
	BankAccount string     `json:"bank_account" db:"bank_account"`
	Avatar      string     `db:"avatar"`
	LastLoginAt *time.Time `db:"last_login_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// Normalize is a method to normalize all field values
func (u *User) Normalize() {
	u.Email = strings.TrimSpace(u.Email)
	u.Name = strings.TrimSpace(u.Name)
	// removes all non-number char, including + sign
	u.Phone = regexp.MustCompile(`\D`).ReplaceAllString(strings.TrimSpace(u.Phone), "")
	r := regexp.MustCompile("^0+")
	if r.MatchString(u.Phone) {
		u.Phone = r.ReplaceAllString(u.Phone, "")
		u.Phone = fmt.Sprintf("62%s", u.Phone)
	}

	u.Password = strings.TrimSpace(u.Password)
}

// Validate is a function to validate user input validity
func (u *User) Validate() error {
	// check name first: must have no special character
	matched, err := regexp.Match("^[A-Za-z0-9\\s]+$", []byte(u.Name))
	if err != nil {
		return err
	}
	if !matched {
		return SejastipError{
			Message: "Nama hanya boleh mengandung karakter alfanumerik",
		}
	}

	// check phone
	return nil
}

// ConvertToPublic converts the User model to public representations
func (u *User) ConvertToPublic() *UserPublic {
	return &UserPublic{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Phone:        u.Phone,
		BankName:     u.BankName,
		BankAccount:  u.BankAccount,
		RegisteredAt: u.CreatedAt,
		Avatar:       u.Avatar,
	}
}

// UserPublic is the collection of user data publicly available
type UserPublic struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	BankName     string    `json:"bank_name"`
	BankAccount  string    `json:"bank_account"`
	RegisteredAt time.Time `json:"registered_at"`
	Avatar       string    `json:"avatar"`
}

// AuthCredentials stores user authentication credentials
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Normalize is a method to normalize all field values
func (a *AuthCredentials) Normalize() {
	a.Email = strings.TrimSpace(a.Email)
	a.Password = strings.TrimSpace(a.Password)
}

// AuthResponse is our structure for token response after user authentication
type AuthResponse struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// UserRepository is a contract for structs implementing user storage
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context, limit, offset int) ([]User, int64, error)
	GetUser(ctx context.Context, ID int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, ID int64, user *User) error
}

// UserUsecase is a contract for usecases related to users
type UserUsecase interface {
	Register(ctx context.Context, user *User) (*UserPublic, error)
	GetUser(ctx context.Context, ID int64) (*UserPublic, error)
}

// AuthUsecase is a contract for usecase related to authentication
type AuthUsecase interface {
	AuthenticateUser(ctx context.Context, auth *AuthCredentials) (*AuthResponse, error)
}
