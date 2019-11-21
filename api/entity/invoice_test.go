package entity_test

import (
	"testing"

	"sejastip.id/api/entity"

	"github.com/stretchr/testify/assert"

	"sejastip.id/api/fixture"
)

func TestInvoiceConvertToPublic(t *testing.T) {
	invoice := fixture.StubbedInvoice()
	invoice.Status = entity.InvoiceStatusPending

	ip := invoice.ConvertToPublic()

	assert.Equal(t, ip.ID, invoice.ID)
	assert.Equal(t, ip.Status, "pending")
}

func TestInvoiceValidate(t *testing.T) {
	testCases := []struct {
		title       string
		invoice     entity.InvoiceCreateForm
		shouldError bool
	}{
		{
			title: "success",
			invoice: entity.InvoiceCreateForm{
				TransactionID: 1,
				PaymentMethod: "transfer",
			},
			shouldError: false,
		},
		{
			title: "transaction ID missing",
			invoice: entity.InvoiceCreateForm{
				TransactionID: 0,
				PaymentMethod: "transfer",
			},
			shouldError: true,
		},
		{
			title: "payment method invalid",
			invoice: entity.InvoiceCreateForm{
				TransactionID: 1,
				PaymentMethod: "t",
			},
			shouldError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			err := test.invoice.Validate()
			assert.Equal(t, err != nil, test.shouldError)
		})
	}
}
