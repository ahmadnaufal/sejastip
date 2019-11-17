package usecase

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sejastip.id/api/storage"
	"sejastip.id/api/util"

	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/entity"
)

type InvoiceProvider struct {
	InvoiceRepo     api.InvoiceRepository
	TransactionRepo api.TransactionRepository
	UserRepo        api.UserRepository

	Storage storage.Storage
}

type InvoiceUsecase struct {
	*InvoiceProvider
}

func NewInvoiceUsecase(pvd *InvoiceProvider) api.InvoiceUsecase {
	return &InvoiceUsecase{pvd}
}

func (uc *InvoiceUsecase) InsertInvoice(ctx context.Context, form *entity.InvoiceCreateForm) (*entity.InvoicePublic, error) {
	err := form.Validate()
	if err != nil {
		// our validation method will always return validation error
		// which is bad request
		return nil, api.SejastipError{
			Message:    err.Error(),
			ErrorCode:  400,
			HTTPStatus: http.StatusBadRequest,
		}
	}

	// check if transaction ID is valid
	transaction, err := uc.TransactionRepo.GetTransaction(ctx, form.TransactionID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching transaction")
	}

	// check transaction ownership: must be the buyer to create invoice
	if transaction.BuyerID != api.GetUserID(ctx) {
		return nil, api.ErrEditInvoiceForbidden
	}

	// then check if the invoice already had invoice
	existingInvoice, err := uc.InvoiceRepo.GetInvoiceFromTransaction(ctx, transaction.ID)
	if err != nil && err != api.ErrNotFound {
		return nil, errors.Wrap(err, "error fetching existing invoice")
	}
	if existingInvoice != nil {
		return nil, api.SejastipError{
			Message:    "Transaksi sudah memiliki invoice",
			ErrorCode:  422,
			HTTPStatus: http.StatusUnprocessableEntity,
		}
	}

	// else, create new invoice
	invoiceCode := entity.InvoiceNumber(fmt.Sprintf("JSTP%s%s", time.Now().Format("20061231"), strconv.FormatInt(transaction.ID, 36)))
	invoice := &entity.Invoice{
		TransactionID: transaction.ID,
		InvoiceCode:   invoiceCode,
		CodedPrice:    transaction.TotalPrice,
		PaymentMethod: form.PaymentMethod,
		Status:        entity.InvoiceStatusPending,
		PaidAt:        nil,
		ReceiptProof:  "",
	}
	err = uc.InvoiceRepo.InsertInvoice(ctx, invoice)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting invoice")
	}

	// update transaction to include invoice ID
	transaction.InvoiceID = &invoice.ID
	err = uc.TransactionRepo.UpdateTransactionState(ctx, transaction.ID, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "error updating transaction")
	}

	invoicePublic := invoice.ConvertToPublic()
	return &invoicePublic, nil
}

func (uc *InvoiceUsecase) GetInvoice(ctx context.Context, invoiceID int64) (*entity.InvoicePublic, error) {
	invoice, err := uc.InvoiceRepo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching invoice")
	}

	invoicePublic := invoice.ConvertToPublic()
	return &invoicePublic, nil
}

func (uc *InvoiceUsecase) UpdateInvoice(ctx context.Context, invoiceID int64, form *entity.InvoiceUpdateForm) (*entity.InvoicePublic, error) {
	invoice, err := uc.InvoiceRepo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching invoice")
	}

	// get transaction data to verify attributes related to transaction
	transaction, err := uc.TransactionRepo.GetTransaction(ctx, invoice.TransactionID)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching transaction")
	}

	// check transaction ownership: must be the buyer to create invoice
	if transaction.BuyerID != api.GetUserID(ctx) {
		return nil, api.ErrEditInvoiceForbidden
	}

	if form.ReceiptProof != "" {
		file, extension, err := util.DecodeUploadedBase64File(form.ReceiptProof)
		if err != nil {
			return nil, api.SejastipError{
				Message:    fmt.Sprintf("Error parsing file: %v", err),
				ErrorCode:  400,
				HTTPStatus: http.StatusBadRequest,
			}
		}

		// upload file
		filename := fmt.Sprintf("%s%s", invoice.InvoiceCode, extension)
		invoice.ReceiptProof, err = uc.uploadReceiptProof(ctx, filename, file)
		if err != nil {
			return nil, errors.Wrap(err, "error uploading payment proof")
		}

		err = uc.InvoiceRepo.UpdateInvoice(ctx, invoiceID, invoice)
		if err != nil {
			return nil, errors.Wrap(err, "error updating invoice")
		}
	}

	if form.Status == "paid" {
		now := time.Now()
		invoice.Status = entity.InvoiceStatusPaid
		invoice.PaidAt = &now
		err = uc.InvoiceRepo.UpdateInvoice(ctx, invoiceID, invoice)
		if err != nil {
			return nil, errors.Wrap(err, "error updating invoice")
		}

		// update transaction to paid
		transaction.Status = entity.TransactionStatusPaid
		transaction.PaidAt = &now
		err = uc.TransactionRepo.UpdateTransactionState(ctx, transaction.ID, transaction)
		if err != nil {
			return nil, errors.Wrap(err, "error updating transaction")
		}
	}

	invoicePublic := invoice.ConvertToPublic()
	return &invoicePublic, nil
}

func (uc *InvoiceUsecase) uploadReceiptProof(ctx context.Context, filename string, content []byte) (string, error) {
	return uc.Storage.Store("invoice_proofs/"+strings.ToLower(filename), content)
}
