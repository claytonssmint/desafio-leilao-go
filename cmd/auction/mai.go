package main

import (
	"context"
	"log"

	"github.com/claytonssmint/desafio-leilao-go/configuration/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	_, err := mongodb.NewMongoDBConnetion(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
