package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct(t *testing.T) {
	p, err := NewProduct("Product 1", 10)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, "Product 1", p.Name)
	assert.Equal(t, 10, p.Price)
}

func TestProduct_WhenNameIsRequired(t *testing.T) {
	p, err := NewProduct("", 10)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProduct_WhenPriceIsRequired(t *testing.T) {
	p, err := NewProduct("Product 1", 0)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProduct_WhenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("Product 1", -1)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProduct_Validate(t *testing.T) {
	p, err := NewProduct("Product 1", 10)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	err = p.Validate()
	assert.Nil(t, err)
}
