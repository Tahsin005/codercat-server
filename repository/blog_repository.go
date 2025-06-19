package repository

import (
	"context"

	"github.com/tahsin005/codercat-server/config"
	"github.com/tahsin005/codercat-server/database"
	"github.com/tahsin005/codercat-server/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type BlogRepository interface {
	Create(ctx context.Context, blog *domain.Blog) error
	FindByID(ctx context.Context, id bson.ObjectID) (*domain.Blog, error)
	Update(ctx context.Context, id bson.ObjectID, blog *domain.Blog) error
	Delete(ctx context.Context, id bson.ObjectID) error
	FindAll(ctx context.Context) ([]*domain.Blog, error)
	FindFeatured(ctx context.Context) ([]*domain.Blog, error)
	FindRecent(ctx context.Context, limit int) ([]*domain.Blog, error)
	FindByCategory(ctx context.Context, category string) ([]*domain.Blog, error)
	Search(ctx context.Context, query string) ([]*domain.Blog, error)
	FindRelated(ctx context.Context, blogID bson.ObjectID, limit int) ([]*domain.Blog, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetPopularCategories(ctx context.Context, limit int) ([]string, error)
}

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(db *database.Database, cfg *config.Config) BlogRepository {
	return &blogRepository{
		collection: db.DB.Collection(cfg.MongoCollNameBlogs),
	}
}

func (r *blogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	blog.ID = bson.NewObjectID()
	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

func (r *blogRepository) FindByID(ctx context.Context, id bson.ObjectID) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&blog)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}
	return &blog, err
}

func (r *blogRepository) Update(ctx context.Context, id bson.ObjectID, blog *domain.Blog) error {
	blog.ID = id
	_, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: blog}})
	return err
}

func (r *blogRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}

func (r *blogRepository) FindAll(ctx context.Context) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

func (r *blogRepository) FindFeatured(ctx context.Context) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	cursor, err := r.collection.Find(ctx, bson.D{{Key: "featured", Value: true}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

func (r *blogRepository) FindRecent(ctx context.Context, limit int) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

func (r *blogRepository) FindByCategory(ctx context.Context, category string) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	filter := bson.D{}
	if category != "All" {
		filter = bson.D{{Key: "category", Value: category}}
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

func (r *blogRepository) Search(ctx context.Context, query string) ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "title", Value: bson.D{{Key: "$regex", Value: query}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "excerpt", Value: bson.D{{Key: "$regex", Value: query}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "content", Value: bson.D{{Key: "$regex", Value: query}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "tags", Value: bson.D{{Key: "$regex", Value: query}, {Key: "$options", Value: "i"}}}},
		}},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

func (r *blogRepository) FindRelated(ctx context.Context, blogID bson.ObjectID, limit int) ([]*domain.Blog, error) {
	currentBlog, err := r.FindByID(ctx, blogID)
	if err != nil {
		return nil, err
	}
	var blogs []*domain.Blog
	filter := bson.D{
		{Key: "_id", Value: bson.D{{Key: "$ne", Value: blogID}}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "category", Value: currentBlog.Category}},
			bson.D{{Key: "tags", Value: bson.D{{Key: "$in", Value: currentBlog.Tags}}}},
		}},
	}
	opts := options.Find().SetLimit(int64(limit))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, cursor.Err()
}

func (r *blogRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categoriesArr []string
	err := r.collection.Distinct(ctx, "category", bson.D{}).Decode(&categoriesArr)
	if err != nil {
		return nil, err
	}
	categories := append([]string{"All"}, categoriesArr...)
	return categories, nil
}

func (r *blogRepository) GetPopularCategories(ctx context.Context, limit int) ([]string, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$category"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Category string `bson:"_id"`
		Count    int    `bson:"count"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	categories := make([]string, len(results))
	for i, result := range results {
		categories[i] = result.Category
	}
	return categories, nil
}
