package delivery

import (
	"encoding/json"
	"net/http"

	"sejastip.id/api"
	"sejastip.id/api/entity"
	"sejastip.id/api/handler"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type DeviceHandler struct {
	uc api.DeviceUsecase
}

func NewDeviceHandler(uc api.DeviceUsecase) DeviceHandler {
	return DeviceHandler{uc}
}

func (h *DeviceHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.PUT("/devices", handler.Decorate(h.UpsertDevice, handler.UserAuth...))

	return nil
}

func (h *DeviceHandler) UpsertDevice(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	userAgent := api.GetUserAgent(r)
	platform := api.GetPlatform(userAgent)

	decoder := json.NewDecoder(r.Body)
	var form entity.DeviceForm
	if err := decoder.Decode(&form); err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	device := entity.Device{
		DeviceID:  form.DeviceID,
		UserAgent: userAgent,
		Platform:  platform,
	}

	ctx := r.Context()
	err := h.uc.UpsertDevice(ctx, &device)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, nil, "User device registered successfully")
	return nil
}
