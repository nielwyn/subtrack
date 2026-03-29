package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level, encoding string) error {
	var config zap.Config

	if encoding == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func Get() *zap.Logger {
	if log == nil {
		log, _ = zap.NewProduction() // fallback if Init was never called
	}
	return log
}

func Info(msg string, fields ...zap.Field)  { Get().Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Get().Error(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Get().Warn(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { Get().Debug(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Get().Fatal(msg, fields...) }

func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
