package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"sejastip.id/api/entity"
	"sejastip.id/api/util"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/handler"
)

type CountryHandler struct {
	uc api.CountryUsecase
}

func NewCountryHandler(uc api.CountryUsecase) CountryHandler {
	return CountryHandler{uc}
}

func (h *CountryHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.POST("/countries", handler.Decorate(h.CreateCountry, handler.AppAuth...))
	r.GET("/countries", handler.Decorate(h.GetCountries, handler.AppAuth...))
	r.GET("/countries/:id", handler.Decorate(h.GetCountry, handler.AppAuth...))

	return nil
}

func (h *CountryHandler) CreateCountry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var countryForm entity.CountryForm
	if err := decoder.Decode(&countryForm); err != nil {
		api.Error(w, err)
		return err
	}

	file, extension, err := util.DecodeUploadedBase64File(countryForm.ImageFile)
	if err != nil {
		err = api.SejastipError{
			Message:    fmt.Sprintf("Error parsing file: %v", err),
			ErrorCode:  400,
			HTTPStatus: http.StatusBadRequest,
		}
		api.Error(w, err)
		return err
	}

	ctx := r.Context()

	// upload file
	filename := fmt.Sprintf("%s%s", slug.Make(countryForm.Name), extension)
	filePath, err := h.uc.UploadCountryImage(ctx, filename, file)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error in uploading file")
	}

	country := entity.Country{
		Name:  countryForm.Name,
		Image: filePath,
	}
	err = h.uc.CreateCountry(ctx, &country)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.Created(w, country, "")
	return nil
}

func (h *CountryHandler) GetCountries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	helper := api.NewQueryHelper(r)
	limit := helper.GetInt("limit", 10)
	offset := helper.GetInt("offset", 0)

	// Get all countries
	ctx := r.Context()
	countries, total, err := h.uc.GetCountries(ctx, limit, offset)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting countries")
	}

	meta := api.NewMetaPagination(http.StatusOK, limit, offset, int(total))
	api.OKWithMeta(w, countries, "", meta)
	return nil
}

func (h *CountryHandler) GetCountry(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	countryID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil || countryID < 1 {
		err = api.ErrInvalidParameter
		api.Error(w, err)
		return err
	}

	// Get all countries
	ctx := r.Context()
	country, err := h.uc.GetCountry(ctx, countryID)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting a country")
	}

	api.OK(w, country, "")
	return nil
}
