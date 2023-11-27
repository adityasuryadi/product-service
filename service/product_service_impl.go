package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"product-service/entity"
	"product-service/exception"
	"product-service/model"
	"product-service/pkg/rabbitmq"
	"product-service/repository"

	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ProductServiceImpl struct {
	Repository   repository.ProductRepository
	RabbitMqConn *amqp091.Connection
}

func NewProductService(repository repository.ProductRepository, rabbitMQConn *amqp091.Connection) ProductService {
	return &ProductServiceImpl{
		Repository:   repository,
		RabbitMqConn: rabbitMQConn,
	}
}

// DeleteProduct implements ProductService.
func (service *ProductServiceImpl) DeleteProduct(id string) (errorCode int) {
	_, err := service.Repository.GetProductById(id)
	fmt.Println("product", err)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404
	}

	err = service.Repository.DeleteProduct(id)
	if err != nil {
		return 500
	}

	exchange := "product.deleted"
	queue := "product.delete"
	routingKey := "delete"

	PublisherConfig := &rabbitmq.PublisherConfig{
		Exchange:    exchange,
		QueueName:   queue,
		RoutingKey:  routingKey,
		ConsumerTag: "",
	}
	body, err := json.Marshal(id)
	if err != nil {
		exception.FailOnError(err, "failed decode product")
	}
	ch := rabbitmq.NewPublisher(service.RabbitMqConn, PublisherConfig)
	defer ch.CloseChannel()
	ch.SetupExchangeAndQueue()
	ch.Publish(body)
	return 200
}

// EditProduct implements ProductService.
func (service *ProductServiceImpl) EditProduct(id string, request model.CreateProductRequest) (errorCode int, product *model.ProductsResponse) {
	exchange := "product.updated"
	queue := "product.update"
	routingKey := "update"

	entityProduct, err := service.Repository.GetProductById(id)
	if err != nil {
		return 500, nil
	}

	if (entityProduct == entity.Product{}) {
		return 404, nil
	}

	entityProduct.Name = request.Name
	entityProduct.Price = request.Price
	entityProduct.Description = request.Description
	entityProduct.Qty = request.Qty

	result, err := service.Repository.UpdateProduct(&entityProduct)
	if err != nil {
		return 500, nil
	}

	product = &model.ProductsResponse{
		Id:          result.Id.String(),
		Name:        result.Name,
		Price:       result.Price,
		Qty:         result.Qty,
		Description: result.Description,
	}

	PublisherConfig := &rabbitmq.PublisherConfig{
		Exchange:    exchange,
		QueueName:   queue,
		RoutingKey:  routingKey,
		ConsumerTag: "",
	}
	body, err := json.Marshal(product)
	if err != nil {
		exception.FailOnError(err, "failed decode product")
	}
	ch := rabbitmq.NewPublisher(service.RabbitMqConn, PublisherConfig)
	defer ch.CloseChannel()
	ch.SetupExchangeAndQueue()
	ch.Publish(body)

	return 200, product
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
	entity := &entity.Product{
		Name:        request.Name,
		Price:       request.Price,
		Qty:         request.Qty,
		Description: request.Description,
	}

	product, err := service.Repository.CreateProduct(entity)
	if err != nil {
		exception.FailOnError(err, "Failed Create Product")
	}

	exchange := "product.created"
	queue := "product.create"
	routingKey := "create"

	PublisherConfig := &rabbitmq.PublisherConfig{
		Exchange:    exchange,
		QueueName:   queue,
		RoutingKey:  routingKey,
		ConsumerTag: "",
	}

	body, err := json.Marshal(product)
	if err != nil {
		exception.FailOnError(err, "failed decode product")
	}
	ch := rabbitmq.NewPublisher(service.RabbitMqConn, PublisherConfig)
	defer ch.CloseChannel()
	ch.SetupExchangeAndQueue()
	ch.Publish(body)
	return 200
}
