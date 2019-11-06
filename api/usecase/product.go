package usecase

import (
	"context"
	"strings"

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
	if err := product.ValidateCreate(); err != nil {
		return nil, api.ValidationError(err)
	}
	_, err := uc.Provider.CountryRepo.GetCountry(ctx, product.CountryID)
	if err != nil {
		return nil, err
	}

	err = uc.Provider.ProductRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, errors.Wrap(err, "error in creating product")
	}

	user, country, _ := uc.fetchProductAdditionalInfo(ctx, *product)

	productPublic := product.ConvertToPublic(country, user)
	return &productPublic, nil
}

func (uc *productUsecase) GetProductsByFilter(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]entity.ProductPublic, int64, error) {
	products, count, err := uc.Provider.ProductRepo.GetProductsByFilter(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error in fetching products by filter")
	}

	publicProducts := []entity.ProductPublic{}
	for _, product := range products {
		user, country, err := uc.fetchProductAdditionalInfo(ctx, product)
		if err != nil {
			return nil, count, err
		}

		publicProducts = append(publicProducts, product.ConvertToPublic(country, user))
	}

	return publicProducts, count, nil
}

func (uc *productUsecase) GetProduct(ctx context.Context, ID int64) (*entity.ProductPublic, error) {
	product, err := uc.Provider.ProductRepo.GetProduct(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in fetching product")
	}

	user, country, err := uc.fetchProductAdditionalInfo(ctx, *product)
	if err != nil {
		return nil, err
	}

	publicProduct := product.ConvertToPublic(country, user)
	return &publicProduct, nil
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, productID, userID int64, newProduct *entity.Product) (*entity.ProductPublic, error) {
	// check first if the product is owned by the user
	product, err := uc.Provider.ProductRepo.GetProduct(ctx, productID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching product")
	}

	if product.SellerID != userID {
		return nil, api.ErrEditProductForbidden
	}

	err = uc.Provider.ProductRepo.UpdateProduct(ctx, productID, newProduct)
	if err != nil {
		return nil, errors.Wrap(err, "error in updating product")
	}

	return uc.GetProduct(ctx, productID)
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, productID, userID int64) error {
	// check first if the product is owned by the user
	product, err := uc.Provider.ProductRepo.GetProduct(ctx, productID)
	if err != nil {
		return errors.Wrap(err, "error fetching product")
	}

	if product.SellerID != userID {
		return api.ErrEditProductForbidden
	}

	// the user requesting product deletion is the owner of the product:
	// proceed with the delete
	err = uc.Provider.ProductRepo.DeleteProduct(ctx, productID)
	if err != nil {
		return errors.Wrap(err, "error in deleting product")
	}

	return nil
}

func (u *productUsecase) UploadProductImage(ctx context.Context, filename string, content []byte) (string, error) {
	return u.Provider.Storage.Store("products/"+strings.ToLower(filename), content)
}

func (uc *productUsecase) fetchProductAdditionalInfo(ctx context.Context, product entity.Product) (*entity.User, *entity.Country, error) {
	user, err := uc.Provider.UserRepo.GetUser(ctx, product.SellerID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error in fetching user details from product")
	}

	country, err := uc.Provider.CountryRepo.GetCountry(ctx, product.CountryID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error in fetching country details from product")
	}

	return user, country, nil
}
