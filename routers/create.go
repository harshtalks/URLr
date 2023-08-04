package routers

import (
	"app/models"
	"app/store"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// ShortenedURL godoc
// @Summary		Create a shortened URL
// @Description	create the shortened url by providing the url, expiration time, and custom tag. All are necessary. for server generated tag, give custom tag as empty string.
// @Tags			Create Shortened URL
// @Accept			json
// @Produce		json
// @Param data body models.CreateRequest true "The input request"
// @Success		200	{object}	models.CreateResponse
// @Failure		400	{object}	models.RequestError
// @Failure		503	{object}	models.BadRequestError
// @Failure		500	{object}	models.BadRequestError
// @Router			/create [post]
func Create(context *fiber.Ctx) error {

	// get the parameter `url` from the request.

	requestBody := new(models.CreateRequest)

	if bodyParseError := context.BodyParser(&requestBody); bodyParseError != nil {
		return context.Status(fiber.StatusBadRequest).JSON(&models.RequestError{
			Error:      "Invalid argument, unable to parse JSON",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// connecting with the database

	redisDb := store.CreateStore()

	// delaying the execution of the following method
	defer redisDb.Close()

	// Implementing Rate Limiting

	if err := redisDb.Get(store.Context.Ctx, context.IP()).Err(); err == redis.Nil {
		// setting the rate limit for the current IP Address.
		// default time is 30 minutes and rate is 30 per 30 minutes
		_ = redisDb.Set(store.Context.Ctx, context.IP(), 30, 30*60*time.Second)
	} else {
		rateLimitStatus, _ := redisDb.Get(store.Context.Ctx, context.IP()).Result()
		valueConvertedToInt, _ := strconv.Atoi(rateLimitStatus)

		if valueConvertedToInt <= 0 {
			// Get the time before which it expires
			limit, _ := redisDb.TTL(store.Context.Ctx, context.IP()).Result()

			// return error as user has finished its quota.
			return context.Status(fiber.StatusServiceUnavailable).JSON(&models.BadRequestError{
				Error:       "You have finished your quota for today. Please Try again later.",
				TimeToReset: int(limit / time.Minute / time.Second),
				StatusCode:  fiber.StatusServiceUnavailable,
			})
		}
	}

	// Checking the valid url.

	if isValid := govalidator.IsURL(requestBody.Url); !isValid {
		return context.Status(fiber.StatusBadRequest).JSON(&models.RequestError{
			Error:      "Bad Request, Please provider a valid URL.",
			StatusCode: fiber.StatusBadRequest,
		})
	}

	// Convert into https

	requestBody.EnforceHttp()

	// Check if the  body time is given as 0

	requestBody.DefaultExpireTime()

	// UUIDs

	// ! No Collision Check here

	var uniqueId string

	if requestBody.CustomTag == "" {
		uniqueId = uuid.New().String()[0:6]
	} else {
		// check if collision exists for custom tag
		// if result > 0 that means key exists

		if result, _ := redisDb.Exists(store.Context.Ctx, requestBody.CustomTag).Result(); result > 0 {
			return context.Status(fiber.StatusBadRequest).JSON(&models.RequestError{
				Error:      "The given key already exists",
				StatusCode: fiber.StatusBadRequest,
			})
		}

		uniqueId = requestBody.CustomTag
	}

	// setting up in database

	if setErr := redisDb.Set(store.Context.Ctx, uniqueId, requestBody.Url, time.Duration(time.Duration(requestBody.ExpireTime)*time.Hour)*time.Second*3600).Err(); setErr != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(&models.RequestError{
			Error:      "Unable to connect to server",
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	// Returning the Status Body

	responseBody := new(models.CreateResponse)

	responseBody.Url = requestBody.Url
	responseBody.CustomTag = "http://localhost:3000/" + uniqueId
	responseBody.ExpireTime = requestBody.ExpireTime

	redisDb.Decr(store.Context.Ctx, context.IP())

	rateLimiting, _ := redisDb.Get(store.Context.Ctx, context.IP()).Result()
	responseBody.XRateRemaining, _ = strconv.Atoi(rateLimiting)

	remainingTime, _ := redisDb.TTL(store.Context.Ctx, context.IP()).Result()
	responseBody.XRateLimitReset = int((remainingTime / time.Nanosecond / time.Minute))

	return context.Status(fiber.StatusOK).JSON(responseBody)
}
