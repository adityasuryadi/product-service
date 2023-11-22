package repository

import "product-service/entity"

type ProductRepository interface {
	CreateProduct(entity *entity.Product) (*entity.Product, error)
	UpdateProduct(entity *entity.Product) (*entity.Product, error)
	DeletProduct(id string) error
	GetProductById(id string) (product entity.Product, err error)
	ListProduct() (products []entity.Product, err error)
}
