package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Subscriber struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
}
