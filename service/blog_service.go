package service

import (
	"context"
	"fmt"

	"github.com/tahsin005/codercat-server/domain"
	"github.com/tahsin005/codercat-server/repository"
	"github.com/tahsin005/codercat-server/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BlogService interface {
	CreateBlog(ctx context.Context, blog *domain.Blog) error
	GetBlogByID(ctx context.Context, id string) (*domain.Blog, error)
	UpdateBlog(ctx context.Context, id string, blog *domain.Blog) error
	DeleteBlog(ctx context.Context, id string) error
	GetAllBlogs(ctx context.Context) ([]*domain.Blog, error)
	GetFeaturedBlogs(ctx context.Context) ([]*domain.Blog, error)
	GetRecentBlogs(ctx context.Context, limit int) ([]*domain.Blog, error)
	GetBlogsByCategory(ctx context.Context, category string) ([]*domain.Blog, error)
	SearchBlogs(ctx context.Context, query string) ([]*domain.Blog, error)
	GetRelatedBlogs(ctx context.Context, id string, limit int) ([]*domain.Blog, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetPopularCategories(ctx context.Context, limit int) ([]string, error)
}

type blogService struct {
	repo              repository.BlogRepository
	subscriberService SubscriberService
	emailConfig       utils.EmailConfig
	templateService   TemplateService
	baseURL           string
}

func NewBlogService(repo repository.BlogRepository, subscriberService SubscriberService, emailConfig utils.EmailConfig, templateService TemplateService, baseURL string) BlogService {
	return &blogService{
		repo:              repo,
		subscriberService: subscriberService,
		emailConfig:       emailConfig,
		templateService:   templateService,
		baseURL:           baseURL,
	}
}

func (s *blogService) CreateBlog(ctx context.Context, blog *domain.Blog) error {
	if err := s.repo.Create(ctx, blog); err != nil {
		return err
	}

	// Notify all subscribers
	subscribers, err := s.subscriberService.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(subscribers) == 0 {
		return nil
	}

	var emails []string
	for _, sub := range subscribers {
		emails = append(emails, sub.Email)
	}

	// Prepare email data
	emailData := domain.EmailData{
		Title:          blog.Title,
		Excerpt:        blog.Excerpt,
		Author:         blog.Author,
		Category:       blog.Category,
		ReadTime:       blog.ReadTime,
		Tags:           blog.Tags,
		BlogURL:        fmt.Sprintf("%s/blogs/%s", s.baseURL, blog.ID.Hex()),
	}

	// Render HTML template
	htmlBody, err := s.templateService.RenderEmailTemplate("new_blog", emailData)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("🚀 New Blog Post: %s", blog.Title)

	// Send HTML email in background
	go utils.SendHTMLEmail(s.emailConfig, emails, subject, htmlBody)

	return nil
}

func (s *blogService) GetBlogByID(ctx context.Context, id string) (*domain.Blog, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, oid)
}

func (s *blogService) UpdateBlog(ctx context.Context, id string, blog *domain.Blog) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, oid, blog)
}

func (s *blogService) DeleteBlog(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, oid)
}

func (s *blogService) GetAllBlogs(ctx context.Context) ([]*domain.Blog, error) {
	return s.repo.FindAll(ctx)
}

func (s *blogService) GetFeaturedBlogs(ctx context.Context) ([]*domain.Blog, error) {
	return s.repo.FindFeatured(ctx)
}

func (s *blogService) GetRecentBlogs(ctx context.Context, limit int) ([]*domain.Blog, error) {
	return s.repo.FindRecent(ctx, limit)
}

func (s *blogService) GetBlogsByCategory(ctx context.Context, category string) ([]*domain.Blog, error) {
	return s.repo.FindByCategory(ctx, category)
}

func (s *blogService) SearchBlogs(ctx context.Context, query string) ([]*domain.Blog, error) {
	return s.repo.Search(ctx, query)
}

func (s *blogService) GetRelatedBlogs(ctx context.Context, id string, limit int) ([]*domain.Blog, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindRelated(ctx, oid, limit)
}

func (s *blogService) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories(ctx)
}

func (s *blogService) GetPopularCategories(ctx context.Context, limit int) ([]string, error) {
	return s.repo.GetPopularCategories(ctx, limit)
}
