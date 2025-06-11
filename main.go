package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tahsin005/codercat-server/config"
	"github.com/tahsin005/codercat-server/database"
	"github.com/tahsin005/codercat-server/handler"
	"github.com/tahsin005/codercat-server/repository"
	"github.com/tahsin005/codercat-server/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Failed to load .env file: %v", err)
	}
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB Atlas: %v", err)
	}
	defer db.Disconnect()

	blogRepo := repository.NewBlogRepository(db, cfg)
	blogService := service.NewBlogService(blogRepo)
	blogHandler := handler.NewBlogHandler(blogService)

	router := mux.NewRouter()
	blogHandler.RegisterRoutes(router)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":" + cfg.Port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}