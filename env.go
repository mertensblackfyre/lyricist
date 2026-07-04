package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	CLIENT_ID     string
	CLIENT_SECRET string
)

func init() {
	err := godotenv.Load(".env")
	
	if err != nil {
		log.Fatalln("Error loading .env file", "error", err)
		return
	}

	CLIENT_ID = os.Getenv("SPOTIFY_ID")
	CLIENT_SECRET = os.Getenv("SPOTIFY_SECRET")

}
