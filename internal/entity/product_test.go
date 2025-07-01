package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Test Product", 100.0)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 100.0, product.Price)
	assert.NotEmpty(t, product.CreatedAt)
}

func TestProductWhenNameIsRequired(t *testing.T) {
	product, err := NewProduct("", 100)

	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	product, err := NewProduct("Test Product", 0)

	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	product, err := NewProduct("Test Product", -100)

	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProduct_Validate(t *testing.T) {
	product, _ := NewProduct("Test Product", 100)

	err := product.Validate()

	assert.Nil(t, err)
}
