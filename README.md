# Backend Path - Go Backend Projesi

Bu proje, Go ile geliştirilmiş modern bir backend API'sidir. PostgreSQL veritabanı, GORM ORM, Gin web framework, Zerolog loglama ve kapsamlı monitoring/observability özellikleri kullanır.

## 🚀 Özellikler

- **Clean Architecture**: Domain, Repository, Service ve Handler katmanları
- **PostgreSQL**: GORM ORM ile veritabanı yönetimi
- **Gin Framework**: Hızlı ve güvenilir HTTP router
- **JWT Authentication**: Güvenli kimlik doğrulama
- **Structured Logging**: Zerolog ile JSON formatında loglama
- **Monitoring & Observability**:
  - Prometheus metrics
  - Grafana dashboards
  - Jaeger distributed tracing
  - Health check endpoints
- **Docker Support**: Multi-stage Dockerfile ve Docker Compose
- **Environment Configuration**: .env dosyası ile yapılandırma

## 📋 Gereksinimler

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose
- Redis (opsiyonel, caching için)

## 🛠️ Kurulum

### 1. Projeyi Klonlayın

```bash
git clone <repository-url>
cd backend-path
```

### 2. Environment Dosyasını Oluşturun

```bash
cp .env.example .env
```

`.env` dosyasını düzenleyerek veritabanı bağlantı bilgilerinizi ayarlayın:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=backend_db

# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s

# Logging
LOG_LEVEL=info

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h
```

### 3. Veritabanını Kurun

PostgreSQL veritabanını oluşturun ve migration'ları çalıştırın:

```bash
# Veritabanını oluşturun
createdb backend_db

# Migration'ları çalıştırın
go run cmd/api/main.go
```

### 4. Uygulamayı Çalıştırın

```bash
# Bağımlılıkları yükleyin
go mod tidy

# Uygulamayı başlatın
go run cmd/api/main.go
```

## 🐳 Docker ile Kurulum

### Docker Compose ile Tüm Servisleri Başlatın

```bash
# Tüm servisleri başlatın
docker-compose up --build -d

# Logları görüntüleyin
docker-compose logs -f app
```

### Servisler

- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686

## 📊 Monitoring & Observability

### Health Check Endpoints

- `GET /health` - Uygulama sağlık durumu
- `GET /ready` - Uygulama hazır olma durumu
- `GET /metrics` - Prometheus metrics

### Prometheus Metrics

Uygulama aşağıdaki metrikleri toplar:

- `http_requests_total` - Toplam HTTP istek sayısı
- `http_request_duration_seconds` - İstek süreleri
- `http_requests_in_flight` - Aktif istek sayısı
- `http_request_size_bytes` - İstek boyutları
- `http_response_size_bytes` - Yanıt boyutları

### Grafana Dashboard

Grafana'da aşağıdaki dashboard'lar otomatik olarak yüklenir:

- Backend API Overview
- HTTP Request Rate
- Response Time
- Error Rate
- Active Connections

### Distributed Tracing

Jaeger ile distributed tracing aktif:

- HTTP isteklerinin trace edilmesi
- Span context propagation
- Performance analizi

## 🏗️ Proje Yapısı

```
backend-path/
├── cmd/
│   └── api/
│       └── main.go              # Uygulama giriş noktası
├── config/
│   └── config.go                # Yapılandırma yönetimi
├── internal/
│   ├── app/
│   │   └── app.go               # Ana uygulama yapısı
│   ├── domain/
│   │   ├── user.go              # Kullanıcı domain modeli
│   │   ├── transaction.go       # İşlem domain modeli
│   │   ├── wallet.go            # Cüzdan domain modeli
│   │   └── payment.go           # Ödeme domain modeli
│   ├── handler/
│   │   ├── user_handler.go      # Kullanıcı HTTP handler'ları
│   │   ├── transaction_handler.go
│   │   ├── wallet_handler.go
│   │   └── payment_handler.go
│   ├── middleware/
│   │   ├── auth.go              # Kimlik doğrulama middleware
│   │   ├── logger.go            # Loglama middleware
│   │   ├── monitor.go           # Monitoring middleware
│   │   └── tracing.go           # Tracing middleware
│   ├── models/
│   │   ├── user.go              # GORM modeli
│   │   ├── transaction.go
│   │   ├── wallet.go
│   │   └── payment.go
│   ├── repository/
│   │   ├── database.go          # Veritabanı bağlantısı
│   │   ├── user_repository.go   # Kullanıcı repository
│   │   ├── transaction_repository.go
│   │   ├── wallet_repository.go
│   │   └── payment_repository.go
│   └── service/
│       ├── user_service.go      # Kullanıcı business logic
│       ├── transaction_service.go
│       ├── wallet_service.go
│       └── payment_service.go
├── monitoring/
│   ├── prometheus.yml           # Prometheus konfigürasyonu
│   └── grafana/                 # Grafana dashboard'ları
├── scripts/
│   └── init.sql                 # Veritabanı init script'i
├── Dockerfile                   # Multi-stage Dockerfile
├── docker-compose.yml           # Docker Compose konfigürasyonu
└── README.md
```

## 🔧 API Endpoints

### Kullanıcı İşlemleri

- `POST /api/v1/users/` - Kullanıcı oluştur
- `GET /api/v1/users/:id` - Kullanıcı getir
- `PUT /api/v1/users/:id` - Kullanıcı güncelle
- `DELETE /api/v1/users/:id` - Kullanıcı sil
- `GET /api/v1/users/` - Tüm kullanıcıları listele

### İşlem İşlemleri

- `POST /api/v1/transactions/` - İşlem oluştur
- `GET /api/v1/transactions/:id` - İşlem getir
- `GET /api/v1/transactions/` - Tüm işlemleri listele

### Cüzdan İşlemleri

- `POST /api/v1/wallets/` - Cüzdan oluştur
- `GET /api/v1/wallets/:id` - Cüzdan getir
- `PUT /api/v1/wallets/:id` - Cüzdan güncelle
- `GET /api/v1/wallets/` - Tüm cüzdanları listele

### Ödeme İşlemleri

- `POST /api/v1/payments/` - Ödeme oluştur
- `GET /api/v1/payments/:id` - Ödeme getir
- `GET /api/v1/payments/` - Tüm ödemeleri listele

## 🧪 Test

```bash
# Tüm testleri çalıştır
go test ./...

# Coverage ile test
go test -cover ./...

# Benchmark testleri
go test -bench=. ./...
```

## 📝 Loglama

Uygulama structured logging kullanır. Loglar JSON formatında çıktı verir:

```json
{
  "level": "info",
  "time": "2024-01-15T10:30:00Z",
  "message": "HTTP Request",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": "15.2ms",
  "ip": "127.0.0.1"
}
```

## 🔒 Güvenlik

- JWT tabanlı kimlik doğrulama
- Password hashing (bcrypt)
- CORS yapılandırması
- Rate limiting (opsiyonel)
- Input validation

## 🚀 Production Deployment

### Environment Variables

Production ortamında aşağıdaki environment variable'ları ayarlayın:

```env
ENV=production
LOG_LEVEL=warn
DB_HOST=your_production_db_host
DB_PASSWORD=your_secure_password
JWT_SECRET=your_very_secure_jwt_secret
```

### Docker Production Build

```bash
# Production image oluştur
docker build -t backend-api:latest .

# Container çalıştır
docker run -d \
  --name backend-api \
  -p 8080:8080 \
  --env-file .env \
  backend-api:latest
```

## 🤝 Katkıda Bulunma

1. Fork yapın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Commit yapın (`git commit -m 'Add amazing feature'`)
4. Push yapın (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır.

## 🆘 Destek

Herhangi bir sorun yaşarsanız:

1. GitHub Issues'da sorun bildirin
2. Documentation'ı kontrol edin
3. Community forum'larına başvurun

---

**Not**: Bu proje eğitim amaçlı geliştirilmiştir. Production kullanımı için ek güvenlik önlemleri alınması gerekebilir. 