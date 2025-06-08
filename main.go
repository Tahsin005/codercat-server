package main

import (
	"log"

	"github.com/tahsin005/codercat-server/config"
	"github.com/tahsin005/codercat-server/database"
	"github.com/tahsin005/codercat-server/repository"
	"github.com/tahsin005/codercat-server/service"
)

func main () {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.MongoURI == "" {
		log.Fatal("MONGO_URI is required for MongoDB Atlas connection")
	}

	db, err := database.NewDatabase(cfg)
	blogRepo := repository.NewBlogRepository(db, cfg)
	blogService := service.NewBlogService(blogRepo)
	log.Println(blogRepo, blogService)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB Atlas: %v", err)
	}

	if err == nil {
		log.Println("Connected to MongoDB Atlas")
	}
	defer db.Disconnect()
}