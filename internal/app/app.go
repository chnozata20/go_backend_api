package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cihan-ozata/backend-path/config"
	"github.com/cihan-ozata/backend-path/internal/repository"
	"github.com/cihan-ozata/backend-path/pkg/logger"
)

// App, uygulama yapısını temsil eder
type App struct {
	config   *config.Config
	logger   logger.Logger
	server   *http.Server
	database *repository.Database
}

// New, yeni bir uygulama örneği oluşturur
func New() (*App, error) {
	// Yapılandırmayı yükle
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config yüklenirken hata oluştu: %w", err)
	}

	// Logger oluştur
	log := logger.New(cfg.Logging.Level, true)

	// Veritabanı bağlantısı oluştur
	db, err := repository.NewDatabase(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("veritabanı bağlantısı oluşturulamadı: %w", err)
	}

	// HTTP sunucusu oluştur
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      http.DefaultServeMux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &App{
		config:   cfg,
		logger:   log,
		server:   server,
		database: db,
	}, nil
}

// Run, uygulamayı başlatır
func (a *App) Run() error {
	a.logger.Info("Uygulama başlatılıyor", map[string]interface{}{
		"port": a.config.Server.Port,
	})

	// Veritabanı migrasyonlarını çalıştır
	if err := a.database.Migrate(); err != nil {
		return fmt.Errorf("veritabanı migrasyonu başarısız: %w", err)
	}

	// HTTP sunucusunu ayrı bir goroutine'de başlat
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("HTTP sunucusu başlatılamadı", err, nil)
		}
	}()

	a.logger.Info("Uygulama başlatıldı", map[string]interface{}{
		"port": a.config.Server.Port,
	})

	return nil
}

// Shutdown, uygulamayı düzgün bir şekilde kapatır
func (a *App) Shutdown() error {
	a.logger.Info("Uygulama kapatılıyor", nil)

	// Graceful shutdown için context oluştur
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// HTTP sunucusunu kapat
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("sunucu kapatılırken hata oluştu: %w", err)
	}

	// Veritabanı bağlantısını kapat
	if err := a.database.Close(); err != nil {
		return fmt.Errorf("veritabanı bağlantısı kapatılırken hata oluştu: %w", err)
	}

	a.logger.Info("Uygulama başarıyla kapatıldı", nil)
	return nil
}

// WaitForShutdown, kapatma sinyallerini bekler
func (a *App) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	a.logger.Info("Kapatma sinyali alındı", map[string]interface{}{
		"signal": sig.String(),
	})

	if err := a.Shutdown(); err != nil {
		a.logger.Error("Uygulama kapatılırken hata oluştu", err, nil)
		os.Exit(1)
	}
} 