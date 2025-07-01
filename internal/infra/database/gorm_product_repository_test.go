package database

import (
	"fmt"
	"math/rand"
	"testing"

	"os"
	"path/filepath"

	"github.com/Mazzael/go-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateNewProduct(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "file:memory:")

	db, err := gorm.Open(sqlite.Open(tmpFile), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 100)
	assert.Nil(t, err)

	gormProductRepository := NewProduct(db)
	err = gormProductRepository.Create(product)
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID)

	os.Remove(tmpFile)
}

func TestFindAllProducts(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "file:memory:")

	db, err := gorm.Open(sqlite.Open(tmpFile), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)

		db.Create(product)
	}

	gormProductRepository := NewProduct(db)

	products, err := gormProductRepository.FindAll(0, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = gormProductRepository.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = gormProductRepository.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)

	os.Remove(tmpFile)
}

func TestFindProductByID(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "file:memory:")

	db, err := gorm.Open(sqlite.Open(tmpFile), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10.00)
	assert.NoError(t, err)

	db.Create(product)

	gormProductRepository := NewProduct(db)

	productFound, err := gormProductRepository.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, product.Name, productFound.Name)

	os.Remove(tmpFile)
}

func TestUpdateProduct(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "file:memory:")

	db, err := gorm.Open(sqlite.Open(tmpFile), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10.00)
	assert.NoError(t, err)

	db.Create(product)

	gormProductRepository := NewProduct(db)

	product.Name = "Updated Product"
	err = gormProductRepository.Update(product)
	assert.NoError(t, err)

	productFound, err := gormProductRepository.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", productFound.Name)

	os.Remove(tmpFile)
}

func TestDeleteProduct(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "file:memory:")

	db, err := gorm.Open(sqlite.Open(tmpFile), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product1, err1 := entity.NewProduct("Product 1", 10.00)
	product2, err2 := entity.NewProduct("Product 2", 20.00)
	assert.NoError(t, err1)
	assert.NoError(t, err2)

	db.Create(product1)
	db.Create(product2)

	gormProductRepository := NewProduct(db)

	err = gormProductRepository.Delete(product1.ID.String())
	assert.NoError(t, err)

	productFound, err := gormProductRepository.FindByID(product1.ID.String())
	assert.Error(t, err)
	assert.Nil(t, productFound)

	os.Remove(tmpFile)
}
