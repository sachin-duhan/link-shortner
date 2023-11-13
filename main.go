package main

import (
	"go-api/env"
	"go-api/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	app := fiber.New()

	app.Use(logger.New())
	env.Init()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
