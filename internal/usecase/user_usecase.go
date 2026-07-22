package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/util"
)

type UserUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	UserRepository        *repository.UserRepository
	RefreshTokenRepository *repository.RefreshTokenRepository
	UserProducer          *messaging.UserProducer
	Redis                 *redis.Client
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, refreshTokenRepository *repository.RefreshTokenRepository,
	userProducer *messaging.UserProducer, redis *redis.Client) *UserUseCase {
	return &UserUseCase{
		DB:                     db,
		Log:                    logger,
		Validate:               validate,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
		UserProducer:           userProducer,
		Redis:                  redis,
	}
}

func (c *UserUseCase) Verify(ctx context.Context, request *model.VerifyUserRequest) (*model.Auth, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	claims, err := util.ValidateAccessToken(request.Token)
	if err != nil {
		c.Log.Warnf("Failed validate JWT : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	blacklisted, err := c.Redis.Exists(ctx, "bl:"+claims.JTI).Result()
	if err != nil {
		c.Log.Warnf("Failed check token blacklist : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if blacklisted > 0 {
		c.Log.Warnf("Token is blacklisted")
		return nil, fiber.ErrUnauthorized
	}

	return &model.Auth{ID: claims.UserID}, nil
}

func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:       uuid.New(),
		Password: string(password),
		Name:     request.Name,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user created event")
		if err = c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user created event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user created event")
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body  : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByName(tx, user, request.Name); err != nil {
		c.Log.Warnf("Failed find user by name : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := util.GenerateAccessToken(user.ID, user.Name)
	if err != nil {
		c.Log.Warnf("Failed generate JWT : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshValue := uuid.New()
	refreshToken := refreshValue.String()
	refreshExpiry := time.Now().Add(util.RefreshExpiryDuration()).UnixMilli()

	refreshEntity := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: util.Sha256Hex(refreshToken),
		ExpiresAt: refreshExpiry,
		CreatedAt: time.Now().UnixMilli(),
	}
	if err := c.RefreshTokenRepository.Create(tx, refreshEntity); err != nil {
		c.Log.Warnf("Failed create refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user login event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user login event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user login event")
	}

	return converter.UserToTokenResponse(user, accessToken, refreshToken), nil
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return false, fiber.ErrNotFound
	}

	claims, err := util.ValidateAccessToken(request.Token)
	if err != nil {
		c.Log.Warnf("Failed validate JWT for blacklist : %+v", err)
		return false, fiber.ErrUnauthorized
	}

	if err := c.Redis.Set(ctx, "bl:"+claims.JTI, "1", 7*24*time.Hour).Err(); err != nil {
		c.Log.Warnf("Failed to blacklist token : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := c.RefreshTokenRepository.DeleteByUserId(tx, request.ID); err != nil {
		c.Log.Warnf("Failed delete refresh tokens : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user logout event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user logout event : %+v", err)
			return false, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user logout event")
	}

	return true, nil
}

func (c *UserUseCase) Refresh(ctx context.Context, request *model.RefreshTokenRequest) (*model.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	refreshHash := util.Sha256Hex(request.RefreshToken)
	refreshEntity := new(entity.RefreshToken)
	if err := c.RefreshTokenRepository.FindByHash(tx, refreshEntity, refreshHash); err != nil {
		c.Log.Warnf("Failed find refresh token : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if time.Now().UnixMilli() > refreshEntity.ExpiresAt {
		c.Log.Warnf("Refresh token expired")
		return nil, fiber.ErrUnauthorized
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, refreshEntity.UserID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := c.RefreshTokenRepository.Delete(tx, refreshEntity); err != nil {
		c.Log.Warnf("Failed delete old refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	accessToken, err := util.GenerateAccessToken(user.ID, user.Name)
	if err != nil {
		c.Log.Warnf("Failed generate JWT : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newRefreshValue := uuid.New()
	newRefreshToken := newRefreshValue.String()
	newRefreshExpiry := time.Now().Add(util.RefreshExpiryDuration()).UnixMilli()

	newRefreshEntity := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: util.Sha256Hex(newRefreshToken),
		ExpiresAt: newRefreshExpiry,
		CreatedAt: time.Now().UnixMilli(),
	}
	if err := c.RefreshTokenRepository.Create(tx, newRefreshEntity); err != nil {
		c.Log.Warnf("Failed create new refresh token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToTokenResponse(user, accessToken, newRefreshToken), nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Log.Warnf("Failed to generate bcrype hash : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
		user.Password = string(password)
	}

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if c.UserProducer != nil {
		event := converter.UserToEvent(user)
		c.Log.Info("Publishing user updated event")
		if err := c.UserProducer.Send(event); err != nil {
			c.Log.Warnf("Failed publish user updated event : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		c.Log.Info("Kafka producer is disabled, skipping user updated event")
	}

	return converter.UserToResponse(user), nil
}
