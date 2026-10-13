package main

import (
	"log"

	"evolve-server/internal/config"
	"evolve-server/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("no .env file found: %s", err.Error())
	}

	config, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	db, err := database.Open(config.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open DB: %s", err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("failed to close DB: ", err.Error())
		}
	}()

	log.Println("database connected successfully!")
	
	err = database.SetupDB(db)
	if err != nil {
		log.Fatalf("failed to setup DB tables: %s", err.Error())
	}
	
	log.Println("database tables setup successfully!")
}
