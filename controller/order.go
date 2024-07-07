package controller

import (
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sir-shalahuddin/synapsis/dto"
	error_helper "github.com/sir-shalahuddin/synapsis/pkg/helper"
	"github.com/sir-shalahuddin/synapsis/pkg/utils"
)

type orderService interface {
	CreateOrder(ctx context.Context, req *dto.OrderRequest, userID string) (string, error)
	PayOrder(ctx context.Context, userID string, orderID string) (bool, error)
}

type orderController struct {
	svc      orderService
	validate *validator.Validate
}

func NewOrderController(svc orderService) *orderController {
	validate := validator.New()

	return &orderController{
		svc:      svc,
		validate: validate,
	}
}

// CreateOrder godoc
// @Summary Create an order
// @Description Creates a new order for the user
// @Tags order
// @Accept  json
// @Produce  json
// @Param request body dto.OrderRequest true "Order request"
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /orders [post]
func (ctl orderController) CreateOrder(ctx *fiber.Ctx) error {
	var req dto.OrderRequest

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := ctl.validate.Struct(req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	orderID, err := ctl.svc.CreateOrder(ctx.Context(), &req, userID)

	if err != nil {
		log.Printf("CreateOrder - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}

	return utils.HandleSuccess(ctx, "create order success", orderID, fiber.StatusOK)

}

// PayOrder godoc
// @Summary Pay for an order
// @Description Processes payment for a specific order
// @Tags order
// @Accept  json
// @Produce  json
// @Param order_id path string true "Order ID"
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /orders/{order_id}/payments [post]
func (ctl orderController) PayOrder(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	orderID := ctx.Params("order_id")

	_, err := ctl.svc.PayOrder(ctx.Context(), userID, orderID)

	if err != nil {
		if errors.Is(err, error_helper.ErrorUnAuthorized) {
			return utils.HandleError(ctx, err, "forbidden", fiber.StatusForbidden)
		}
		log.Printf("PayOrder - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}

	return utils.HandleSuccess(ctx, "order payment successfull", orderID, fiber.StatusOK)

}
