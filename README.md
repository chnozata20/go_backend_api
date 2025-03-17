# Go Backend API

Bu proje, Go ile geliştirilmiş bir backend API örneğidir. Proje, modern Go uygulamaları için temel yapıyı ve en iyi uygulamaları göstermektedir.

## Özellikler

- Modüler proje yapısı
- Environment variables ile yapılandırma
- Zerolog ile yapılandırılabilir loglama
- Graceful shutdown desteği
- HTTP sunucusu
- PostgreSQL veritabanı entegrasyonu
- GORM ORM ile veritabanı işlemleri
- Otomatik migrasyon sistemi

## Proje Yapısı

```
.
├── cmd/
│   └── api/            # Ana uygulama giriş noktası
├── config/             # Yapılandırma kodları
├── docs/               # Dokümantasyon
├── internal/           # Dışa açık olmayan paketler
│   ├── app/            # Uygulama başlatma ve yapılandırma
│   ├── handlers/       # HTTP isteklerini işleyen fonksiyonlar
│   ├── middleware/     # HTTP ara yazılımları
│   ├── models/         # Veri modelleri
│   └── repository/     # Veritabanı işlemleri
└── pkg/                # Dışa açık paketler
    └── logger/         # Loglama paketi
```

## Veritabanı Şeması

Veritabanı şeması hakkında detaylı bilgi için [docs/database_schema.md](docs/database_schema.md) dosyasına bakın.

## Başlangıç

### Gereksinimler

- Go 1.16 veya üzeri
- PostgreSQL 12 veya üzeri

### Kurulum

1. Projeyi klonlayın:
```bash
git clone https://github.com/cihan-ozata/go-backend-api.git
cd go-backend-api
```

2. Bağımlılıkları yükleyin:
```bash
go mod download
```

3. `.env.example` dosyasını `.env` olarak kopyalayın ve gerekli ayarları yapın:
```bash
cp .env.example .env
```

4. PostgreSQL veritabanını oluşturun:
```bash
createdb go_backend_api
```

5. Uygulamayı çalıştırın:
```bash
go run cmd/api/main.go
```

## Yapılandırma

Uygulama, aşağıdaki environment variables ile yapılandırılabilir:

| Değişken | Açıklama | Varsayılan Değer |
|----------|----------|------------------|
| SERVER_PORT | HTTP sunucusu portu | 8080 |
| SERVER_READ_TIMEOUT | HTTP okuma zaman aşımı | 10s |
| SERVER_WRITE_TIMEOUT | HTTP yazma zaman aşımı | 10s |
| SERVER_IDLE_TIMEOUT | HTTP boşta kalma zaman aşımı | 60s |
| DB_HOST | Veritabanı sunucusu | localhost |
| DB_PORT | Veritabanı portu | 5432 |
| DB_USER | Veritabanı kullanıcısı | postgres |
| DB_PASSWORD | Veritabanı şifresi | postgres |
| DB_NAME | Veritabanı adı | go_backend_api |
| DB_SSL_MODE | Veritabanı SSL modu | disable |
| LOG_LEVEL | Loglama seviyesi (debug, info, warn, error) | info |

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın. 