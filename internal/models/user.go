package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole, kullanıcı rollerini temsil eder
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

// User, kullanıcı modelini temsil eder
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // JSON çıktısında gösterilmez
	Role      UserRole       `gorm:"size:20;default:user" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete için
}

// TableName, tablonun adını belirtir
func (User) TableName() string {
	return "users"
} 