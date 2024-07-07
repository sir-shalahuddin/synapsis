package controller

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/synapsis/model"
	"github.com/sir-shalahuddin/synapsis/pkg/utils"
)

type productService interface {
	GetProducts(ctx context.Context, category string) ([]model.Product, error)
}

type productController struct {
	svc      productService
	validate *validator.Validate
}

func NewProductController(svc productService) *productController {
	validate := validator.New()

	return &productController{
		svc:      svc,
		validate: validate,
	}
}

// GetProducts godoc
// @Summary Get products
// @Description Retrieves products, optionally filtered by category
// @Tags product
// @Accept  json
// @Produce  json
// @Param category query string false "Category filter"
// @Success 200 {object} dto.Response
// @Failure 500 {object} dto.ErrorMessage
// @Router /products [get]
func (ctl productController) GetProducts(ctx *fiber.Ctx) error {

	category := ctx.Query("category")
	products, err := ctl.svc.GetProducts(ctx.Context(), category)

	if err != nil {
		log.Printf("GetProducts - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}
	return utils.HandleSuccess(ctx, "get product success", products, fiber.StatusOK)

}
