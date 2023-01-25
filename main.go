package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/xbizzybone/go-clean-code/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var client *mongo.Client
var usersCollection *mongo.Collection

var zapLogger *zap.Logger

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.TimeKey = "timestamp"
	logConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	zapLogger, err = logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	zapLogger = zapLogger.With(zap.String("service", "go-clean-code"))
	zapLogger.Info("Logger initialized")

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING"))
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	usersCollection = client.Database(dbName).Collection(os.Getenv("MONGO_USERS_COLLECTION_NAME"))
	zapLogger.Info("MongoDB initialized")
}

func main() {
	defer zapLogger.Sync()

	zapLogger.Info("Starting server")

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	users.ApplyRoutes(app, zapLogger, usersCollection)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
