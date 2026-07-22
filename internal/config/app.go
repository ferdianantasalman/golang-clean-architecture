package config

import (
	"golang-clean-architecture/internal/delivery/http"
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/route"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/internal/util"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Producer *kafka.Producer
}

func Bootstrap(config *BootstrapConfig) {
	userRepository := repository.NewUserRepository(config.Log)
	contactRepository := repository.NewContactRepository(config.Log)
	addressRepository := repository.NewAddressRepository(config.Log)
	refreshTokenRepository := repository.NewRefreshTokenRepository(config.Log)

	var userProducer *messaging.UserProducer
	var contactProducer *messaging.ContactProducer
	var addressProducer *messaging.AddressProducer

	if config.Producer != nil {
		userProducer = messaging.NewUserProducer(config.Producer, config.Log)
		contactProducer = messaging.NewContactProducer(config.Producer, config.Log)
		addressProducer = messaging.NewAddressProducer(config.Producer, config.Log)
	}

	redisClient := NewRedis()

	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, refreshTokenRepository, userProducer, redisClient)
	contactUseCase := usecase.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository, contactProducer)
	addressUseCase := usecase.NewAddressUseCase(config.DB, config.Log, config.Validate, contactRepository, addressRepository, addressProducer)

	userController := http.NewUserController(userUseCase, config.Log)
	contactController := http.NewContactController(contactUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)

	rateLimiterUtil := util.NewRateLimiterUtil(redisClient)

	authMiddleware := middleware.NewAuth(userUseCase, rateLimiterUtil)

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
		ContactController: contactController,
		AddressController: addressController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}
