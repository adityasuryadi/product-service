package service

import (
	"product-service/model"
	"product-service/repository"
)

type ProductServiceImpl struct {
	Repository repository.ProductRepository
}

// FindProductById implements ProductService.
func (service *ProductServiceImpl) FindProductById(id string) (errorCode int, product model.ProductsResponse) {
	productEntity, err := service.Repository.GetProductById(id)
	if err != nil {
		return 404, product
	}

	product = model.ProductsResponse{
		Id:          productEntity.Id.String(),
		Name:        productEntity.Name,
		Price:       productEntity.Price,
		Qty:         productEntity.Qty,
		Description: productEntity.Description,
	}
	return 200, product
}

// GetAllProduct implements ProductService.
func (service *ProductServiceImpl) GetAllProduct() (errorCode int, products []model.ProductsResponse) {
	result, err := service.Repository.ListProduct()
	if err != nil {
		return 500, nil
	}
	for _, v := range result {
		products = append(products, model.ProductsResponse{
			Id:          v.Id.String(),
			Name:        v.Name,
			Qty:         v.Qty,
			Price:       v.Price,
			Description: v.Description,
		})
	}

	return 200, products
}

// InsertProduct implements ProductService.
func (service *ProductServiceImpl) InsertProduct(request model.CreateProductRequest) (errorCode int) {
	panic("unimplemented")
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		Repository: repository,
	}
}
