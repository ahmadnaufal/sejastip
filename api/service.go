package api

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

// ResourceClaims for claims
type ResourceClaims struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
	jwt.StandardClaims
}

// AuthResponse is our structure for token response after user authentication
type AuthResponse struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

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

// Product stores database row representations of a product data
type Product struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Price       uint      `db:"price"`
	SellerID    int64     `db:"seller_id"`
	CountryID   int64     `db:"country_id"`
	Status      uint      `db:"status"`
	FromDate    time.Time `db:"from_date"`
	ToDate      time.Time `db:"to_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (p *Product) ConvertToPublic(c *Country, u *User) *ProductPublic {
	return &ProductPublic{
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
	Seller      *UserPublic `json:"seller"`
	Country     *Country    `json:"country"`
	Status      string      `json:"status"`
	FromDate    string      `json:"from_date"`
	ToDate      string      `json:"to_date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

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

// UserRepository is a contract for structs implementing user storage
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUsers(ctx context.Context, limit, offset int) ([]User, int64, error)
	GetUser(ctx context.Context, ID int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, ID int64, user *User) error
}

// BankRepository is a contract for structs implementing banks storage
type BankRepository interface {
	CreateBank(ctx context.Context, bank *Bank) error
	GetBanks(ctx context.Context, limit, offset int) ([]Bank, int64, error)
	GetBankByName(ctx context.Context, name string) (*Bank, error)
}

// CountryRepository is a contract for structs implementing country repository
type CountryRepository interface {
	CreateCountry(ctx context.Context, country *Country) error
	GetCountries(ctx context.Context, limit, offset int) ([]Country, int64, error)
	GetCountry(ctx context.Context, ID int64) (*Country, error)
}

// ProductRepository is a contract for structs implementing product storage
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProductsByKeyword(ctx context.Context, keyword string, limit, offset int) ([]Product, error)
	GetProductsByUser(ctx context.Context, userID int64, limit, offset int) ([]Product, error)
	GetProductsByCountry(ctx context.Context, countryID int64, limit, offset int) ([]Product, error)
	GetProduct(ctx context.Context, ID int64) (*Product, error)
	UpdateProduct(ctx context.Context, ID int64, newProduct *Product) error
	DeleteProduct(ctx context.Context, ID int64) error
}

// UserAddressRepository is a contract for structs implementing user address storage
type UserAddressRepository interface {
	CreateAddress(ctx context.Context, address *UserAddress) error
	GetUserAddresses(ctx context.Context, userID int64, limit, offset int) ([]UserAddress, error)
	GetUserAddress(ctx context.Context, ID int64) (*UserAddress, error)
	UpdateAddress(ctx context.Context, ID int64, newAddress *UserAddress) error
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

// BankUsecase is a contract for usecase related to bank data
type BankUsecase interface {
	CreateBank(ctx context.Context, bank *Bank) error
	GetBanks(ctx context.Context, limit, offset int) ([]Bank, int64, error)
	UploadBankImage(ctx context.Context, filename string, content []byte) (string, error)
}

// CountryUsecase is a contract for usecase related to countries
type CountryUsecase interface {
	CreateCountry(ctx context.Context, country *Country) error
	GetCountries(ctx context.Context, limit, offset int) ([]Country, int64, error)
	GetCountry(ctx context.Context, ID int64) (*Country, error)
	UploadCountryImage(ctx context.Context, filename string, content []byte) (string, error)
}

// ProductUsecase is a contract for structs implementing product usecase
type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *Product) (*ProductPublic, error)
	GetProductsByFilter(ctx context.Context, filters map[string]string, limit, offset int) ([]ProductPublic, int64, error)
	GetProduct(ctx context.Context, ID int64) (*ProductPublic, error)
	UpdateProduct(ctx context.Context, ID int64, newProduct *Product) (*ProductPublic, error)
	DeleteProduct(ctx context.Context, ID int64) error
}

// UserAddressUsecase is a contract for structs implementing user address usecase
type UserAddressUsecase interface {
	CreateAddress(ctx context.Context, address *UserAddress) (*UserAddressPublic, error)
	GetUserAddresses(ctx context.Context, userID int64, limit, offset int) ([]UserAddressPublic, int64, error)
	GetUserAddress(ctx context.Context, ID int64) (*UserAddressPublic, error)
	UpdateAddress(ctx context.Context, ID int64, newAddress *UserAddress) (*UserAddressPublic, error)
}
