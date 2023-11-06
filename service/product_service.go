package service

import "product-service/model"

type ProductService interface {
	GetAllProduct() (errorCode int, products []model.ProductsResponse)
	FindProductById(id string) (errorCode int, product model.ProductsResponse)
	InsertProduct(request model.CreateProductRequest) (errorCode int)
}
