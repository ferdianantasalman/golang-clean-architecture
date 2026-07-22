package http

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
)

type AddressController struct {
	UseCase *usecase.AddressUseCase
	Log     *logrus.Logger
}

func NewAddressController(useCase *usecase.AddressUseCase, log *logrus.Logger) *AddressController {
	return &AddressController{
		Log:     log,
		UseCase: useCase,
	}
}

func (c *AddressController) Create(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.CreateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = uuid.MustParse(ctx.Params("contactId"))

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to create address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{
		Status:  "success",
		Message: "Resource created successfully",
		Data:    response,
	})
}

func (c *AddressController) List(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	contactId := uuid.MustParse(ctx.Params("contactId"))

	request := &model.ListAddressRequest{
		UserId:    auth.ID,
		ContactId: contactId,
	}

	responses, err := c.UseCase.List(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to list addresses")
		return err
	}

	return ctx.JSON(model.WebResponse[[]model.AddressResponse]{
		Status:  "success",
		Message: "Resources fetched successfully",
		Data:    responses,
	})
}

func (c *AddressController) Get(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetAddressRequest{
		UserId:    auth.ID,
		ContactId: uuid.MustParse(ctx.Params("contactId")),
		ID:        uuid.MustParse(ctx.Params("addressId")),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to get address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{
		Status:  "success",
		Message: "Resource fetched successfully",
		Data:    response,
	})
}

func (c *AddressController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateAddressRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("failed to parse request body")
		return fiber.ErrBadRequest
	}

	request.UserId = auth.ID
	request.ContactId = uuid.MustParse(ctx.Params("contactId"))
	request.ID = uuid.MustParse(ctx.Params("addressId"))

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("failed to update address")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AddressResponse]{
		Status:  "success",
		Message: "Resource updated successfully",
		Data:    response,
	})
}

func (c *AddressController) Delete(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.DeleteAddressRequest{
		UserId:    auth.ID,
		ContactId: uuid.MustParse(ctx.Params("contactId")),
		ID:        uuid.MustParse(ctx.Params("addressId")),
	}

	if err := c.UseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("failed to delete address")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Status:  "success",
		Message: "Resource deleted successfully",
		Data:    true,
	})
}
