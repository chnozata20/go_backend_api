package models

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog, denetim kaydı modelini temsil eder
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EntityType string    `gorm:"size:50;not null;index" json:"entity_type"`
	EntityID   uint      `gorm:"not null;index" json:"entity_id"`
	Action     string    `gorm:"size:50;not null;index" json:"action"`
	Details    string    `gorm:"type:text" json:"details"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     *uint     `gorm:"index" json:"user_id,omitempty"` // İşlemi yapan kullanıcı (opsiyonel)

	// İlişkiler
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName, tablonun adını belirtir
func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate, kayıt oluşturulmadan önce çalışır
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	return nil
} 