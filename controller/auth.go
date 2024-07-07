package controller

import (
	"context"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/synapsis/dto"
	error_helper "github.com/sir-shalahuddin/synapsis/pkg/helper"
	"github.com/sir-shalahuddin/synapsis/pkg/utils"
)

type userService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, error)
}

type userController struct {
	svc      userService
	validate *validator.Validate
}

func NewUserController(svc userService) *userController {
	validate := validator.New()
	validate.RegisterValidation("password", utils.ValidatePassword)

	return &userController{
		svc:      svc,
		validate: validate,
	}
}

// Register godoc
// @Summary User registration
// @Description Registers a new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /auth/register [post]
func (ctl userController) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := ctl.validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			switch e.Field() {
			case "Password":
				return utils.HandleError(ctx, err, "invalid password", fiber.StatusBadRequest)
			case "Email":
				return utils.HandleError(ctx, err, "invalid email", fiber.StatusBadRequest)
			default:
				return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
			}
		}
	}

	user, err := ctl.svc.Register(ctx.Context(), &req)

	if err != nil {
		if errors.Is(err, error_helper.ErrorDuplicateEmail) {
			return utils.HandleError(ctx, err, "email already registered", fiber.StatusConflict)
		}
		log.Printf("Register Service - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}

	return utils.HandleSuccess(ctx, "register success", user, fiber.StatusOK)

}

// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /auth/login [post]
func (ctl userController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
	}

	if err := ctl.validate.Struct(req); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			switch e.Field() {
			case "Password":
				return utils.HandleError(ctx, err, "invalid password", fiber.StatusBadRequest)
			case "Email":
				return utils.HandleError(ctx, err, "invalid email", fiber.StatusBadRequest)
			default:
				return utils.HandleError(ctx, err, "invalid payload", fiber.StatusBadRequest)
			}
		}
	}

	token, err := ctl.svc.Login(ctx.Context(), &req)

	if err != nil {
		if errors.Is(err, error_helper.ErrorAuthentication) {
			return utils.HandleError(ctx, err, "", fiber.StatusBadRequest)
		}
		log.Printf("Login - Service error: %v", err)
		return utils.HandleError(ctx, err, "internal error", fiber.StatusInternalServerError)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "login successful",
		"token":   token,
	})
}
