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
	"github.com/cihan-ozata/backend-path/internal/handler"
	"github.com/cihan-ozata/backend-path/internal/middleware"
	"github.com/cihan-ozata/backend-path/internal/repository"
	"github.com/cihan-ozata/backend-path/internal/service"
	"github.com/cihan-ozata/backend-path/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App, uygulama yapısını temsil eder
type App struct {
	config   *config.Config
	logger   logger.Logger
	server   *http.Server
	database *repository.Database
	router   *gin.Engine
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

	// App instance oluştur
	app := &App{
		config:   cfg,
		logger:   log,
		database: db,
	}

	// Gin router oluştur
	router := gin.New()

	// Jaeger tracer'ı başlat
	if err := middleware.InitTracer("backend-api", "jaeger:14268"); err != nil {
		log.Warn("Jaeger tracer başlatılamadı", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Middleware'leri ekle
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.MonitorMiddleware())
	router.Use(middleware.TracingMiddleware())

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Monitoring endpoint'leri
	router.GET("/health", app.healthCheck)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/ready", app.readyCheck)

	// Pprof endpoint'leri (development için)
	if cfg.Logging.Level == "debug" {
		pprof.Register(router)
	}

	// Services ve handlers oluştur
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	transactionRepo := repository.NewTransactionRepository(db.DB)
	walletRepo := repository.NewWalletRepository(db.DB)
	transactionService := service.NewTransactionService(transactionRepo, walletRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	paymentRepo := repository.NewPaymentRepository(db.DB)
	paymentService := service.NewPaymentService(paymentRepo, walletRepo, transactionRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// API routes
	api := router.Group("/api/v1")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/", userHandler.ListUsers)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("/", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.GET("/", transactionHandler.ListTransactions)
		}

		// Wallet routes
		wallets := api.Group("/wallets")
		{
			wallets.POST("/", walletHandler.CreateWallet)
			wallets.GET("/:id", walletHandler.GetWallet)
			wallets.PUT("/:id", walletHandler.UpdateWallet)
			wallets.GET("/", walletHandler.ListWallets)
		}

		// Payment routes
		payments := api.Group("/payments")
		{
			payments.POST("/", paymentHandler.CreatePayment)
			payments.GET("/:id", paymentHandler.GetPayment)
			payments.GET("/", paymentHandler.ListPayments)
		}
	}

	// HTTP sunucusu oluştur
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	app.server = server
	app.router = router

	return app, nil
}

// healthCheck endpoint'i
func (a *App) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "backend-api",
	})
}

// readyCheck endpoint'i
func (a *App) readyCheck(c *gin.Context) {
	// Veritabanı bağlantısını kontrol et
	if err := a.database.DB.Raw("SELECT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"timestamp": time.Now().Format(time.RFC3339),
			"error":     "database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "backend-api",
	})
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