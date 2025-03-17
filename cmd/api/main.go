package main

import (
	"log"
	"os"

	"github.com/cihan-ozata/go-backend-api/internal/app"
)

func main() {
	// Uygulama örneği oluştur
	application, err := app.New()
	if err != nil {
		log.Fatalf("Uygulama başlatılamadı: %v", err)
		os.Exit(1)
	}

	// Uygulamayı başlat
	if err := application.Run(); err != nil {
		log.Fatalf("Uygulama çalıştırılamadı: %v", err)
		os.Exit(1)
	}

	// Kapatma sinyallerini bekle
	application.WaitForShutdown()
} 