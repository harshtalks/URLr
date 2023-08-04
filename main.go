package main

import (
	_ "app/docs"
	"app/routers"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/joho/godotenv"
)

//	@title		URLr APIs
//	@version	1.0
//	@description.markdown
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Harsh Pareeek
//	@contact.url	https://www.github.com/harshtalks
//	@contact.email	harshpareek91@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/
func main() {

	// loading the env files here

	if envFileLoadError := godotenv.Load(); envFileLoadError != nil {
		fmt.Print("Error while reading env error")
		panic("No Env Found")
	}

	// creating the instance of a new go fiber

	app := fiber.New()

	// Adding Middlewares
	// 1. Middleware for the logging.
	app.Use(logger.New())

	// 2 API Docs
	app.Get("/docs/*", swagger.HandlerDefault)

	// API Handlers

	app.Post("/create", routers.Create)

	app.Get("/:tag", routers.Fetch)

	// listening on the Port 3000
	log.Fatal(app.Listen(":3000"))
}
