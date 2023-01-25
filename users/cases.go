package users

import (
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

func (c *cases) Register(user *User) (User, error) {
	if _, err := c.repository.GetUserByEmail(user.Email); err == nil {
		return User{}, errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return User{}, errors.New("error while hashing password")
	}
	user.Password = hashedPassword

	if err = c.repository.CreateUser(user); err != nil {
		return User{}, errors.New("error while creating user")
	}

	createUser, err := c.repository.GetUserByEmail(user.Email)
	if err != nil {
		return User{}, errors.New("error while getting user")
	}

	return createUser, nil
}

func (c *cases) Login(user *User) (User, error) {
	result, err := c.repository.GetUserByEmail(user.Email)
	if err != nil {
		return User{}, errors.New("user not found")
	}

	if !validatePassword(result.Password, user.Password) {
		return User{}, errors.New("invalid password")
	}

	result.LastLogin = primitive.NewDateTimeFromTime(time.Now())
	if err = c.repository.UpdateUser(&result); err != nil {
		return User{}, errors.New("error while updating user")
	}

	return result, nil
}

func (c *cases) GetUserById(id string) (User, error) {
	return c.repository.GetUserById(id)
}

func (c *cases) GetUserByEmail(email string) (User, error) {
	return c.repository.GetUserByEmail(email)
}

func (c *cases) Deactivate(id string) error {
	return c.repository.Deactivate(id)
}

func (c *cases) Activate(id string) error {
	return c.repository.Activate(id)
}

func (c *cases) Update(id string, user *User) error {
	userToUpdate, err := c.repository.GetUserById(id)
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

	return c.repository.UpdateUser(user)
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
