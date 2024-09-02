package main

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"tincho.dev/rest-ws/handlers"
	"tincho.dev/rest-ws/server"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseURL: DATABASE_URL,
	})

	if err != nil {
		log.Fatal("Error creating server: ", err)
	}

	s.Start(binder)
}

func binder(s server.Server, r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods("GET")
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods("POST")
	r.HandleFunc("/signin", handlers.SignInHandler(s)).Methods("POST")
	r.HandleFunc("/users", handlers.FindAllUsersHandler(s)).Methods("GET")
}
