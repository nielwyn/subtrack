package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	Name         string         `gorm:"not null" json:"name"`
	Vendor       string         `gorm:"not null" json:"vendor"`
	Platform     string         `gorm:"not null" json:"platform"`
	BillingCycle string         `gorm:"not null;default:'monthly'" json:"billing_cycle"` // monthly, annual
	Amount       float64        `gorm:"not null;default:0" json:"amount"`
	Currency     string         `gorm:"not null;default:'USD'" json:"currency"`
	Category     string         `json:"category"`
	AutoRenews   bool           `gorm:"not null;default:true" json:"auto_renews"`
	RenewalDate  *time.Time     `json:"renewal_date"`
	Status       string         `gorm:"not null;default:'active'" json:"status"` // active, trial, cancelled
	Notes        string         `json:"notes"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type CreateSubscriptionRequest struct {
	Name         string     `json:"name" binding:"required,min=1,max=200"`
	Vendor       string     `json:"vendor" binding:"required,min=1,max=100"`
	Platform     string     `json:"platform" binding:"required,min=1,max=100"`
	BillingCycle string     `json:"billing_cycle" binding:"required,oneof=monthly annual"`
	Amount       float64    `json:"amount" binding:"non_negative"`
	Currency     string     `json:"currency" binding:"omitempty,len=3"`
	Category     string     `json:"category" binding:"max=100"`
	AutoRenews   bool       `json:"auto_renews"`
	RenewalDate  *time.Time `json:"renewal_date"`
	Status       string     `json:"status" binding:"omitempty,oneof=active trial cancelled"`
	Notes        string     `json:"notes" binding:"max=2000"`
}

type UpdateSubscriptionRequest struct {
	Name         *string    `json:"name" binding:"omitempty,min=1,max=200"`
	Vendor       *string    `json:"vendor" binding:"omitempty,min=1,max=100"`
	Platform     *string    `json:"platform" binding:"omitempty,min=1,max=100"`
	BillingCycle *string    `json:"billing_cycle" binding:"omitempty,oneof=monthly annual"`
	Amount       *float64   `json:"amount" binding:"omitempty,non_negative"`
	Currency     *string    `json:"currency" binding:"omitempty,len=3"`
	Category     *string    `json:"category" binding:"omitempty,max=100"`
	AutoRenews   *bool      `json:"auto_renews"`
	RenewalDate  *time.Time `json:"renewal_date"`
	Status       *string    `json:"status" binding:"omitempty,oneof=active trial cancelled"`
	Notes        *string    `json:"notes" binding:"omitempty,max=2000"`
}

type SubscriptionQuery struct {
	Page         int     `form:"page"`
	Limit        int     `form:"limit"`
	Vendor       string  `form:"vendor"`
	Platform     string  `form:"platform"`
	Category     string  `form:"category"`
	Status       string  `form:"status"`
	BillingCycle string  `form:"billing_cycle"`
	MinAmount    float64 `form:"min_amount"`
	MaxAmount    float64 `form:"max_amount"`
	RenewingSoon int     `form:"renewing_soon"` // renewing within N days
}

type PaginatedSubscriptions struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int64          `json:"total"`
	Page          int            `json:"page"`
	Limit         int            `json:"limit"`
	TotalPages    int            `json:"total_pages"`
}
