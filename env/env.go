package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

var REDIS_CACHE_DB, REDIS_FREQUENCY_DB int

var BASE_TTL = 30 * 60 * time.Second

func Init() {
    initRedisVariable("REDIS_CACHE_DB", &REDIS_CACHE_DB)
    initRedisVariable("REDIS_FREQUENCY_DB", &REDIS_FREQUENCY_DB)
}

func initRedisVariable(variableKey string, target *int) {
    dbStr := os.Getenv(variableKey)
    if dbStr == "" {
        *target = 0 // Default value, change it as needed
    } else {
        db, err := strconv.Atoi(dbStr)
        if err != nil {
            log.Fatalf("Failed to convert %s to integer: %v", variableKey, err)
        }
        *target = db
    }
}
