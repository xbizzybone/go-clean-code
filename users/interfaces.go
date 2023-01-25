package users

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	GetUserById(ctx *fiber.Ctx) error
	GetUserByEmail(ctx *fiber.Ctx) error
	Deactivate(ctx *fiber.Ctx) error
	Activate(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type Cases interface {
	Register(user *User) (User, error)
	Login(user *User) (User, error)
	GetUserById(id string) (User, error)
	GetUserByEmail(email string) (User, error)
	Deactivate(id string) error
	Activate(id string) error
	Update(id string, user *User) error
}

type Repository interface {
	CreateUser(user *User) error
	UpdateUser(user *User) error
	GetUserById(id string) (User, error)
	GetUserByEmail(email string) (User, error)
	Deactivate(id string) error
	Activate(id string) error
}
