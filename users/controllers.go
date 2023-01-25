package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xbizzybone/go-clean-code/utils"
)

type controller struct {
	cases Cases
}

func NewController(cases Cases) Controller {
	return &controller{cases}
}

func (c *controller) Register(ctx *fiber.Ctx) error {
	requestBody := new(UserCreateRequest)

	if err := utils.BodyParser(ctx, requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error validando el cuerpo de la petición :" + err.Error(),
		})
	}

	user := new(User)
	user.Email = requestBody.Email
	user.Password = requestBody.Password
	user.Name = requestBody.Name

	createUser, err := c.cases.Register(user)
	if err != nil {
		if err.Error() == "user already exists" {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "El usuario ya existe",
			})
		}

		if err.Error() == "error while hashing password" {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error encriptando la contraseña",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creando el usuario",
		})
	}

	userResponse := new(UserCreateResponse)
	userResponse.Name = createUser.Name
	userResponse.Email = createUser.Email
	userResponse.ID = createUser.ID

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario creado correctamente",
		"user":    userResponse,
	})
}

func (c *controller) Login(ctx *fiber.Ctx) error {
	requestBody := new(UserRequest)

	if err := utils.BodyParser(ctx, requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error validando el cuerpo de la petición :" + err.Error(),
		})
	}

	user := new(User)
	user.Email = requestBody.Email
	user.Password = requestBody.Password

	userResult, err := c.cases.Login(user)
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "El usuario no existe",
			})
		}

		if err.Error() == "invalid password" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Contraseña incorrecta",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error iniciando sesión",
		})
	}

	if userResult.IsDeleted {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Cuenta de usuario deshabilitada",
		})
	}

	response := new(UserLoginResponse)
	response.Email = userResult.Email
	response.Name = userResult.Name
	response.ID = userResult.ID
	response.LastLogin = userResult.LastLogin
	response.CreatedAt = userResult.CreatedAt
	response.UpdatedAt = userResult.UpdatedAt

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *controller) GetUserById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userResult, err := c.cases.GetUserById(id)
	if err != nil {
		if err.Error() == "invalid id" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Usuario no existe",
			})
		}

		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "El usuario no existe",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error obteniendo el usuario",
		})
	}

	response := new(UserResponse)
	response.Email = userResult.Email
	response.Name = userResult.Name
	response.ID = userResult.ID
	response.IsDeleted = userResult.IsDeleted
	response.LastLogin = userResult.LastLogin
	response.CreatedAt = userResult.CreatedAt
	response.UpdatedAt = userResult.UpdatedAt
	response.DeleteAt = userResult.DeleteAt

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *controller) GetUserByEmail(ctx *fiber.Ctx) error {
	email := ctx.Params("email")

	userResult, err := c.cases.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "El usuario no existe",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error obteniendo el usuario",
		})
	}

	response := new(UserResponse)
	response.Email = userResult.Email
	response.Name = userResult.Name
	response.ID = userResult.ID
	response.IsDeleted = userResult.IsDeleted
	response.LastLogin = userResult.LastLogin
	response.CreatedAt = userResult.CreatedAt
	response.UpdatedAt = userResult.UpdatedAt
	response.DeleteAt = userResult.DeleteAt

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *controller) Deactivate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.cases.Deactivate(id)
	if err != nil {
		if err.Error() == "invalid id" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Usuario no existe",
			})
		}

		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "El usuario no existe",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error eliminando el usuario",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario desactivado correctamente",
	})
}

func (c *controller) Activate(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.cases.Activate(id)
	if err != nil {
		if err.Error() == "invalid id" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Usuario no existe",
			})
		}

		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "El usuario no existe",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error activando el usuario",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario activado correctamente",
	})
}

func (c *controller) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	requestBody := new(UserUpdateRequest)
	err := ctx.BodyParser(requestBody)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parseando el body",
		})
	}

	if requestBody.Email != "" && !utils.IsValidEmail(requestBody.Email) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email invalido",
		})
	}

	user := new(User)
	user.Email = requestBody.Email
	user.Name = requestBody.Name
	user.Password = requestBody.Password

	err = c.cases.Update(id, user)
	if err != nil {
		if err.Error() == "invalid id" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Usuario no existe",
			})
		}

		if err.Error() == "user not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "El usuario no existe",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error actualizando el usuario",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Usuario actualizado correctamente",
	})
}
