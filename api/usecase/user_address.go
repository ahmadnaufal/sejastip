package usecase

import (
	"context"

	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

type UserAddressProvider struct {
	UserAddressRepo api.UserAddressRepository
}

type userAddressUsecase struct {
	*UserAddressProvider
}

func NewUserAddressUsecase(pvd *UserAddressProvider) api.UserAddressUsecase {
	return &userAddressUsecase{pvd}
}

func (uc *userAddressUsecase) CreateAddress(ctx context.Context, address *entity.UserAddress) (*entity.UserAddressPublic, error) {
	err := uc.UserAddressRepo.CreateAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	addressPublic := address.ConvertToPublic()
	return &addressPublic, nil
}

func (uc *userAddressUsecase) GetUserAddresses(ctx context.Context, userID int64, limit, offset int) ([]entity.UserAddressPublic, int64, error) {
	addresses, count, err := uc.UserAddressRepo.GetUserAddresses(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	addressesPublic := []entity.UserAddressPublic{}
	for _, address := range addresses {
		addressesPublic = append(addressesPublic, address.ConvertToPublic())
	}

	return addressesPublic, count, nil
}

func (uc *userAddressUsecase) GetUserAddress(ctx context.Context, ID int64) (*entity.UserAddressPublic, error) {
	address, err := uc.UserAddressRepo.GetUserAddress(ctx, ID)
	if err != nil {
		return nil, errors.Wrap(err, "error in fetching address")
	}

	addressPublic := address.ConvertToPublic()
	return &addressPublic, nil
}

func (uc *userAddressUsecase) UpdateAddress(ctx context.Context, ID int64, newAddress *entity.UserAddress) (*entity.UserAddressPublic, error) {
	err := uc.UserAddressRepo.UpdateAddress(ctx, ID, newAddress)
	if err != nil {
		return nil, errors.Wrap(err, "error in updating product")
	}

	return uc.GetUserAddress(ctx, ID)
}
