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

// AuthHandler holds AuthUsecase to be used in the handler
type AuthHandler struct {
	uc api.AuthUsecase
}

// NewAuthHandler creates a new instance of UserHandler
// with the provided UserUsecase
func NewAuthHandler(uc api.AuthUsecase) AuthHandler {
	return AuthHandler{uc}
}

// RegisterHandler registers all route for this handler
func (h *AuthHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	d := handler.DefaultMiddlewares()
	r.POST("/auth", handler.Decorate(h.Authenticate, d...))

	return nil
}

// Authenticate is a handler for user authentication
func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var auth entity.AuthCredentials
	if err := decoder.Decode(&auth); err != nil {
		err = api.ErrInvalidParameter
		api.Error(w, err)
		return err
	}

	ctx := r.Context()
	authResponse, err := h.uc.AuthenticateUser(ctx, &auth)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, authResponse, nil)
	return nil
}
