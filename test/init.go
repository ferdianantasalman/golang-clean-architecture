package test

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/config"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var log *logrus.Logger

var validate *validator.Validate

func init() {
	godotenv.Load() // ponytail: silent fail if .env missing
	log = config.NewLogger()
	validate = config.NewValidator()
	app = config.NewFiber()
	db = config.NewDatabase(log)
	producer := config.NewKafkaProducer(log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Producer: producer,
	})
}
