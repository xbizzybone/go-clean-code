package utils

import (
	"errors"
	"log"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
	return c.Next()
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
