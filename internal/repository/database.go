package repository

import (
	"fmt"

	"github.com/cihan-ozata/go-backend-api/config"
	"github.com/cihan-ozata/go-backend-api/internal/models"
	"github.com/cihan-ozata/go-backend-api/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Database, veritabanı bağlantısını ve işlemlerini yönetir
type Database struct {
	DB     *gorm.DB
	Config *config.Config
	Logger logger.Logger
}

// NewDatabase, yeni bir veritabanı bağlantısı oluşturur
func NewDatabase(cfg *config.Config, log logger.Logger) (*Database, error) {
	log.Info("Veritabanına bağlanılıyor", map[string]interface{}{
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"name": cfg.Database.DBName,
	})

	// DSN (Data Source Name) oluştur
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// GORM yapılandırması
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Tablo adlarını tekil kullan
		},
		Logger: NewGormLogger(log), // Özel logger
	}

	// Veritabanına bağlan
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("veritabanına bağlanılamadı: %w", err)
	}

	log.Info("Veritabanına başarıyla bağlanıldı", nil)

	return &Database{
		DB:     db,
		Config: cfg,
		Logger: log,
	}, nil
}

// Migrate, veritabanı şemasını oluşturur veya günceller
func (d *Database) Migrate() error {
	d.Logger.Info("Veritabanı migrasyonu başlatılıyor", nil)

	// Tüm modelleri migrate et
	if err := d.DB.AutoMigrate(models.Models...); err != nil {
		return fmt.Errorf("veritabanı migrasyonu başarısız: %w", err)
	}

	d.Logger.Info("Veritabanı migrasyonu başarıyla tamamlandı", nil)
	return nil
}

// Close, veritabanı bağlantısını kapatır
func (d *Database) Close() error {
	d.Logger.Info("Veritabanı bağlantısı kapatılıyor", nil)

	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("veritabanı bağlantısı alınamadı: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("veritabanı bağlantısı kapatılamadı: %w", err)
	}

	d.Logger.Info("Veritabanı bağlantısı başarıyla kapatıldı", nil)
	return nil
} 