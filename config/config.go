package config

import (
	"log"
	"os"
)

// JWTKey is read from environment variable
var JWTKey []byte

func init() {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Fatal("JWT_KEY environment variable not set")
	}
	JWTKey = []byte(key)
}
