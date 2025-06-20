package config

import (
	"os"
)

type Config struct {
	MongoURI                 string
	MongoDBName              string
	MongoCollNameBlogs       string
	MongoCollNameSubscribers string
	Port                     string
	SMTPEmail                string
	SMTPPassword             string
	SMTPHost                 string
	SMTPPort                 string
	BaseURL                  string
}

func LoadConfig() (*Config, error) {
	return &Config{
		MongoURI:                 getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:              getEnv("MONGO_DB_NAME", "codercat"),
		MongoCollNameBlogs:       getEnv("MONGO_COLLECTION_NAME_BLOG", "blogs"),
		MongoCollNameSubscribers: getEnv("MONGO_COLLECTION_NAME_SUBSCRIBERS", "subscribers"),
		Port:                     getEnv("PORT", "8080"),
		SMTPEmail:                getEnv("SMTP_EMAIL", ""),
		SMTPPassword:             getEnv("SMTP_PASSWORD", ""),
		SMTPHost:                 getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:                 getEnv("SMTP_PORT", "587"),
		BaseURL:                  getEnv("BASE_URL", "https://codercat-server.onrender.com"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
