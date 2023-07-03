package users

import (
	"context"

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
	Register(ctx context.Context, user *User) (User, error)
	Login(ctx context.Context, user *User) (User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
	Update(ctx context.Context, id string, user *User) error
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
}
