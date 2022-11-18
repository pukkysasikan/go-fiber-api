package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = configs.GetCollection(configs.DB, "menus")
var validate = validator.New()

func CreateMenu(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var menu models.Menu
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.MenuResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&menu); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.MenuResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newMenu := models.Menu{
		Name:        menu.Name,
		Category:    menu.Category,
		Price:       menu.Price,
		Description: menu.Description,
		Image:       menu.Image,
	}

	result, err := menuCollection.InsertOne(ctx, newMenu)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.MenuResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAMenu(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	menuId := c.Params("menuId")
	var menu models.Menu
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(menuId)

	err := menuCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&menu)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.MenuResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": menu}})
}

func EditAMenu(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	menuId := c.Params("menuId")
	var menu models.Menu
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(menuId)

	//validate the request body
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.MenuResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&menu); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.MenuResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": menu.Name, "price": menu.Price, "description": menu.Description}

	result, err := menuCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedMenu models.Menu
	if result.MatchedCount == 1 {
		err := menuCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedMenu)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.MenuResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedMenu}})
}

func DeleteAMenu(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	menuId := c.Params("menuId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(menuId)

	result, err := menuCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.MenuResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Menu with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.MenuResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Menu successfully deleted!"}},
	)
}

func GetAllMenus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var menus []models.Menu
	defer cancel()

	results, err := menuCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.Menu
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.MenuResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		menus = append(menus, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.MenuResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": menus}},
	)
}
