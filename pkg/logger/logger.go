package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger, uygulama genelinde kullanılacak logger arayüzü
type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, err error, fields ...map[string]interface{})
	Fatal(msg string, err error, fields ...map[string]interface{})
}

// ZerologLogger, zerolog implementasyonu
type ZerologLogger struct {
	logger zerolog.Logger
}

// New, yeni bir logger oluşturur
func New(level string, pretty bool) Logger {
	var l zerolog.Level

	switch level {
	case "debug":
		l = zerolog.DebugLevel
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	default:
		l = zerolog.InfoLevel
	}

	var output io.Writer = os.Stdout

	// Pretty logging için konsol writer kullan
	if pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Global logger'ı yapılandır
	zerolog.SetGlobalLevel(l)
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	return &ZerologLogger{
		logger: log.Logger,
	}
}

// Debug, debug seviyesinde log mesajı
func (l *ZerologLogger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.logger.Debug()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Info, info seviyesinde log mesajı
func (l *ZerologLogger) Info(msg string, fields ...map[string]interface{}) {
	event := l.logger.Info()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Warn, warn seviyesinde log mesajı
func (l *ZerologLogger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.logger.Warn()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Error, error seviyesinde log mesajı
func (l *ZerologLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Fatal, fatal seviyesinde log mesajı
func (l *ZerologLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
} 