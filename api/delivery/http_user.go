package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sejastip.id/api"
	"sejastip.id/api/entity"
	"sejastip.id/api/handler"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// UserHandler holds UserUsecase to be used in the handler
type UserHandler struct {
	uc api.UserUsecase
}

// NewUserHandler creates a new instance of UserHandler
// with the provided UserUsecase
func NewUserHandler(uc api.UserUsecase) UserHandler {
	return UserHandler{uc}
}

// RegisterHandler registers all route for this handler
func (h *UserHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.POST("/users", handler.Decorate(h.Register, handler.AppAuth...))
	r.GET("/users/:id", handler.Decorate(h.GetUser, handler.AppAuth...))
	r.GET("/me", handler.Decorate(h.GetMe, handler.UserAuth...))

	return nil
}

// Register is a handler for user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var user entity.User
	if err := decoder.Decode(&user); err != nil {
		api.Error(w, err)
		return err
	}

	ctx := r.Context()
	publicUser, err := h.uc.Register(ctx, &user)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, publicUser, nil)
	return nil
}

// GetUser is a handler for get a single user data
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		err = api.ErrInvalidParameter
		api.Error(w, err)
		return err
	}

	// Get user by id
	ctx := r.Context()
	user, err := h.uc.GetUser(ctx, id)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, user, nil)
	return nil
}

// GetMe is a handler to get the logged in user's data
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	ctx := r.Context()
	claims := api.MetaFromContext(ctx)

	// get user from claim in context
	user, err := h.uc.GetUser(ctx, claims.ID)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, user, nil)
	return nil
}
