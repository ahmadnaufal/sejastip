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

type InvoiceHandler struct {
	invoiceUsecase api.InvoiceUsecase
}

func NewInvoiceHandler(uc api.InvoiceUsecase) InvoiceHandler {
	return InvoiceHandler{uc}
}

func (h *InvoiceHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.POST("/invoices", handler.Decorate(h.CreateInvoice, handler.UserAuth...))
	r.GET("/invoices/:id", handler.Decorate(h.GetInvoice, handler.UserAuth...))
	r.PATCH("/invoices/:id", handler.Decorate(h.UpdateInvoice, handler.UserAuth...))

	return nil
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var invoiceForm entity.InvoiceCreateForm
	if err := decoder.Decode(&invoiceForm); err != nil {
		api.Error(w, err)
		return err
	}

	ctx := r.Context()
	invoice, err := h.invoiceUsecase.InsertInvoice(ctx, &invoiceForm)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.Created(w, invoice, "")
	return nil
}

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	invoiceID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, api.ErrInvalidParameter)
		return err
	}

	// get invoice
	ctx := r.Context()
	invoice, err := h.invoiceUsecase.GetInvoice(ctx, invoiceID)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, invoice, "")
	return nil
}

func (h *InvoiceHandler) UpdateInvoice(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	invoiceID, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		api.Error(w, err)
		return err
	}

	decoder := json.NewDecoder(r.Body)
	var form entity.InvoiceUpdateForm
	if err := decoder.Decode(&form); err != nil {
		api.Error(w, err)
		return err
	}

	ctx := r.Context()
	invoice, err := h.invoiceUsecase.UpdateInvoice(ctx, invoiceID, &form)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, invoice, "Transaksi berhasil diperbarui")
	return nil
}
