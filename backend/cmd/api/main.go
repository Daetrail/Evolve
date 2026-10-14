package main

import (
	"log"
	"net/http"

	"evolve-server/internal/config"
	"evolve-server/internal/database"
	"evolve-server/internal/middleware"
	"evolve-server/internal/user"

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
	
	// create stores
	userStore := user.NewUserStore(db)
	
	// create services
	userService := user.NewUserService(userStore, config.JWTSecret)
	
	// create handlers
	userHandler := user.NewUserHandler(userService)
	
	// create middleware
	authMiddleware := middleware.AuthMiddleware(config.JWTSecret)
	
	
	// routes
	mux := http.NewServeMux()
	
	// public
	// user
	mux.HandleFunc("POST /auth/sign-up", userHandler.SignUp)
	mux.HandleFunc("POST /auth/sign-in", userHandler.SignIn)
	
	// protected
	mux.Handle("GET /auth/verify", authMiddleware(http.HandlerFunc(userHandler.Verify)))
	mux.Handle("GET /user/me", authMiddleware(http.HandlerFunc(userHandler.Me)))
	
	log.Printf("Listening on :%s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, mux))
}
