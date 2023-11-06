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
