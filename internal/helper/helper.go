package helper

import (
	"log"
	"os"
)

func GetEnv(key string) (value string) {
	value = os.Getenv(key)
	if value == "" {
		log.Fatalf("%s must be set", key)
	}
	return value
}
