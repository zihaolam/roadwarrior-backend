package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zihaolam/roadwarrior-backend/entities"
	"github.com/zihaolam/roadwarrior-backend/repos"
	"github.com/zihaolam/roadwarrior-backend/schemas"
)

type IPointHandler struct{}

func NewPointHandler() IPointHandler {
	return IPointHandler{}
}

func (controller *IPointHandler) GetOneHandler(c *fiber.Ctx) error {
	pointId := c.Params("pointId")
	product, err := repos.PointRepo.GetOne(&entities.PointKey{PK: repos.PointRepo.TablePK, SK: pointId})

	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return OKResponse(c, product)
}

func (controller *IPointHandler) GetAllHandler(c *fiber.Ctx) error {
	products, err := repos.PointRepo.GetAll()

	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(products)
}

func (controller *IPointHandler) CreateHandler(c *fiber.Ctx) error {
	var body schemas.CreatePointSchema
	if err := c.BodyParser(&body); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err)
	}

	newPoint, err := repos.PointRepo.Create(&body)
	if err != nil {
		return err
	}
	fmt.Println("reacher hered")

	return OKResponse(c, newPoint)
}

func (controller *IPointHandler) UpdateHandler(c *fiber.Ctx) error {
	pointId := c.Params("pointId")

	var body schemas.UpdatePointSchema
	if err := c.BodyParser(&body); err != nil {
		return ErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	updatedPoint, err := repos.PointRepo.UpdateOne(&entities.PointKey{PK: repos.PointRepo.TablePK, SK: pointId}, &body)

	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err)

	}

	return OKResponse(c, updatedPoint)
}

func (controller *IPointHandler) DeleteHandler(c *fiber.Ctx) error {
	pointId := c.Params("pointId")

	err := repos.PointRepo.DeleteOne(&entities.PointKey{PK: repos.PointRepo.TablePK, SK: pointId})

	if err != nil {
		return ErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	return OKResponse(c, fmt.Sprintf("Successfully Deleted Record: %s", pointId))
}
