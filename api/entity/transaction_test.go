package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sejastip.id/api/entity"
)

func TestTransactionGetStatusString(t *testing.T) {
	txPlaced := entity.Transaction{Status: entity.TransactionStatusInit}
	txPaid := entity.Transaction{Status: entity.TransactionStatusPaid}
	txDelivered := entity.Transaction{Status: entity.TransactionStatusDelivered}
	txFinished := entity.Transaction{Status: entity.TransactionStatusFinished}

	assert.Equal(t, txPlaced.GetStatusString(), "placed")
	assert.Equal(t, txPaid.GetStatusString(), "paid")
	assert.Equal(t, txDelivered.GetStatusString(), "delivered")
	assert.Equal(t, txFinished.GetStatusString(), "finished")
}

func TestTransactionFormValidate(t *testing.T) {
	testCases := []struct {
		title       string
		form        entity.TransactionForm
		shouldError bool
	}{
		{
			title: "success",
			form: entity.TransactionForm{
				ProductID: 1,
				Quantity:  1,
				AddressID: 2,
			},
			shouldError: false,
		},
		{
			title: "no product ID",
			form: entity.TransactionForm{
				Quantity:  1,
				AddressID: 2,
			},
			shouldError: true,
		},
		{
			title: "quantity is zero",
			form: entity.TransactionForm{
				ProductID: 1,
				Quantity:  0,
				AddressID: 2,
			},
			shouldError: true,
		},
		{
			title: "no address selected",
			form: entity.TransactionForm{
				ProductID: 1,
				Quantity:  1,
			},
			shouldError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			err := test.form.Validate()
			assert.Equal(t, err != nil, test.shouldError)
		})
	}
}

func TestUpdateTransactionFormValidate(t *testing.T) {
	testCases := []struct {
		title       string
		form        entity.UpdateTransactionForm
		shouldError bool
	}{
		{
			title: "success",
			form: entity.UpdateTransactionForm{
				Status: "paid",
			},
			shouldError: false,
		},
		{
			title: "status undefined",
			form: entity.UpdateTransactionForm{
				Status: "invalid",
			},
			shouldError: true,
		},
		{
			title: "status delivered success",
			form: entity.UpdateTransactionForm{
				Status:    "delivered",
				AWBNumber: "1010",
				Courier:   "Test Kurir",
			},
			shouldError: false,
		},
		{
			title: "status delivered but no awb",
			form: entity.UpdateTransactionForm{
				Status:  "delivered",
				Courier: "Test Kurir",
			},
			shouldError: true,
		},
		{
			title: "status delivered but no courier",
			form: entity.UpdateTransactionForm{
				Status:    "delivered",
				AWBNumber: "1010",
			},
			shouldError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			err := test.form.Validate()
			assert.Equal(t, err != nil, test.shouldError)
		})
	}
}
