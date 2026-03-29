package service

import (
	"errors"
	"math"

	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/repository"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionService interface {
	CreateSubscription(req *models.CreateSubscriptionRequest) (*models.Subscription, error)
	GetAllSubscriptions(query models.SubscriptionQuery) (*models.PaginatedSubscriptions, error)
	GetSubscriptionByID(id uint) (*models.Subscription, error)
	UpdateSubscription(id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error)
	DeleteSubscription(id uint) error
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	sub := &models.Subscription{
		Name:         req.Name,
		Vendor:       req.Vendor,
		Plan:         req.Plan,
		BillingCycle: req.BillingCycle,
		CostPerSeat:  req.CostPerSeat,
		Seats:        req.Seats,
		CurrentUsers: req.CurrentUsers,
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

func (s *subscriptionService) GetAllSubscriptions(query models.SubscriptionQuery) (*models.PaginatedSubscriptions, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	subs, total, err := s.repo.FindAll(query)
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

func (s *subscriptionService) GetSubscriptionByID(id uint) (*models.Subscription, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, ErrSubscriptionNotFound
	}
	return sub, nil
}

func (s *subscriptionService) UpdateSubscription(id uint, req *models.UpdateSubscriptionRequest) (*models.Subscription, error) {
	sub, err := s.repo.FindByID(id)
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
	if req.Plan != nil {
		sub.Plan = *req.Plan
	}
	if req.BillingCycle != nil {
		sub.BillingCycle = *req.BillingCycle
	}
	if req.CostPerSeat != nil {
		sub.CostPerSeat = *req.CostPerSeat
	}
	if req.Seats != nil {
		sub.Seats = *req.Seats
	}
	if req.CurrentUsers != nil {
		sub.CurrentUsers = *req.CurrentUsers
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

func (s *subscriptionService) DeleteSubscription(id uint) error {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if sub == nil {
		return ErrSubscriptionNotFound
	}

	return s.repo.Delete(id)
}
