package routers

import (
	_ "app/docs"
	"app/models"
	"app/store"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// FetchURL godoc
// @Summary		Fetch the shortened Url from the redis and then redirect to the original url upon success.
// @Description	redirect to original url by tag (shortened url tag)
// @Tags			Fetch URL
// @ID				redirect-from-custom-url
// @Accept			json
// @Produce		    json
// @Param			tag	path		string	true	"Custom Tag"
// @Success		301
// @Failure		404	{object}	models.RequestError
// @Failure		500	{object}	models.RequestError
// @Router			/{tag} [get]
func Fetch(context *fiber.Ctx) error {

	// API handler to redirect the user to the link that is created.

	customTag := context.Params("tag")

	// we will see if the custom tag exists in the database or not.
	// it will return us the original URL
	// ! We will strictly do the 301 redirect upon fetching url successfully

	// Create Database, Connect
	redisDb := store.CreateStore()

	// delay the execution of the following method by the end of the execution
	defer redisDb.Close()

	// get the url

	url, urlErr := redisDb.Get(store.Context.Ctx, customTag).Result()

	if urlErr == redis.Nil {
		return context.Status(fiber.StatusNotFound).JSON(&models.RequestError{
			Error:      "the url you are looking for is not found.",
			StatusCode: fiber.StatusNotFound,
		})
	} else if urlErr != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(&models.RequestError{
			Error:      "Can not connect to the database.",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Redirect to the original URL

	return context.Redirect(url, fiber.StatusMovedPermanently)
}
