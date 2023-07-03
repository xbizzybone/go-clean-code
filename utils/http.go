package utils

import (
	"context"
	"errors"
	"log"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var Validator = validator.New()

func validateSchema(data interface{}) error {
	if err := Validator.Struct(data); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func BodyParser(ctx *fiber.Ctx, data interface{}) error {
	if err := ctx.BodyParser(data); err != nil {
		log.Println(err)
		return errors.New("the request body is invalid")
	}
	if err := validateSchema(data); err != nil {
		return err
	}
	return nil
}

func GetNextMiddleWare(c *fiber.Ctx) error {
	request_id := uuid.New().String()
	ctx := context.WithValue(c.Context(), "request_id", request_id)
	c.SetUserContext(ctx)

	return c.Next()
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
