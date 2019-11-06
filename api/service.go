package api

import (
	"context"

	"sejastip.id/api/entity"
)

// UserRepository is a contract for structs implementing user storage
type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUsers(ctx context.Context, limit, offset int) ([]entity.User, int64, error)
	GetUser(ctx context.Context, ID int64) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, ID int64, user *entity.User) error
}

// BankRepository is a contract for structs implementing banks storage
type BankRepository interface {
	CreateBank(ctx context.Context, bank *entity.Bank) error
	GetBanks(ctx context.Context, limit, offset int) ([]entity.Bank, int64, error)
	GetBankByName(ctx context.Context, name string) (*entity.Bank, error)
}

// CountryRepository is a contract for structs implementing country repository
type CountryRepository interface {
	CreateCountry(ctx context.Context, country *entity.Country) error
	GetCountries(ctx context.Context, limit, offset int) ([]entity.Country, int64, error)
	GetCountry(ctx context.Context, ID int64) (*entity.Country, error)
	BulkCreateCountries(ctx context.Context, countries []entity.Country) error
}

// ProductRepository is a contract for structs implementing product storage
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductsByUser(ctx context.Context, userID int64, limit, offset int) ([]entity.Product, int64, error)
	GetProductsByFilter(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.Product, int64, error)
	GetProduct(ctx context.Context, ID int64) (*entity.Product, error)
	UpdateProduct(ctx context.Context, ID int64, newProduct *entity.Product) error
	DeleteProduct(ctx context.Context, ID int64) error
}

// UserAddressRepository is a contract for structs implementing user address storage
type UserAddressRepository interface {
	CreateAddress(ctx context.Context, address *entity.UserAddress) error
	GetUserAddresses(ctx context.Context, userID int64, limit, offset int) ([]entity.UserAddress, int64, error)
	GetUserAddress(ctx context.Context, ID int64) (*entity.UserAddress, error)
	UpdateAddress(ctx context.Context, ID int64, newAddress *entity.UserAddress) error
}

// TransactionRepository is a contract for structs implementing transaction storage
type TransactionRepository interface {
	GetTransactions(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.Transaction, int64, error)
	GetTransaction(ctx context.Context, transactionID int64) (*entity.Transaction, error)
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) error
}

// UserUsecase is a contract for usecases related to users
type UserUsecase interface {
	Register(ctx context.Context, user *entity.User) (*entity.UserPublic, error)
	GetUser(ctx context.Context, ID int64) (*entity.UserPublic, error)
}

// AuthUsecase is a contract for usecase related to authentication
type AuthUsecase interface {
	AuthenticateUser(ctx context.Context, auth *entity.AuthCredentials) (*entity.AuthResponse, error)
}

// BankUsecase is a contract for usecase related to bank data
type BankUsecase interface {
	CreateBank(ctx context.Context, bank *entity.Bank) error
	GetBanks(ctx context.Context, limit, offset int) ([]entity.Bank, int64, error)
	UploadBankImage(ctx context.Context, filename string, content []byte) (string, error)
}

// CountryUsecase is a contract for usecase related to countries
type CountryUsecase interface {
	CreateCountry(ctx context.Context, country *entity.Country) error
	GetCountries(ctx context.Context, limit, offset int) ([]entity.Country, int64, error)
	GetCountry(ctx context.Context, ID int64) (*entity.Country, error)
	UploadCountryImage(ctx context.Context, filename string, content []byte) (string, error)
	BulkCreateCountries(ctx context.Context, countries []entity.Country) error
}

// ProductUsecase is a contract for structs implementing product usecase
type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.ProductPublic, error)
	GetProductsByFilter(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.ProductPublic, int64, error)
	GetProduct(ctx context.Context, ID int64) (*entity.ProductPublic, error)
	UpdateProduct(ctx context.Context, productID, userID int64, newProduct *entity.Product) (*entity.ProductPublic, error)
	DeleteProduct(ctx context.Context, productID, userID int64) error
	UploadProductImage(ctx context.Context, filename string, content []byte) (string, error)
}

// UserAddressUsecase is a contract for structs implementing user address usecase
type UserAddressUsecase interface {
	CreateAddress(ctx context.Context, address *entity.UserAddress) (*entity.UserAddressPublic, error)
	GetUserAddresses(ctx context.Context, userID int64, limit, offset int) ([]entity.UserAddressPublic, int64, error)
	GetUserAddress(ctx context.Context, ID int64) (*entity.UserAddressPublic, error)
	UpdateAddress(ctx context.Context, ID int64, newAddress *entity.UserAddress) (*entity.UserAddressPublic, error)
}

// TransactionUsecase is a contract for structs implementing transactions usecase
type TransactionUsecase interface {
	GetTransactions(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]*entity.TransactionPublic, int64, error)
	GetTransaction(ctx context.Context, transactionID int64) (*entity.TransactionPublic, error)
	CreateTransaction(ctx context.Context, transactionForm *entity.TransactionForm, userID int64) (*entity.TransactionPublic, error)
}
