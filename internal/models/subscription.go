package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Vendor       string         `gorm:"not null" json:"vendor"`
	Plan         string         `json:"plan"`
	BillingCycle string         `gorm:"not null;default:'monthly'" json:"billing_cycle"` // monthly, annual
	CostPerSeat  float64        `gorm:"not null;default:0" json:"cost_per_seat"`
	Seats        int            `gorm:"not null;default:1" json:"seats"`
	CurrentUsers int            `gorm:"not null;default:0" json:"current_users"`
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
	Plan         string     `json:"plan" binding:"max=100"`
	BillingCycle string     `json:"billing_cycle" binding:"required,oneof=monthly annual"`
	CostPerSeat  float64    `json:"cost_per_seat" binding:"non_negative"`
	Seats        int        `json:"seats" binding:"positive"`
	CurrentUsers int        `json:"current_users" binding:"non_negative"`
	AutoRenews   bool       `json:"auto_renews"`
	RenewalDate  *time.Time `json:"renewal_date"`
	Status       string     `json:"status" binding:"omitempty,oneof=active trial cancelled"`
	Notes        string     `json:"notes" binding:"max=2000"`
}

type UpdateSubscriptionRequest struct {
	Name         *string    `json:"name" binding:"omitempty,min=1,max=200"`
	Vendor       *string    `json:"vendor" binding:"omitempty,min=1,max=100"`
	Plan         *string    `json:"plan" binding:"omitempty,max=100"`
	BillingCycle *string    `json:"billing_cycle" binding:"omitempty,oneof=monthly annual"`
	CostPerSeat  *float64   `json:"cost_per_seat" binding:"omitempty,non_negative"`
	Seats        *int       `json:"seats" binding:"omitempty,positive"`
	CurrentUsers *int       `json:"current_users" binding:"omitempty,non_negative"`
	AutoRenews   *bool      `json:"auto_renews"`
	RenewalDate  *time.Time `json:"renewal_date"`
	Status       *string    `json:"status" binding:"omitempty,oneof=active trial cancelled"`
	Notes        *string    `json:"notes" binding:"omitempty,max=2000"`
}

type SubscriptionQuery struct {
	Page         int     `form:"page"`
	Limit        int     `form:"limit"`
	Vendor       string  `form:"vendor"`
	Status       string  `form:"status"`
	BillingCycle string  `form:"billing_cycle"`
	MinCost      float64 `form:"min_cost"`
	MaxCost      float64 `form:"max_cost"`
	LowCapacity  int     `form:"low_capacity"`   // subscriptions where (seats - current_users) <= this
	RenewingSoon int     `form:"renewing_soon"`  // renewing within N days
}

type PaginatedSubscriptions struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int64          `json:"total"`
	Page          int            `json:"page"`
	Limit         int            `json:"limit"`
	TotalPages    int            `json:"total_pages"`
}
