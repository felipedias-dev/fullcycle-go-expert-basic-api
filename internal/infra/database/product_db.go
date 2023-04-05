package database

import (
	"github.com/felipedias-dev/fullcycle-go-expert-basic-api/internal/entity"
	"gorm.io/gorm"
)

type Product struct {
	DB *gorm.DB
}

func (p *Product) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *Product) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	err := p.DB.Order(sort).Offset((page - 1) * limit).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *Product) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := p.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *Product) Update(product *entity.Product) error {
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return p.DB.Save(product).Error
}

func (p *Product) Delete(id string) error {
	_, err := p.FindByID(id)
	if err != nil {
		return err
	}
	return p.DB.Where("id = ?", id).Delete(&entity.Product{}).Error
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{DB: db}
}
