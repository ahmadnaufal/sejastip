package usecase

import (
	"context"

	"github.com/pkg/errors"

	"sejastip.id/api"
	"sejastip.id/api/entity"
	"sejastip.id/api/storage"
)

type ProductProvider struct {
	ProductRepo api.ProductRepository
	UserRepo    api.UserRepository
	CountryRepo api.CountryRepository

	Storage storage.Storage
}

type productUsecase struct {
	Provider *ProductProvider
}

func NewProductUsecase(pvd *ProductProvider) api.ProductUsecase {
	return &productUsecase{pvd}
}

func (uc *productUsecase) CreateProduct(ctx context.Context, product *entity.Product) (*entity.ProductPublic, error) {
	err := uc.Provider.ProductRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, errors.Wrap(err, "error in creating product")
	}

	productPublic := product.ConvertToPublic(nil, nil)
	return &productPublic, nil
}

func (uc *productUsecase) GetProductsByFilter(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.ProductPublic, int64, error) {
	products, count, err := uc.Provider.ProductRepo.GetProductsByFilter(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error in fetching products by filter")
	}

	publicProducts := []entity.ProductPublic{}
	for _, product := range products {
		user, err := uc.Provider.UserRepo.GetUser(ctx, product.SellerID)
		if err != nil {
			return nil, 0, errors.Wrap(err, "error in fetching user details from product")
		}

		country, err := uc.Provider.CountryRepo.GetCountry(ctx, product.CountryID)
		if err != nil {
			return nil, 0, errors.Wrap(err, "error in fetching country details from product")
		}

		publicProducts = append(publicProducts, product.ConvertToPublic(country, user))
	}

	return publicProducts, count, nil
}

func (uc *productUsecase) GetProduct(ctx context.Context, ID int64) (*entity.ProductPublic, error) {
	return nil, nil
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, ID int64, newProduct *entity.Product) (*entity.ProductPublic, error) {
	return nil, nil
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, ID int64) error {
	return nil
}
