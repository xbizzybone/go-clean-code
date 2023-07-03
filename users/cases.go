package users

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type cases struct {
	logger     *zap.Logger
	repository Repository
}

func NewCases(logger *zap.Logger, repository Repository) Cases {
	return &cases{
		logger:     logger,
		repository: repository,
	}
}

func (c *cases) Register(ctx context.Context, user *User) (User, error) {
	if _, err := c.repository.GetUserByEmail(ctx, user.Email); err == nil {
		return User{}, errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return User{}, errors.New("error while hashing password")
	}
	user.Password = hashedPassword

	if err = c.repository.CreateUser(ctx, user); err != nil {
		return User{}, errors.New("error while creating user")
	}

	createUser, err := c.repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return User{}, errors.New("error while getting user")
	}

	return createUser, nil
}

func (c *cases) Login(ctx context.Context, user *User) (User, error) {
	result, err := c.repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	if !validatePassword(result.Password, user.Password) {
		return User{}, errors.New("invalid password")
	}

	result.LastLogin = primitive.NewDateTimeFromTime(time.Now())
	if err = c.repository.UpdateUser(ctx, &result); err != nil {
		return User{}, errors.New("error while updating user")
	}

	return result, nil
}

func (c *cases) GetUserById(ctx context.Context, id string) (User, error) {
	return c.repository.GetUserById(ctx, id)
}

func (c *cases) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return c.repository.GetUserByEmail(ctx, email)
}

func (c *cases) Deactivate(ctx context.Context, id string) error {
	return c.repository.Deactivate(ctx, id)
}

func (c *cases) Activate(ctx context.Context, id string) error {
	return c.repository.Activate(ctx, id)
}

func (c *cases) Update(ctx context.Context, id string, user *User) error {
	userToUpdate, err := c.repository.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	if user.Name != "" {
		userToUpdate.Name = user.Name
	}

	if user.Email != "" {
		userToUpdate.Email = user.Email
	}

	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return errors.New("error while hashing password")
		}
		userToUpdate.Password = hashedPassword
	}

	return c.repository.UpdateUser(ctx, user)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error while hashing password")
	}
	return string(hashedPassword), nil
}

func validatePassword(user_password, comming_password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user_password), []byte(comming_password)); err != nil {
		return false
	}
	return true
}
