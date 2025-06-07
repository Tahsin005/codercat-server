package database

import (
	"context"
	"log"

	"github.com/tahsin005/codercat-server/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	if cfg.MongoURI == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Printf("Failed to connect to MongoDB Atlas: %v", err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB Atlas: %v", err)
		return nil, err
	}

	db := client.Database(cfg.MongoDBName)

	return &Database{Client: client, DB: db}, nil
}

func (d *Database) Disconnect() {
	if err := d.Client.Disconnect(context.TODO()); err != nil {
		log.Printf("Error disconnecting from MongoDB Atlas: %v", err)
	}
}