package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sir-shalahuddin/synapsis/dto"
)

func HandleError(ctx *fiber.Ctx, err error, message string, statusCode int) error {
	if message == "" {
		message = err.Error()
	}
	return ctx.Status(statusCode).JSON(dto.ErrorMessage{
		Error:      true,
		Message:    message,
		StatusCode: statusCode,
	})
}

func HandleSuccess(ctx *fiber.Ctx, message string, data interface{}, statusCode int) error {
	return ctx.Status(statusCode).JSON(dto.Response{
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
	})
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	if len(password) < 8 {
		return false
	}
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasNumber && hasUppercase && hasLowercase && hasSpecialChar
}
