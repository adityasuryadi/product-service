package repository

import (
	"product-service/entity"

	"gorm.io/gorm"
)

func NewProductRepository(database *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		db: database,
	}
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

// DeletProduct implements ProductRepository.
func (repository *ProductRepositoryImpl) DeleteProduct(id string) error {
	var entity entity.Product
	err := repository.db.Where("id = ?", id).Delete(&entity).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateProduct implements ProductRepository.
func (repository *ProductRepositoryImpl) UpdateProduct(entity *entity.Product) (*entity.Product, error) {
	err := repository.db.Save(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// GetProductById implements ProductRepository.
func (repository *ProductRepositoryImpl) GetProductById(id string) (product entity.Product, err error) {
	err = repository.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

// CreateProduct implements ProductRepository.
func (repository *ProductRepositoryImpl) CreateProduct(entity *entity.Product) (*entity.Product, error) {
	err := repository.db.Create(entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// ListProduct implements ProductRepository.
func (repository *ProductRepositoryImpl) ListProduct() (products []entity.Product, err error) {
	var entity []entity.Product
	result := repository.db.Find(&entity).Debug()
	if result.Error != nil {
		return entity, result.Error
	}
	return entity, nil
}
