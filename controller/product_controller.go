package controller

import (
	"log"

	"product-service/model"
	"product-service/service"

	"github.com/gofiber/fiber/v2"
)

func (controller *ProductController) Route(app *fiber.App) {
	app.Get("products", controller.ListProduct)
	app.Get("product/:id", controller.FindProduct)
	app.Post("product", controller.Store)
	app.Put("product/:id", controller.Update)
	app.Delete("product/:id", controller.Delete)
}

type ProductController struct {
	Service service.ProductService
}

func NewProductController(service service.ProductService) ProductController {
	return ProductController{Service: service}
}

func (controller *ProductController) ListProduct(ctx *fiber.Ctx) error {
	errorCode, products := controller.Service.GetAllProduct()
	response := model.GetResponse(errorCode, products, "Success")
	log.Print(products)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (controller *ProductController) FindProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	errorCode, product := controller.Service.FindProductById(id)
	response := model.GetResponse(errorCode, product, "Success")

	if errorCode == 404 {
		response := model.GetResponse(errorCode, nil, "Not Found")
		return ctx.Status(fiber.StatusNotFound).JSON(response)
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (controller *ProductController) Store(ctx *fiber.Ctx) error {
	var request model.CreateProductRequest
	ctx.BodyParser(&request)
	responseCode := controller.Service.InsertProduct(request)
	response := model.GetResponse(responseCode, nil, "")
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (controller *ProductController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var request model.CreateProductRequest
	ctx.BodyParser(&request)
	responseCode, product := controller.Service.EditProduct(id, request)

	if responseCode == 200 {
		response := model.GetResponse(responseCode, product, "Success Update")
		return ctx.Status(fiber.StatusOK).JSON(response)
	}

	response := model.GetResponse(responseCode, nil, "Failed Update")
	return ctx.Status(responseCode).JSON(response)
}

func (controller *ProductController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	responseCode := controller.Service.DeleteProduct(id)

	if responseCode == 200 {
		response := model.GetResponse(responseCode, nil, "Success Delete")
		return ctx.Status(fiber.StatusOK).JSON(response)
	}

	response := model.GetResponse(responseCode, nil, "Failed delete")
	return ctx.Status(responseCode).JSON(response)
}
