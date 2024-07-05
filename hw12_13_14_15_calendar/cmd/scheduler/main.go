package main

import (
	"context"
	"flag"
	"fmt"
	p "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/producer"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	configFile string
	migrate    bool
)

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
	producer := p.NewProducer(
		slog.With("service", "producer"),
		config.Broker.URI,
	)
	err = producer.Init(
		config.Exchange.Name,
		config.Exchange.Type,
		config.Exchange.Durable,
		config.Exchange.AutoDelete,
		config.Exchange.Internal,
		config.Exchange.NoWait,
		nil,
		config.Exchange.Reliable,
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
	// create ticker with one minute interval for ListEventsForNotification
	ticker := time.NewTicker(time.Second)
	eventsCh := make(chan []schemas.Event)

	// run ListEventsForNotification in goroutine
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// run ListEventsForNotification in infinite loop
		for range ticker.C {
			slog.Info("call ListEventsForNotification")
			events, err := storage.ListEventsForNotification(context.Background())
			if err != nil {
				slog.Error("failed to list events: " + err.Error())
				continue
			}
			eventsCh <- events
		}
	}()
	wg.Add(1)
	// run SendNotification in goroutine
	// SendNotification will be called when eventsCh will be filled with events
	go func() {
		defer wg.Done()
		for events := range eventsCh {
			slog.Info("call SendNotification")
			for _, event := range events {
				slog.Info("send notification for event: " + event.String())
				err := producer.Publish(
					config.Exchange.Name,
					config.Exchange.Key,
					false,
					false,
					event.String(),
				)
				if err != nil {
					slog.Error("failed to publish event " + event.ID + ", error: " + err.Error())
					continue
				}
			}
		}
	}()

	// wait for ticker to stop
	<-ticker.C
	wg.Wait()
}
