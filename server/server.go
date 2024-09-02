package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"tincho.dev/rest-ws/database"
	"tincho.dev/rest-ws/repositories"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	if config.DatabaseURL == "" {
		return nil, errors.New("database url is required")
	}

	return &Broker{
		config: config,
		router: mux.NewRouter(),
	}, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	binder(b, b.router)
	repo, err := database.NewPostgres(b.config.DatabaseURL)

	if err != nil {
		log.Fatal("Error creating database: ", err)
	}

	repositories.SetUserRepository(repo)
	log.Println("Server is running on port", b.config.Port)

	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
