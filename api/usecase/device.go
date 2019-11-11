package usecase

import (
	"context"

	"github.com/pkg/errors"

	"sejastip.id/api"
	"sejastip.id/api/entity"
)

type DeviceProvider struct {
	DeviceRepo api.DeviceRepository
}

type deviceUsecase struct {
	*DeviceProvider
}

func NewDeviceUsecase(pvd *DeviceProvider) api.DeviceUsecase {
	return &deviceUsecase{pvd}
}

func (uc *deviceUsecase) UpsertDevice(ctx context.Context, device *entity.Device) error {
	meta := api.MetaFromContext(ctx)
	if meta.ID < 1 {
		return api.ErrForbidden
	}
	existingDevice, err := uc.DeviceRepo.GetUserDevice(ctx, meta.ID)
	if err != nil && err != api.ErrNotFound {
		return errors.Wrap(err, "error fetching device")
	}

	if existingDevice != nil {
		err = uc.DeviceRepo.UpdateUserDevice(ctx, existingDevice.ID, device)
		if err != nil {
			return errors.Wrap(err, "error updating existing device")
		}
	} else {
		device.UserID = meta.ID
		err = uc.DeviceRepo.InsertUserDevice(ctx, device)
		if err != nil {
			return errors.Wrap(err, "error inserting new device")
		}
	}

	return nil
}
