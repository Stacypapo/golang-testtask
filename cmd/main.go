package main

import (
	"golangTestTask/configs"
	"golangTestTask/internal/handler"
	"golangTestTask/internal/repository"
	"golangTestTask/internal/service"
	"log"
	"net/http"

	_ "golangTestTask/docs"
)

// @title Payment System API
// @version 1.0
// @description API для управления транзакциями и кошельками
// @host localhost:8080
// @BasePath /
func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := repository.NewPostgresDB(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatal(err)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	services.BaseWallets(10, 100)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.InitRoutes()))
}
