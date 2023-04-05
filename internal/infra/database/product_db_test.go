package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10)
	productDB := NewProduct(db)

	// Act
	err = productDB.Create(product)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, product.ID)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, 10, product.Price)
	assert.NotEmpty(t, product.CreatedAt)
}

func TestProduct_FindAll(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	db.AutoMigrate(&entity.Product{})

	for i := 1; i < 28; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Intn(100)+1)
		assert.Nil(t, err)
		db.Create(product)
	}

	productDB := NewProduct(db)

	products, err := productDB.FindAll(1, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 7)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 27", products[6].Name)

	products, err = productDB.FindAll(3, 10, "")
	assert.Nil(t, err)
	assert.Len(t, products, 7)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 27", products[6].Name)

	products, err = productDB.FindAll(3, 9, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 7)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 27", products[6].Name)

	products, err = productDB.FindAll(3, 10, "desc")
	assert.Nil(t, err)
	assert.Len(t, products, 7)
	assert.Equal(t, "Product 7", products[0].Name)
	assert.Equal(t, "Product 1", products[6].Name)

	products, err = productDB.FindAll(3, 12, "desc")
	assert.Nil(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 3", products[0].Name)
	assert.Equal(t, "Product 1", products[2].Name)
}

func TestProduct_FindByID(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10)
	productDB := NewProduct(db)
	productDB.Create(product)

	// Act
	productFound, err := productDB.FindByID(product.ID.String())

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, "Product 1", productFound.Name)
	assert.Equal(t, 10, productFound.Price)
	assert.NotEmpty(t, productFound.CreatedAt)
}

func TestProduct_Update(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10)
	productDB := NewProduct(db)
	productDB.Create(product)

	// Act
	product.Name = "Product 2"
	product.Price = 20
	err = productDB.Update(product)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "Product 2", product.Name)
	assert.Equal(t, 20, product.Price)
}

func TestProduct_Delete(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	product, _ := entity.NewProduct("Product 1", 10)
	productDB := NewProduct(db)
	productDB.Create(product)

	// Act
	err = productDB.Delete(product.ID.String())

	// Assert
	assert.Nil(t, err)
}
