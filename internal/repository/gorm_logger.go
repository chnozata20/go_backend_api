package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cihan-ozata/backend-path/pkg/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger, GORM için özel logger
type GormLogger struct {
	Logger        logger.Logger
	SlowThreshold time.Duration
}

// NewGormLogger, yeni bir GORM logger oluşturur
func NewGormLogger(log logger.Logger) gormlogger.Interface {
	return &GormLogger{
		Logger:        log,
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LogMode, log modunu ayarlar (GORM logger arayüzü için gerekli)
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info, bilgi mesajlarını loglar
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, data...), nil)
}

// Warn, uyarı mesajlarını loglar
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, data...), nil)
}

// Error, hata mesajlarını loglar
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, data...), errors.New("gorm error"), nil)
}

// Trace, SQL sorgularını loglar
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// Hata varsa logla
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Logger.Error("SQL Error", err, map[string]interface{}{
			"sql":      sql,
			"rows":     rows,
			"elapsed":  elapsed,
			"duration": fmt.Sprintf("%s", elapsed),
		})
		return
	}

	// Yavaş sorguları logla
	if elapsed > l.SlowThreshold {
		l.Logger.Warn("Slow SQL", map[string]interface{}{
			"sql":      sql,
			"rows":     rows,
			"elapsed":  elapsed,
			"duration": fmt.Sprintf("%s", elapsed),
		})
		return
	}

	// Debug modunda tüm sorguları logla
	l.Logger.Debug("SQL", map[string]interface{}{
		"sql":      sql,
		"rows":     rows,
		"elapsed":  elapsed,
		"duration": fmt.Sprintf("%s", elapsed),
	})
} 