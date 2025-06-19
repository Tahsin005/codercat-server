package service

import (
	"context"

	"github.com/tahsin005/codercat-server/domain"
	"github.com/tahsin005/codercat-server/repository"
)

type SubscriberService interface {
	CreateSubscriber(ctx context.Context, subscriber *domain.Subscriber) error
}

type subscriberService struct {
	repo repository.SubscriberRepository
}

func NewSubscriberService(repo repository.SubscriberRepository) SubscriberService {
	return &subscriberService{repo: repo}
}

func (s *subscriberService) CreateSubscriber(ctx context.Context, subscriber *domain.Subscriber) error {
	return s.repo.CreateSubscriber(ctx, subscriber)
}