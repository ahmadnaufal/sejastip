package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/entity"
	"sejastip.id/api/handler"
)

type UserAddressHandler struct {
	uc api.UserAddressUsecase
}

func NewUserAddressHandler(uc api.UserAddressUsecase) UserAddressHandler {
	return UserAddressHandler{uc}
}

func (h *UserAddressHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.POST("/users/:id/addresses", handler.Decorate(h.CreateAddressForUser, handler.UserAuth...))
	r.GET("/users/:id/addresses", handler.Decorate(h.GetUserAddresses, handler.AppAuth...))
	r.GET("/addresses/:id", handler.Decorate(h.GetUserAddress, handler.AppAuth...))
	r.PUT("/addresses/:id", handler.Decorate(h.UpdateAddress, handler.UserAuth...))

	return nil
}

func (h *UserAddressHandler) CreateAddressForUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var form entity.UserAddressForm
	if err := decoder.Decode(&form); err != nil {
		api.Error(w, err)
		return err
	}

	userID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}
	ctx := r.Context()
	if userID != api.MetaFromContext(ctx).ID {
		err = api.ErrForbidden
		api.Error(w, err)
		return err
	}

	address := entity.UserAddress{
		Address:     form.Address,
		Phone:       form.Phone,
		AddressName: form.AddressName,
		UserID:      userID,
	}
	addressPublic, err := h.uc.CreateAddress(ctx, &address)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.Created(w, addressPublic, "")
	return nil
}

func (h *UserAddressHandler) GetUserAddresses(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	helper := api.NewQueryHelper(r)
	limit := helper.GetInt("limit", 10)
	offset := helper.GetInt("offset", 0)

	userID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	// Get all banks
	ctx := r.Context()
	addresses, total, err := h.uc.GetUserAddresses(ctx, userID, limit, offset)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting user addresses")
	}

	meta := api.NewMetaPagination(http.StatusOK, limit, offset, int(total))
	api.OKWithMeta(w, addresses, "", meta)
	return nil
}

func (h *UserAddressHandler) GetUserAddress(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	// Get all banks
	ctx := r.Context()
	address, err := h.uc.GetUserAddress(ctx, id)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting user address")
	}

	api.OK(w, address, "")
	return nil
}

func (h *UserAddressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var form entity.UserAddressForm
	if err := decoder.Decode(&form); err != nil {
		api.Error(w, err)
		return err
	}

	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	ctx := r.Context()

	address := entity.UserAddress{
		Address:     form.Address,
		Phone:       form.Phone,
		AddressName: form.AddressName,
	}
	addressPublic, err := h.uc.UpdateAddress(ctx, id, &address)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, addressPublic, "")
	return nil
}
