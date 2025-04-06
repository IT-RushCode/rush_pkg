package utils

import (
	"os"
	"path/filepath"

	"github.com/IT-RushCode/rush_pkg/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitZLogger(cfg *config.Config) *zap.Logger {
	cfgLog := cfg.LOGGER

	// Установка уровня логирования
	var zapLevel zapcore.Level
	if err := zapLevel.Set(cfgLog.Level); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// Базовые настройки энкодера
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Настройка writer'ов
	var cores []zapcore.Core

	// Console core (всегда с цветами)
	consoleEncoderConfig := encoderConfig
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
	cores = append(cores, zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		zapLevel,
	))

	// File cores (если включено)
	if cfgLog.FileLog {
		logDir := "logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("ошибка создания папки logs: " + err.Error())
		}

		// File encoder (без цветов)
		fileEncoderConfig := encoderConfig
		fileEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		var fileEncoder zapcore.Encoder
		if cfgLog.FileJsonFormat {
			fileEncoder = zapcore.NewJSONEncoder(fileEncoderConfig)
		} else {
			fileEncoder = zapcore.NewConsoleEncoder(fileEncoderConfig)
		}

		// Access log
		accessWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logDir, cfgLog.AccessLog.Filename),
			MaxSize:    cfgLog.AccessLog.MaxSize,
			MaxBackups: cfgLog.AccessLog.MaxBackups,
			MaxAge:     cfgLog.AccessLog.MaxAge,
			Compress:   cfgLog.AccessLog.Compress,
		})
		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			accessWriter,
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapLevel && lvl < zapcore.ErrorLevel
			}),
		))

		// Error log
		errorWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logDir, cfgLog.ErrorLog.Filename),
			MaxSize:    cfgLog.ErrorLog.MaxSize,
			MaxBackups: cfgLog.ErrorLog.MaxBackups,
			MaxAge:     cfgLog.ErrorLog.MaxAge,
			Compress:   cfgLog.ErrorLog.Compress,
		})
		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			errorWriter,
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		))
	}

	var options []zap.Option
	options = append(options, zap.AddCaller())

	if cfgLog.EnableStackTrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// Создание логгера
	logger := zap.New(zapcore.NewTee(cores...), options...)

	return logger
}

// WithRID создает поле zap.Field с request_id из контекста
func WithRID(ctx *fiber.Ctx) zap.Field {
	return zap.String("request_id", ctx.GetRespHeader(fiber.HeaderXRequestID))
}
