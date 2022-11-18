package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "categorys")

func CreateCategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var category models.Category
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CategoryResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&category); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CategoryResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newCategory := models.Category{
		Name: category.Name,
	}

	result, err := categoryCollection.InsertOne(ctx, newCategory)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.CategoryResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetACategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	categoryId := c.Params("categoryId")
	var category models.Category
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(categoryId)

	err := categoryCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&category)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.CategoryResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": category}})
}

func EditACategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	categoryId := c.Params("categoryId")
	var category models.Category
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(categoryId)

	//validate the request body
	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CategoryResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&category); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.CategoryResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": category.Name}

	result, err := categoryCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	//get updated user details
	var updatedCategory models.Category
	if result.MatchedCount == 1 {
		err := categoryCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedCategory)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.CategoryResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedCategory}})
}

func DeleteACategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	categoryId := c.Params("categoryId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(categoryId)

	result, err := categoryCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.CategoryResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Category with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CategoryResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Category successfully deleted!"}},
	)
}

func GetAllCategorys(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var categorys []models.Category
	defer cancel()

	results, err := categoryCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.Category
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.CategoryResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		categorys = append(categorys, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.CategoryResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": categorys}},
	)
}
