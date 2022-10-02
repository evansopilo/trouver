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

// CreateReview creates a new review, handler for adding a new review to the application.
func (app *Application) CreateReview(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var review data.Review

	// decode the request body to review variable declared and continue with the request flow
	// when the decode is successfull otherwise return a status bad request back to the client.
	if err := c.BodyParser(&review); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid request",
		})
	}

	// add review user id to user id obtained from auth token claims.
	review.UserID = c.Locals("user_id").(string)

	// add timestamp of current time to the create time of review object.
	review.CreatedAt = time.Now()

	// add review id from a random generated uuid.
	review.ID = uuid.New().String()

	// insert the new review record to the database within the defined context with timeout. When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	if err := app.Models.Review.InsertOne(ctx, "trouver", "reviews", &review); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "create review failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "create review operation success",
		"data": map[string]interface{}{
			"id": review.ID,
		},
	})
}

// ListPlace lists places, handler from updating a place in the application.
func (app *Application) ListReview(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	page_size, _ := strconv.Atoi(c.Query("size", "10"))
	skip := (page-1)*page_size + 1

	// get review records with the provided filters from the database within the defined context with timeout.When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	reviews, err := app.Models.Review.List(ctx, "trouver", "reviews", c.Params("place_id"), data.Filter{Skip: skip, Limit: page_size})
	if err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "read review failed",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   reviews,
	})
}

func (app *Application) UpdateReview(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var review data.Review

	// add review id to id obtained from prams
	review.ID = c.Params("review_id")

	// add user id to user id obtained from auth token claims.
	review.UserID = c.Locals("user_id").(string)

	// decode the request body to review variable declared and continue with the request flow
	// when the decode is successfull otherwise return a status bad request back to the client.
	if err := c.BodyParser(&review); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user request",
		})
	}

	existingReview, err := app.Models.Review.FindOne(ctx, "trouver", "reviews", c.Params("review_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	// for successfull update, the place to be updated must be specific to the user or the user must be an admin.
	// otherwise returns a status forbidden(user has no permission to update the record).
	if !(existingReview.UserID == c.Locals("user_id").(string) || c.Locals("user_role").(string) == "admin") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	// update the place record to the database within the defined context with timeout. When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	if err := app.Models.Review.UpdateOne(ctx, "trouver", "places", &review); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "update place failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "update review success",
		"data": map[string]interface{}{
			"id": c.Params("review_id"),
		},
	})
}

func (app *Application) DeleteReview(c *fiber.Ctx) error {

	// create a context with a 5-second timeout deadline. The entire request response cycle will be
	// tied to this context, therefore response should be returned within the defined context timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the review with the provided id from the database.
	existingReview, err := app.Models.Review.FindOne(ctx, "trouver", "reviews", c.Params("review_id"))
	if err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "delete review failed",
		})
	}

	// for successfull delete, the review to be delete must be specific to the user or the user must be an admin.
	// otherwise returns a status forbidden(user has no permission to delete the record).
	if !(existingReview.UserID == c.Locals("user_id").(string) || c.Locals("user_role") == "admin") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "update review failed",
		})
	}

	// delete the review record to the database within the defined context with timeout. When the
	// operation fails for any reason ie. elapsed context timeout reponse with a 500 Internal Server error
	// is returned.
	if err := app.Models.Review.DeleteOne(ctx, "trouver", "reviews", c.Params("review_id")); err != nil {
		logrus.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "delete review failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "delete place success",
		"data": map[string]interface{}{
			"id": c.Params("review_id"),
		},
	})
}
