package main

import (
	"mkgame-go/web"
	"log"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	testEnvValue := os.Getenv("ENV_PARAM")
	log.Println("ENV : ", testEnvValue)
	log.Println("start")
	web.ServerOn()
	log.Println("done")
}
