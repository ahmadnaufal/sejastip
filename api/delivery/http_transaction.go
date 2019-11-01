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

type TransactionHandler struct {
	transactionUsecase api.TransactionUsecase
}

func NewTransactionHandler(uc api.TransactionUsecase) TransactionHandler {
	return TransactionHandler{uc}
}

func (h *TransactionHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.POST("/transactions", handler.Decorate(h.CreateTransaction, handler.UserAuth...))
	r.GET("/transactions", handler.Decorate(h.GetTransactions, handler.UserAuth...))
	r.GET("/transactions/:id", handler.Decorate(h.GetTransaction, handler.AppAuth...))

	return nil
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var transactionForm entity.TransactionForm
	if err := decoder.Decode(&transactionForm); err != nil {
		api.Error(w, err)
		return err
	}

	ctx := r.Context()
	transaction, err := h.transactionUsecase.CreateTransaction(ctx, &transactionForm)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.Created(w, transaction, "")
	return nil
}

func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	helper := api.NewQueryHelper(r)
	limit := helper.GetInt("limit", 10)
	offset := helper.GetInt("offset", 0)

	// get filters
	filters := helper.GetFilters()

	ctx := r.Context()
	transactions, total, err := h.transactionUsecase.GetTransactions(ctx, filters, limit, offset)
	if err != nil {
		api.Error(w, err)
		return err
	}

	meta := api.NewMetaPagination(http.StatusOK, limit, offset, int(total))
	api.OKWithMeta(w, transactions, "", meta)
	return nil
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	transactionID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, err)
		return err
	}

	ctx := r.Context()
	transaction, err := h.transactionUsecase.GetTransaction(ctx, transactionID)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, transaction, "")
	return nil
}
