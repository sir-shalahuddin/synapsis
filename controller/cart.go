package controller

import (
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sir-shalahuddin/synapsis/dto"
	"github.com/sir-shalahuddin/synapsis/model"
	error_helper "github.com/sir-shalahuddin/synapsis/pkg/helper"
	"github.com/sir-shalahuddin/synapsis/pkg/utils"
)

type cartService interface {
	AddProduct(ctx context.Context, req *dto.CartRequest, userID string) (*model.Cart, error)
	GetProducts(ctx context.Context, userID string) ([]model.Product, error)
	DeleteProduct(ctx context.Context, userID string, productID string) error
}

type cartController struct {
	svc      cartService
	validate *validator.Validate
}

func NewCartController(svc cartService) *cartController {
	validate := validator.New()

	return &cartController{
		svc:      svc,
		validate: validate,
	}
}

// AddProduct godoc
// @Summary Add product to cart
// @Description Adds a product to the user's cart
// @Tags cart
// @Accept  json
// @Produce  json
// @Param request body dto.CartRequest true "Cart request"
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /carts [post]
func (ctl cartController) AddProduct(ctx *fiber.Ctx) error {
	var req dto.CartRequest

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := ctl.validate.Struct(req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	cart, err := ctl.svc.AddProduct(ctx.Context(), &req, userID)

	if err != nil {
		if errors.Is(err, error_helper.ErrorDuplicate) {
			return utils.HandleError(ctx, err, "product already added to cart", fiber.StatusConflict)
		}
		if errors.Is(err, error_helper.ErrorNotFound) {
			return utils.HandleError(ctx, err, "product not found", fiber.StatusNotFound)
		}
		log.Printf("AddProductCart - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}
	return utils.HandleSuccess(ctx, "add product to cart success", cart, fiber.StatusOK)

}

// GetProducts godoc
// @Summary Get products in cart
// @Description Retrieves all products in the user's cart
// @Tags cart
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 500 {object} dto.ErrorMessage
// @Router /carts [get]
func (ctl cartController) GetProducts(ctx *fiber.Ctx) error {

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	carts, err := ctl.svc.GetProducts(ctx.Context(), userID)

	if err != nil {
		log.Printf("GetProductsCart - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}
	return utils.HandleSuccess(ctx, "get products success", carts, fiber.StatusOK)

}

// DeleteProduct godoc
// @Summary Delete product from cart
// @Description Deletes a product from the user's cart
// @Tags cart
// @Accept  json
// @Produce  json
// @Param request body dto.CartRequest true "Cart request"
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.ErrorMessage
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /carts [delete]
func (ctl cartController) DeleteProduct(ctx *fiber.Ctx) error {
	var req dto.CartRequest

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := ctl.validate.Struct(req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	err := ctl.svc.DeleteProduct(ctx.Context(), userID, req.ProductID)

	if err != nil {
		if errors.Is(err, error_helper.ErrorNotFound) {
			return utils.HandleError(ctx, err, "product not found", fiber.StatusNotFound)
		}
		log.Printf("DeleteProductCart - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}
	return utils.HandleSuccess(ctx, "delete product in cart success", nil, fiber.StatusOK)

}
