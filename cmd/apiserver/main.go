package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/UrcaDeLima/backend_golang_journal/internal/app/apiserver"
)

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
