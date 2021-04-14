package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(name, defaultValue string) string {
	value := os.Getenv(name)

	if value != "" {
		return value
	}

	return defaultValue
}

func GetInt(name string, defaultValue int) int {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("WARNING: %v is not an int. Using default value of %v", name, defaultValue)
		return defaultValue
	}

	return n
}
