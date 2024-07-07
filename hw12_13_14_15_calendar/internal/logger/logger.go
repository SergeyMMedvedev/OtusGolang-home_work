package logger

import (
	"fmt"
	"log/slog"
	"os"

	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
)

func Init(conf c.LoggerConf) error {
	// Установка уровня логирования
	var logLevel slog.Level
	switch conf.Level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		return fmt.Errorf("unknown log level: %s", conf.Level)
	}

	logConfig := &slog.HandlerOptions{
		AddSource:   false,
		Level:       logLevel,
		ReplaceAttr: nil,
	}
	logHandler := slog.NewTextHandler(os.Stdout, logConfig)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	slog.Info("log initialized successfully!")
	return nil
}
