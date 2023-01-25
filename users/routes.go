package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xbizzybone/go-clean-code/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var repo Repository
var ctrl Controller
var cs Cases

func bootstrap(logger *zap.Logger, mongoCollection *mongo.Collection) {
	repo = NewRepository(logger, mongoCollection)
	cs = NewCases(logger, repo)
	ctrl = NewController(cs)
}

func ApplyRoutes(app *fiber.App, logger *zap.Logger, mongoCollection *mongo.Collection) {
	bootstrap(logger, mongoCollection)
	group := app.Group("/auth", utils.GetNextMiddleWare)
	group.Post("/register", ctrl.Register)
	group.Post("/login", ctrl.Login)
	group.Get("/user/:id", ctrl.GetUserById)
	group.Get("/user/email/:email", ctrl.GetUserByEmail)
	group.Put("/user/activate/:id", ctrl.Activate)
	group.Delete("/user/deactivate/:id", ctrl.Deactivate)
	group.Put("/user/:id", ctrl.Update)
}
