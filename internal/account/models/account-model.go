package models

import (
	"time"
)

type Account struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Owner     string    `json:"owner"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
