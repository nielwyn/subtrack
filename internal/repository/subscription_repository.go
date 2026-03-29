package repository

import (
	"errors"
	"time"

	"github.com/nielwyn/inventory-system/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub *models.Subscription) error
	FindAll(query models.SubscriptionQuery) ([]models.Subscription, int64, error)
	FindByID(id uint) (*models.Subscription, error)
	Update(sub *models.Subscription) error
	Delete(id uint) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(sub *models.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *subscriptionRepository) FindAll(query models.SubscriptionQuery) ([]models.Subscription, int64, error) {
	var subs []models.Subscription
	var total int64

	db := r.db.Model(&models.Subscription{})

	if query.Vendor != "" {
		db = db.Where("vendor = ?", query.Vendor)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.BillingCycle != "" {
		db = db.Where("billing_cycle = ?", query.BillingCycle)
	}
	if query.MinCost > 0 {
		db = db.Where("cost_per_seat >= ?", query.MinCost)
	}
	if query.MaxCost > 0 {
		db = db.Where("cost_per_seat <= ?", query.MaxCost)
	}
	if query.LowCapacity > 0 {
		db = db.Where("(seats - current_users) <= ?", query.LowCapacity)
	}
	if query.RenewingSoon > 0 {
		deadline := time.Now().AddDate(0, 0, query.RenewingSoon)
		db = db.Where("renewal_date IS NOT NULL AND renewal_date <= ?", deadline)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit
	err := db.Offset(offset).Limit(query.Limit).Find(&subs).Error
	return subs, total, err
}

func (r *subscriptionRepository) FindByID(id uint) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.First(&sub, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) Update(sub *models.Subscription) error {
	return r.db.Save(sub).Error
}

func (r *subscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subscription{}, id).Error
}
