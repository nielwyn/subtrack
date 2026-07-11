package service

import (
	"errors"
	"math"
	"strings"

	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/repository"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionService interface {
	CreateSubscription(userID uint, req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetAllSubscriptions(userID uint, query models.SubscriptionQuery) (*models.PaginatedSubscriptions, error)
	GetSubscriptionByID(userID, id uint) (*models.Subscription, error)
	UpdateSubscription(userID, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	DeleteSubscription(userID, id uint) error
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(userID uint, req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if currency == "" {
		currency = "USD"
	}

	sub := &models.Subscription{
		UserID:       userID,
		Name:         req.Name,
		Vendor:       req.Vendor,
		Platform:     req.Platform,
		BillingCycle: req.BillingCycle,
		Amount:       req.Amount,
		Currency:     currency,
		Category:     req.Category,
		AutoRenews:   req.AutoRenews,
		RenewalDate:  req.RenewalDate,
		Status:       status,
		Notes:        req.Notes,
	}

	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) GetAllSubscriptions(userID uint, query models.SubscriptionQuery) (*models.PaginatedSubscriptions, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	subs, total, err := s.repo.FindAll(userID, query)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedSubscriptions{
		Subscriptions: subs,
		Total:         total,
		Page:          query.Page,
		Limit:         query.Limit,
		TotalPages:    int(math.Ceil(float64(total) / float64(query.Limit))),
	}, nil
}

func (s *subscriptionService) GetSubscriptionByID(userID, id uint) (*models.Subscription, error) {
	sub, err := s.repo.FindByID(userID, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, ErrSubscriptionNotFound
	}
	return sub, nil
}

func (s *subscriptionService) UpdateSubscription(userID, id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	sub, err := s.repo.FindByID(userID, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, ErrSubscriptionNotFound
	}

	if req.Name != nil {
		sub.Name = *req.Name
	}
	if req.Vendor != nil {
		sub.Vendor = *req.Vendor
	}
	if req.Platform != nil {
		sub.Platform = *req.Platform
	}
	if req.BillingCycle != nil {
		sub.BillingCycle = *req.BillingCycle
	}
	if req.Amount != nil {
		sub.Amount = *req.Amount
	}
	if req.Currency != nil {
		sub.Currency = strings.ToUpper(strings.TrimSpace(*req.Currency))
	}
	if req.Category != nil {
		sub.Category = *req.Category
	}
	if req.AutoRenews != nil {
		sub.AutoRenews = *req.AutoRenews
	}
	if req.RenewalDate != nil {
		sub.RenewalDate = req.RenewalDate
	}
	if req.Status != nil {
		sub.Status = *req.Status
	}
	if req.Notes != nil {
		sub.Notes = *req.Notes
	}

	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) DeleteSubscription(userID, id uint) error {
	sub, err := s.repo.FindByID(userID, id)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionNotFound
	}

	return s.repo.Delete(userID, id)
}
