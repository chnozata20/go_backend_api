package models

import (
	"time"

	"gorm.io/gorm"
)

// Balance, kullanıcı bakiye modelini temsil eder
type Balance struct {
	UserID       uint      `gorm:"primaryKey" json:"user_id"`
	Amount       float64   `gorm:"type:decimal(20,2);not null;default:0" json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`

	// İlişkiler
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName, tablonun adını belirtir
func (Balance) TableName() string {
	return "balances"
}

// BeforeCreate, kayıt oluşturulmadan önce çalışır
func (b *Balance) BeforeCreate(tx *gorm.DB) error {
	b.LastUpdatedAt = time.Now()
	return nil
}

// BeforeUpdate, kayıt güncellenmeden önce çalışır
func (b *Balance) BeforeUpdate(tx *gorm.DB) error {
	b.LastUpdatedAt = time.Now()
	return nil
} 