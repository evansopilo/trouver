package main

import (
	"context"
	"strconv"
	"time"

	"github.com/evansopilo/trouver/internal/data"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreatePlace creates a new place, handler for adding a new place to the application.
func (app *Application) CreatePlace(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var place data.Place

	// decode the request body to place variable declared and continue with the request flow
	// when the decode is successfull otherwise return a status bad request back to the client.
	if err := c.BodyParser(&place); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request",
		})
	}

	// add place user id to user id obtained from auth token claims.
	place.UserID = c.Locals("user_id").(string)

	// add timestamp of current time to the create time of place object.
	place.CreatedAt = time.Now()

	// add place id from a random generated uuid.
	place.ID = uuid.New().String()

	// insert the new place record to the database within the defined context with timeout. When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	if err := app.Models.Place.InsertOne(ctx, "trouver", "places", &place); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "create place failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "create place operation success",
		"data": map[string]interface{}{
			"id": place.ID,
		},
	})
}

// GetPlace gets place, handler for getting a place from the application by given id.
func (app *Application) GetPlace(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get place record with the provided id from the database within the defined context with timeout.When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	place, err := app.Models.Place.FindOne(ctx, "trouver", "places", c.Params("place_id"))
	if err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "read place failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "read place operation success",
		"data": map[string]interface{}{
			"place": place,
		},
	})
}

// ListPlace lists places, handler from updating a place in the application.
func (app *Application) ListPlace(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	page_size, _ := strconv.Atoi(c.Query("size", "10"))
	skip := (page-1)*page_size + 1

	// get place records with the provided filters from the database within the defined context with timeout.When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	places, err := app.Models.Place.List(ctx, "trouver", "places", data.Filter{Skip: skip, Limit: page_size})
	if err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "read place failed",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   places,
	})
}

// UpdatePlace updates place, handler from updating a place in the application.
func (app *Application) UpdateOne(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var place data.Place

	// add place id to id obtained from prams
	place.ID = c.Params("place_id")

	// add user id to user id obtained from auth token claims.
	place.UserID = c.Locals("user_id").(string)

	// decode the request body to place variable declared and continue with the request flow
	// when the decode is successfull otherwise return a status bad request back to the client.
	if err := c.BodyParser(&place); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user request",
		})
	}

	existingPlace, err := app.Models.Place.FindOne(ctx, "trouver", "places", c.Params("place_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	// for successfull update, the place to be updated must be specific to the user or the user must be an admin.
	// otherwise returns a status forbidden(user has no permission to update the record).
	if !(existingPlace.UserID == c.Locals("user_id").(string) || c.Locals("user_role").(string) == "admin") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	// update the place record to the database within the defined context with timeout. When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	if err := app.Models.Place.UpdateOne(ctx, "trouver", "places", &place); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "update place success",
		"data": map[string]interface{}{
			"id": c.Params("place_id"),
		},
	})
}

// DeletePlace delete place, handler from deleting a place in the application.
func (app *Application) DeletePlace(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the place with the provided id from the database.
	existingPlace, err := app.Models.Place.FindOne(ctx, "trouver", "places", c.Params("place_id"))
	if err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "delete place failed",
		})
	}

	// for successfull delete, the place to be delete must be specific to the user or the user must be an admin.
	// otherwise returns a status forbidden(user has no permission to delete the record).
	if !(existingPlace.UserID == c.Locals("user_id").(string) || c.Locals("user_role") == "admin") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	// delete the place record to the database within the defined context with timeout. When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	if err := app.Models.Place.DeleteOne(ctx, "trouver", "places", c.Params("place_id")); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "delete place failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "delete place success",
		"data": map[string]interface{}{
			"id": c.Params("place_id"),
		},
	})
}
