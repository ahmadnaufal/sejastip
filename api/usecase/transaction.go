package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"sejastip.id/api"
	"sejastip.id/api/entity"
	"sejastip.id/api/infra"
)

type TransactionProvider struct {
	TransactionRepo api.TransactionRepository
	ShippingRepo    api.ShippingRepository
	UserRepo        api.UserRepository
	ProductRepo     api.ProductRepository
	AddressRepo     api.UserAddressRepository
	CountryRepo     api.CountryRepository
	DeviceRepo      api.DeviceRepository
	Pubsub          *infra.PubsubClient
}

type TransactionUsecase struct {
	*TransactionProvider
}

func NewTransactionUsecase(pvd *TransactionProvider) api.TransactionUsecase {
	return &TransactionUsecase{pvd}
}

// GetTransactions
func (uc *TransactionUsecase) GetTransactions(ctx context.Context, filter entity.DynamicFilter, limit, offset int) ([]*entity.TransactionPublic, int64, error) {
	transactions, total, err := uc.TransactionRepo.GetTransactions(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error in fetching transactions by filter")
	}

	transactionsPublic := []*entity.TransactionPublic{}
	for _, transaction := range transactions {
		transactionPublic, err := uc.convertToPublic(ctx, &transaction)
		if err != nil {
			return nil, total, err
		}

		transactionsPublic = append(transactionsPublic, transactionPublic)
	}

	return transactionsPublic, total, nil
}

// GetTransaction
func (uc *TransactionUsecase) GetTransaction(ctx context.Context, transactionID int64) (*entity.TransactionPublic, error) {
	transaction, err := uc.TransactionRepo.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return uc.convertToPublic(ctx, transaction)
}

func (uc *TransactionUsecase) convertToPublic(ctx context.Context, transaction *entity.Transaction) (*entity.TransactionPublic, error) {
	if transaction == nil {
		return nil, nil
	}

	// get buyer data
	buyer, err := uc.UserRepo.GetUser(ctx, transaction.BuyerID)
	if err != nil {
		return nil, err
	}

	// get buyer address data
	address, err := uc.AddressRepo.GetUserAddress(ctx, transaction.BuyerAddressID)
	if err != nil {
		return nil, err
	}

	// get product detail
	product, err := uc.ProductRepo.GetProduct(ctx, transaction.ProductID)
	if err != nil {
		return nil, err
	}

	// get seller user data
	seller, err := uc.UserRepo.GetUser(ctx, transaction.SellerID)
	if err != nil {
		return nil, err
	}

	// get country details
	country, err := uc.CountryRepo.GetCountry(ctx, product.CountryID)
	if err != nil {
		return nil, err
	}

	// get shipping
	shipping, err := uc.ShippingRepo.GetShipping(ctx, transaction.ID)
	if err != nil && err != api.ErrNotFound {
		return nil, err
	}

	productPublic := product.ConvertToPublic(country, seller)
	buyerPublic := buyer.ConvertToPublic()
	buyerAddressPublic := address.ConvertToPublic()

	var shippingPublic *entity.TransactionShippingPublic
	if shipping != nil {
		shippingPublic = &entity.TransactionShippingPublic{
			AWBNumber: shipping.AWBNumber,
			Courier:   shipping.Courier,
			CreatedAt: shipping.CreatedAt,
			UpdatedAt: shipping.UpdatedAt,
		}
	}

	// build the public transaction
	transactionPublic := &entity.TransactionPublic{
		ID:           transaction.ID,
		Product:      &productPublic,
		Buyer:        buyerPublic,
		BuyerAddress: &buyerAddressPublic,
		Quantity:     transaction.Quantity,
		Notes:        transaction.Notes,
		TotalPrice:   transaction.TotalPrice,
		Status:       transaction.GetStatusString(),
		Shipping:     shippingPublic,
		InvoiceID:    transaction.InvoiceID,
		PaidAt:       transaction.PaidAt,
		FinishedAt:   transaction.FinishedAt,
		CreatedAt:    transaction.CreatedAt,
		UpdatedAt:    transaction.UpdatedAt,
	}

	return transactionPublic, nil
}

// CreateTransaction
func (uc *TransactionUsecase) CreateTransaction(ctx context.Context, transactionForm *entity.TransactionForm, userID int64) (*entity.TransactionPublic, error) {
	// validate transaction form
	err := transactionForm.Validate()
	if err != nil {
		// our validation method will always return validation error
		// which is bad request
		return nil, api.ValidationError(err)
	}

	// then, do product validation next
	product, err := uc.ProductRepo.GetProduct(ctx, transactionForm.ProductID)
	if err != nil {
		return nil, errors.Wrap(err, "error in product checking")
	}
	// making sure the requesting user does not order his/her own product(s)
	if product.SellerID == userID {
		return nil, api.ErrBuyOwnProduct
	}

	// next, do address validation
	address, err := uc.AddressRepo.GetUserAddress(ctx, transactionForm.AddressID)
	if err != nil {
		return nil, errors.Wrap(err, "error in address checking")
	}

	// making sure the address is owned by the requesting user
	if address.UserID != userID {
		return nil, api.ErrTransactionAddressNotOwned
	}

	// finally, after lots of relational validations, we create our transaction object
	transaction := entity.Transaction{
		ProductID:      transactionForm.ProductID,
		BuyerID:        userID,
		SellerID:       product.SellerID,
		InvoiceID:      nil,
		BuyerAddressID: transactionForm.AddressID,
		Quantity:       transactionForm.Quantity,
		Notes:          transactionForm.Notes,
		TotalPrice:     int64(transactionForm.Quantity * product.Price),
	}
	err = uc.TransactionRepo.CreateTransaction(ctx, &transaction)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction")
	}

	// notify
	device, _ := uc.DeviceRepo.GetUserDevice(ctx, transaction.SellerID)
	if device != nil {
		if user, _ := uc.UserRepo.GetUser(ctx, transaction.SellerID); user != nil {
			notification := &entity.NotificationRequest{
				Device: device.DeviceID,
				UserID: transaction.SellerID,
			}
			notification.Data.Title = fmt.Sprintf("Hi %s, ada transaksi baru!", user.Name)
			notification.Data.Content = fmt.Sprintf("Ada yang ingin membeli %s dari kamu.", product.Title)
			uc.Pubsub.PublishNotification(ctx, notification)
		}
	}

	return uc.GetTransaction(ctx, transaction.ID)
}

func (uc *TransactionUsecase) UpdateTransaction(ctx context.Context, transactionID int64, form *entity.UpdateTransactionForm) error {
	transaction, err := uc.TransactionRepo.GetTransaction(ctx, transactionID)
	if err != nil {
		return errors.Wrap(err, "error fetching transaction")
	}

	// check transaction owner. reject any edit request from unauthorized users
	meta := api.MetaFromContext(ctx)
	userID := meta.ID
	// for now, even for the buyer
	if userID != transaction.SellerID {
		return api.ErrEditTransactionForbidden
	}

	if err := form.Validate(); err != nil {
		// our validation method will always return validation error
		// which is bad request
		return api.ValidationError(err)
	}

	// TODO need also validation on transaction state change

	// update the transaction status
	statusInt, ok := entity.MapStatusToStringReverse[strings.ToLower(form.Status)]
	if !ok {
		return api.ErrInvalidTransactionStateTransition
	}
	now := time.Now()
	transaction.Status = statusInt
	switch transaction.Status {
	case entity.TransactionStatusPaid:
		transaction.PaidAt = &now
	case entity.TransactionStatusFinished:
		transaction.FinishedAt = &now
	case entity.TransactionStatusDelivered:
		shipping := &entity.TransactionShipping{
			TransactionID: transactionID,
			AWBNumber:     form.AWBNumber,
			Courier:       form.Courier,
		}
		err = uc.ShippingRepo.InsertShipping(ctx, shipping)
		if err != nil {
			return errors.Wrap(err, "error inserting shipping info")
		}
	}
	err = uc.TransactionRepo.UpdateTransactionState(ctx, transactionID, transaction)
	if err != nil {
		return errors.Wrap(err, "error updating transaction state")
	}

	return nil
}
