package config

import (
	"log"
	"os"
)

var JWTKey []byte

func init() {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Fatal("JWT_KEY environment variable not set")
	}
	JWTKey = []byte(key)
}
