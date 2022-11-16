package server

import (
	"context"
	"errors"
	"francocorrea/go/rest-ws/database"
	"francocorrea/go/rest-ws/repositories"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port string
	JWTSecret string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

//Return Config
func (broker *Broker) Config() *Config {
	return broker.config
}

//Boot Server
func (broker *Broker) Start(binder func(server Server, router *mux.Router)) {
	broker.router = mux.NewRouter()
	binder(broker, broker.router)

	//Conection Database
	repo, err := database.NewPostgresRepository(broker.config.DatabaseUrl)

	if err != nil {
		log.Fatal(err)
	}

	repositories.SetRepository(repo)
	
	log.Println("Starting server on port:", broker.Config().Port)

	if err := http.ListenAndServe(broker.config.Port, broker.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Server Constructor
func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("The Port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("The JWTSecret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("The DatabaseUrl is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	return broker, nil
}