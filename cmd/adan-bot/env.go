package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Getenv(key, def string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	log.Printf("Missing environment variable '%s', using default value", key)

	return def
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Cannot load envfile: %v", err)
	}
}
