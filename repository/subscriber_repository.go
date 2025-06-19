package repository

import (
	"context"

	"github.com/tahsin005/codercat-server/config"
	"github.com/tahsin005/codercat-server/database"
	"github.com/tahsin005/codercat-server/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubscriberRepository interface {
	CreateSubscriber(ctx context.Context, subscriber *domain.Subscriber) error
	GetAll(ctx context.Context) ([]*domain.Subscriber, error)
}

type subscriberRepository struct {
	collection *mongo.Collection
}

func NewSubscriberRepository(db *database.Database, cfg *config.Config) SubscriberRepository {
	return &subscriberRepository{
		collection: db.DB.Collection(cfg.MongoCollNameSubscribers),
	}
}

func (r *subscriberRepository) CreateSubscriber(ctx context.Context, subscriber *domain.Subscriber) error {
	subscriber.ID = bson.NewObjectID()
	_, err := r.collection.InsertOne(ctx, subscriber)
	return err
}

func (r *subscriberRepository) GetAll(ctx context.Context) ([]*domain.Subscriber, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscribers []*domain.Subscriber
	for cursor.Next(ctx) {
		var sub domain.Subscriber
		if err := cursor.Decode(&sub); err != nil {
			return nil, err
		}
		subscribers = append(subscribers, &sub)
	}
	return subscribers, nil
}
