package delivery

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/handler"
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

	return nil
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	helper := api.NewQueryHelper(r)
	limit := helper.GetInt("limit", 10)
	offset := helper.GetInt("offset", 0)

	filters := helper.GetFilters()

	// Get all banks
	ctx := r.Context()
	banks, total, err := h.uc.GetProductsByFilter(ctx, filters, limit, offset)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting products")
	}

	meta := api.NewPaginationMeta(limit, offset, int(total))
	api.OK(w, banks, meta)
	return nil
}
