package main

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"tincho.dev/rest-ws/handlers"
	"tincho.dev/rest-ws/middlewares"
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
	r.Use(middlewares.AuthMiddleware(s))

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods("GET")

	// Auth routes
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods("POST")
	r.HandleFunc("/signin", handlers.SignInHandler(s)).Methods("POST")

	// User routes
	r.HandleFunc("/users", handlers.FindAllUsersHandler(s)).Methods("GET")
	r.HandleFunc("/users/list", handlers.ListUsersHandler(s)).Methods("GET")
	r.HandleFunc("/users/me", handlers.MeHandler(s)).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.FindOneUserHandler(s)).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.UpdateUserHandler(s)).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.DeleteUserHandler(s)).Methods("DELETE")

	// Post routes
	r.HandleFunc("/posts", handlers.CreatePostHandler(s)).Methods("POST")
	r.HandleFunc("/posts/{id}", handlers.FindOnePostHandler(s)).Methods("GET")
	r.HandleFunc("/posts", handlers.FindAllPostsHandler(s)).Methods("GET")
	r.HandleFunc("/posts/{id}", handlers.UpdateOnePostHandler(s)).Methods("PUT")
	r.HandleFunc("/posts/{id}", handlers.DeleteOnePostHandler(s)).Methods("DELETE")
}
