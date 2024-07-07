package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	p "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/producer"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := c.NewSchedulerConfig()

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
	// run scheduler
	producer, err := p.NewProducer(
		slog.With("service", "producer"),
		config.Broker.URI,
		config.Exchange,
	)
	if err != nil {
		slog.Error("failed to init producer: " + err.Error())
		os.Exit(1)
	}
	storage, err := storage.NewStorage(config.Storage)
	if err != nil {
		slog.Error("failed to create storage: " + err.Error())
		os.Exit(1)
	}
	err = storage.Connect(context.Background())
	if err != nil {
		slog.Error("failed to connect to database: " + err.Error())
		os.Exit(1)
	}

	nTicker := time.NewTicker(time.Second)
	dTicker := time.NewTicker(time.Second * 2)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go producer.SearchEventsForNotification(context.Background(), nTicker, storage, wg)
	wg.Add(1)
	go producer.PublishEventsForNotification(wg, config.Exchange)
	wg.Add(1)
	go producer.SearchAndRemoveOldEvents(context.Background(), dTicker, storage, wg)

	wg.Wait()
}
