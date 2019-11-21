package entity_test

import (
	"testing"
	"time"

	"sejastip.id/api/fixture"

	"github.com/stretchr/testify/assert"

	"sejastip.id/api/entity"
)

func TestProductNormalizeCreate(t *testing.T) {
	product := entity.Product{
		ID:          1,
		Title:       "  test ",
		Description: "  barang mantap",
		Image:       "http://image.test/test.jpg  ",
	}

	product.NormalizeCreate()

	assert.Equal(t, product.Title, "test")
	assert.Equal(t, product.Description, "barang mantap")
	assert.Equal(t, product.Image, "http://image.test/test.jpg")
}

func TestProductValidateCreate(t *testing.T) {
	testCases := []struct {
		title       string
		product     entity.Product
		shouldError bool
	}{
		{
			title: "success",
			product: entity.Product{
				Title:     "barang mantap",
				Price:     100,
				SellerID:  1,
				ToDate:    time.Now().Add(10 * time.Hour),
				Image:     "http://image.test/product.jpg",
				CountryID: 1,
			},
			shouldError: false,
		},
		{
			title: "title less than 3 char",
			product: entity.Product{
				Title:     "",
				Price:     100,
				SellerID:  1,
				ToDate:    time.Now().Add(10 * time.Hour),
				Image:     "http://image.test/product.jpg",
				CountryID: 1,
			},
			shouldError: true,
		},
		{
			title: "price is less than 1",
			product: entity.Product{
				Title:     "barang mantap",
				Price:     0,
				SellerID:  1,
				ToDate:    time.Now().Add(10 * time.Hour),
				Image:     "http://image.test/product.jpg",
				CountryID: 1,
			},
			shouldError: true,
		},
		{
			title: "no seller",
			product: entity.Product{
				Title:     "barang mantap",
				Price:     100,
				SellerID:  0,
				ToDate:    time.Now().Add(10 * time.Hour),
				Image:     "http://image.test/product.jpg",
				CountryID: 1,
			},
			shouldError: true,
		},
		{
			title: "to date is already passed",
			product: entity.Product{
				Title:     "barang mantap",
				Price:     100,
				SellerID:  1,
				ToDate:    time.Now().Add(-10 * time.Hour),
				Image:     "http://image.test/product.jpg",
				CountryID: 1,
			},
			shouldError: true,
		},
		{
			title: "to date is already passed",
			product: entity.Product{
				Title:     "barang mantap",
				Price:     100,
				SellerID:  1,
				ToDate:    time.Now().Add(10 * time.Hour),
				Image:     "",
				CountryID: 1,
			},
			shouldError: true,
		},
		{
			title: "country not defined",
			product: entity.Product{
				Title:     "barang mantap",
				Price:     100,
				SellerID:  1,
				ToDate:    time.Now().Add(10 * time.Hour),
				Image:     "http://image.test/product.jpg",
				CountryID: 0,
			},
			shouldError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			err := test.product.ValidateCreate()
			assert.Equal(t, err != nil, test.shouldError)
		})
	}
}

func TestProductConvertToPublic(t *testing.T) {
	product := fixture.StubbedProduct()
	product.Status = entity.ProductStatusOffered
	user := fixture.StubbedUser()
	country := fixture.StubbedCountry()

	pp := product.ConvertToPublic(&country, &user)

	assert.Equal(t, pp.ID, product.ID)
	assert.Equal(t, pp.Status, "offered")
	assert.Equal(t, pp.Seller.ID, user.ID)
	assert.Equal(t, pp.Country.ID, country.ID)
}
