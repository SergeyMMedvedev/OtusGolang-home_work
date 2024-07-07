package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/consumer"
	cfg "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := cfg.NewSenderConfig()
	err := config.Read(configFile)
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}
	err = logger.Init(config.Logger)
	if err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}
	slog.Info("config: " + config.String())

	consumer, err := c.NewConsumer(
		slog.With("service", "consumer"),
		config.Broker.URI,
		config.Consumer,
		config.Exchange,
		config.Queue,
		config.Binding,
	)
	if err != nil {
		fmt.Printf("failed to create consumer: %v\n", err)
		os.Exit(1)
	}

	if config.Consumer.Lifetime > 0 {
		slog.Info(
			fmt.Sprintf("running for %s", config.Consumer.Lifetime),
		)
		time.Sleep(config.Consumer.Lifetime)
	} else {
		slog.Info("running forever")
		select {}
	}
	slog.Info("shutting down")

	if err := consumer.Shutdown(); err != nil {
		slog.Error(
			fmt.Sprintf("error during shutdown: %s", err),
		)
	}
}
