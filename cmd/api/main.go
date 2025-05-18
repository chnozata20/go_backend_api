package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cihan-ozata/backend-path/config"
	"github.com/cihan-ozata/backend-path/internal/domain"
	"github.com/cihan-ozata/backend-path/internal/handler"
	"github.com/cihan-ozata/backend-path/internal/middleware"
	"github.com/cihan-ozata/backend-path/internal/repository"
	"github.com/cihan-ozata/backend-path/internal/service"
	"github.com/cihan-ozata/backend-path/internal/worker"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Reset database tables
func resetDatabase(db *gorm.DB) error {
	tables := []string{"users", "transactions", "wallets", "payments"}
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// Config yükle
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config yüklenirken hata oluştu:", err)
	}

	// Database connection
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	// Veritabanını sıfırlamak için bu satırı ekleyin
	if err := resetDatabase(db); err != nil {
		log.Fatal("Failed to reset database:", err)
	}

	// Auto migrate database schema
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Transaction{},
		&domain.Wallet{},
		&domain.Payment{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	walletService := service.NewWalletService(walletRepo)
	transactionService := service.NewTransactionService(transactionRepo, walletRepo)
	paymentService := service.NewPaymentService(paymentRepo, walletRepo, transactionRepo)

	// Initialize workers
	transactionWorker := worker.NewTransactionWorker(5, transactionService)
	taskProcessor := worker.NewTaskProcessor(3)

	// Start workers
	ctx := context.Background()
	transactionWorker.Start(ctx)
	taskProcessor.Start(ctx)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	walletHandler := handler.NewWalletHandler(walletService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.SecurityMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RateLimitMiddleware(100, 200)) // 100 requests per second, burst of 200
	router.Use(middleware.MonitorMiddleware())

	// Public routes
	router.POST("/users/register", middleware.ValidationMiddleware(&domain.User{}), userHandler.Register)
	router.POST("/users/login", middleware.ValidationMiddleware(&domain.LoginRequest{}), userHandler.Login)

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// User routes
	api.GET("/users/:id", userHandler.GetUser)
	api.PUT("/users/:id", middleware.ValidationMiddleware(&domain.User{}), userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	// Wallet routes
	api.POST("/wallets", middleware.ValidationMiddleware(&domain.Wallet{}), walletHandler.CreateWallet)
	api.GET("/wallets/user/:user_id", walletHandler.GetUserWallet)
	api.GET("/wallets/user/:user_id/balance", walletHandler.GetWalletBalance)

	// Transaction routes
	api.POST("/transactions", middleware.ValidationMiddleware(&domain.Transaction{}), transactionHandler.CreateTransaction)
	api.GET("/transactions/:id", transactionHandler.GetTransaction)
	api.GET("/transactions/user/:user_id", transactionHandler.GetUserTransactions)
	api.POST("/transactions/:id/process", transactionHandler.ProcessTransaction)

	// Payment routes
	api.POST("/payments", middleware.ValidationMiddleware(&domain.Payment{}), paymentHandler.CreatePayment)
	api.GET("/payments/:id", paymentHandler.GetPayment)
	api.GET("/payments/user/:user_id", paymentHandler.GetUserPayments)
	api.POST("/payments/:id/process", paymentHandler.ProcessPayment)
	api.POST("/payments/:id/refund", paymentHandler.RefundPayment)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.RoleMiddleware(domain.UserRoleAdmin))
	admin.GET("/users", userHandler.GetAllUsers)
	admin.GET("/transactions", transactionHandler.GetAllTransactions)
	admin.GET("/payments", paymentHandler.GetAllPayments)

	// Metrics endpoint
	admin.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, middleware.GetMetrics())
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 