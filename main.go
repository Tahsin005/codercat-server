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
	"github.com/tahsin005/codercat-server/utils"
)

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:5173":        true,
			"https://codercat.vercel.app": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	log.Println("Connected to MongoDB Atlas")

	blogRepo := repository.NewBlogRepository(db, cfg)
	subscriberRepo := repository.NewSubscriberRepository(db, cfg)
	subscriberService := service.NewSubscriberService(subscriberRepo)
	templateService := service.NewTemplateService("templates")

	emailCfg := utils.EmailConfig{
		From:     cfg.SMTPEmail,
		Password: cfg.SMTPPassword,
		SMTPHost: cfg.SMTPHost,
		SMTPPort: cfg.SMTPPort,
	}

	blogService := service.NewBlogService(blogRepo, subscriberService, emailCfg, templateService, cfg.BaseURL)
	blogHandler := handler.NewBlogHandler(blogService)
	subscriberHandler := handler.NewSubscriberHandler(subscriberService)

	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(corsMiddleware)

	blogHandler.RegisterRoutes(router)
	subscriberHandler.RegisterRoutes(router)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
