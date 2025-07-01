package database

import (
	"github.com/Mazzael/go-api/internal/entity"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{DB: db}
}

func (p *GormProductRepository) Create(product *entity.Product) error {
	return p.DB.Create(product).Error
}

func (p *GormProductRepository) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	if err := p.DB.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *GormProductRepository) Update(product *entity.Product) error {
	_, err := p.FindByID(product.ID.String())
	if err != nil {
		return err
	}

	return p.DB.Save(product).Error
}

func (p *GormProductRepository) Delete(id string) error {
	if err := p.DB.Delete(&entity.Product{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (p *GormProductRepository) FindAll(page, limit int, sort string) ([]*entity.Product, error) {
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	var products []*entity.Product
	offset := page * limit

	query := p.DB.Offset(offset).Limit(limit)

	query = query.Order("created_at " + sort)

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
