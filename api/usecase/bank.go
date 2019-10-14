package usecase

import (
	"context"

	"sejastip.id/api/storage"

	"github.com/pkg/errors"
	"sejastip.id/api"
)

type BankProvider struct {
	BankRepo api.BankRepository
	Storage  storage.Storage
}

type bankUsecase struct {
	*BankProvider
}

func NewBankUsecase(pvd *BankProvider) api.BankUsecase {
	return &bankUsecase{pvd}
}

func (u *bankUsecase) CreateBank(ctx context.Context, bank *api.Bank) error {
	err := u.BankProvider.BankRepo.CreateBank(ctx, bank)
	if err != nil {
		return errors.Wrap(err, "error in creating bank")
	}

	return nil
}

func (u *bankUsecase) GetBanks(ctx context.Context, limit, offset int) ([]api.Bank, int64, error) {
	banks, count, err := u.BankProvider.BankRepo.GetBanks(ctx, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error in getting all banks")
	}
	return banks, count, nil
}

func (u *bankUsecase) UploadBankImage(ctx context.Context, filename string, content []byte) (string, error) {
	return u.BankProvider.Storage.Store("banks/"+filename, content)
}
