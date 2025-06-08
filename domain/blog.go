package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Blog struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string        `bson:"title" json:"title"`
	Excerpt     string        `bson:"excerpt" json:"excerpt"`
	Content     string        `bson:"content" json:"content"`
	Author      string        `bson:"author" json:"author"`
	AuthorImage string        `bson:"authorImage" json:"authorImage"`
	Date        string        `bson:"date" json:"date"`
	ReadTime    string        `bson:"readTime" json:"readTime"`
	Category    string        `bson:"category" json:"category"`
	Tags        []string      `bson:"tags" json:"tags"`
	Image       string        `bson:"image" json:"image"`
	Featured    bool          `bson:"featured" json:"featured"`
}