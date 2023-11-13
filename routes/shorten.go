package routes

import (
	"go-api/database"
	"go-api/env"
	"go-api/helpers"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Request struct {
	URL string `json:"url"`
	CustomShort string `json:"short"`
	Expiry time.Duration `json:"expiry"`
}


type Response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}


func ShortenUrl(c *fiber.Ctx) error {
	body := new(Request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
		})
	}


	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid URL",
		})
	}

	// implement rate limiting
	// everytime a user queries, check if the IP is already in database,
	// if yes, decrement the calls remaining by one, else add the IP to database
	// with expiry of `30mins`. So in this case the user will be able to send 10
	// requests every 30 minutes
	r2 := database.RedisClient(env.REDIS_FREQUENCY_DB)
	defer r2.Close()
	val, err := r2.Get(database.RedisCtx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.RedisCtx, c.IP(), os.Getenv("API_QUOTE"), env.BASE_TTL).Err()
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 && len(val) > 0 {
			limit, _ := r2.TTL(database.RedisCtx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// check for the domain error
	// users may abuse the shortener by shorting the domain `localhost:3000` itself
	// leading to a inifite loop, so don't accept the domain for shortening
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "haha... nice try",
		})
	}

	body.URL = helpers.EnforceHTTP(body.URL)
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.RedisClient(env.REDIS_CACHE_DB)
	defer r.Close()

	val, _ = r.Get(database.RedisCtx, id).Result()
	// check if the user provided short is already in use
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL short already in use",
		})
	}
	if body.Expiry == 0 {
		body.Expiry = 24 // default expiry of 24 hours
	}
	err = r.Set(database.RedisCtx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		log.Println("Failed -", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}
	// respond with the url, short, expiry in hours, calls remaining and time to reset
	resp := Response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	r2.Decr(database.RedisCtx, c.IP())
	val, _ = r2.Get(database.RedisCtx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.RedisCtx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}