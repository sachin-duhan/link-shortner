package routes

import (
	"go-api/database"
	"go-api/env"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")

	if len(url) < 3 {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid URL identifier.",
		})
	}

	// create db instance.
	db := database.RedisClient(env.REDIS_CACHE_DB)
	defer db.Close()

	value, err := db.Get(database.RedisCtx, url).Result()
	if err != redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No record exist in our DB.",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	rInr := database.RedisClient(env.REDIS_FREQUENCY_DB)
	defer rInr.Close()
	_ = rInr.Incr(database.RedisCtx, "counter")

	log.Fatalf("Found record for %s - {%s}", url, value)
	return c.Redirect(value, fiber.StatusMovedPermanently)
}
