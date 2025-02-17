package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/entity"
	"sejastip.id/api/handler"
	"sejastip.id/api/util"
)

type ProductHandler struct {
	uc api.ProductUsecase
}

func NewProductHandler(uc api.ProductUsecase) ProductHandler {
	return ProductHandler{uc}
}

func (h *ProductHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.GET("/products", handler.Decorate(h.GetProducts, handler.AppAuth...))
	r.GET("/products/:id", handler.Decorate(h.GetProduct, handler.AppAuth...))
	r.POST("/products", handler.Decorate(h.CreateProduct, handler.UserAuth...))
	r.PUT("/products/:id", handler.Decorate(h.UpdateProduct, handler.UserAuth...))
	r.DELETE("/products/:id", handler.Decorate(h.DeleteProduct, handler.UserAuth...))

	return nil
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	helper := api.NewQueryHelper(r)
	limit := helper.GetInt("limit", 10)
	offset := helper.GetInt("offset", 0)

	filters := helper.GetFilters()

	// get products by filter
	ctx := r.Context()
	banks, total, err := h.uc.GetProductsByFilter(ctx, filters, limit, offset)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting products")
	}

	meta := api.NewMetaPagination(http.StatusOK, limit, offset, int(total))
	api.OKWithMeta(w, banks, "", meta)
	return nil
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var productForm entity.ProductForm
	if err := decoder.Decode(&productForm); err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	// parse the date first
	fromDateInTime, err := time.Parse(time.RFC3339, productForm.FromDate+"T00:00:00Z")
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	toDateInTime, err := time.Parse(time.RFC3339, productForm.ToDate+"T23:59:59Z")
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	ctx := r.Context()
	filePath := ""
	if productForm.ImageFile != "" {
		file, extension, err := util.DecodeUploadedBase64File(productForm.ImageFile)
		if err != nil {
			err = api.ValidationError(fmt.Errorf("Error parsing file: %v", err))
			api.Error(w, err)
			return err
		}

		// upload file
		filename := fmt.Sprintf("%s%s", uuid.New().String(), extension)
		filePath, err = h.uc.UploadProductImage(ctx, filename, file)
		if err != nil {
			api.Error(w, err)
			return errors.Wrap(err, "error in uploading file")
		}
	}

	// Create product, but normalize the inputs first
	meta := api.MetaFromContext(ctx)
	product := entity.Product{
		Title:       productForm.Title,
		Description: productForm.Description,
		Price:       productForm.Price,
		SellerID:    meta.ID, // get the user ID from meta acquired from context
		CountryID:   productForm.CountryID,
		FromDate:    fromDateInTime,
		ToDate:      toDateInTime,
		Image:       filePath,
	}
	product.NormalizeCreate()

	productPublic, err := h.uc.CreateProduct(ctx, &product)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.Created(w, productPublic, "")
	return nil
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		err = api.ErrInvalidParameter
		api.Error(w, err)
		return err
	}

	// Get product by id
	ctx := r.Context()
	product, err := h.uc.GetProduct(ctx, id)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, product, "")
	return nil
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	productID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var product entity.Product
	if err := decoder.Decode(&product); err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	// update product
	ctx := r.Context()
	meta := api.MetaFromContext(ctx)
	productPublic, err := h.uc.UpdateProduct(ctx, productID, meta.ID, &product)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, productPublic, "product successfully updated")
	return nil
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	productID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		err = api.ErrInvalidParameter
		api.Error(w, err)
		return err
	}

	// delete product by id
	ctx := r.Context()
	meta := api.MetaFromContext(ctx)
	err = h.uc.DeleteProduct(ctx, productID, meta.ID)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, nil, "product has been successfully deleted")
	return nil
}
