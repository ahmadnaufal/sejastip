package usecase

import (
	"context"
	"strings"

	"sejastip.id/api/storage"

	"github.com/pkg/errors"
	"sejastip.id/api"
)

type CountryProvider struct {
	CountryRepo api.CountryRepository
	Storage     storage.Storage
}

type countryUsecase struct {
	*CountryProvider
}

func NewCountryUsecase(pvd *CountryProvider) api.CountryUsecase {
	return &countryUsecase{pvd}
}

func (u *countryUsecase) CreateCountry(ctx context.Context, country *api.Country) error {
	err := u.CountryProvider.CountryRepo.CreateCountry(ctx, country)
	if err != nil {
		return errors.Wrap(err, "error in creating country")
	}

	return nil
}

func (u *countryUsecase) GetCountries(ctx context.Context, limit, offset int) ([]api.Country, int64, error) {
	countries, count, err := u.CountryProvider.CountryRepo.GetCountries(ctx, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error in getting all countries")
	}
	return countries, count, nil
}

func (u *countryUsecase) GetCountry(ctx context.Context, ID int64) (*api.Country, error) {
	country, err := u.CountryProvider.CountryRepo.GetCountry(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in getting a country with ID")
	}
	return country, nil
}

func (u *countryUsecase) UploadCountryImage(ctx context.Context, filename string, content []byte) (string, error) {
	return u.CountryProvider.Storage.Store("countries/"+strings.ToLower(filename), content)
}
