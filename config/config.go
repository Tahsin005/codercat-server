package config

import (
	"os"
)

type Config struct {
	MongoURI      string
	MongoDBName   string
	MongoCollNameBlogs string
	MongoCollNameSubscribers string
	Port          string
}

func LoadConfig() (*Config, error) {
	return &Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:   getEnv("MONGO_DB_NAME", "codercat"),
		MongoCollNameBlogs: getEnv("MONGO_COLLECTION_NAME_BLOG", "blogs"),
		MongoCollNameSubscribers: getEnv("MONGO_COLLECTION_NAME_SUBSCRIBERS", "subscribers"),
		Port:          getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}