package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/app"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/server/grpcserver"
	internalhttp "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/server/http"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
)

var (
	configFile string
	migrate    bool
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.BoolVar(&migrate, "migrate", false, "Run database migrations")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := c.NewCalendarConfig()
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
	err = storage.Migrate(context.Background())
	if err != nil {
		slog.Error("failed to migrate database: " + err.Error())
		os.Exit(1)
	}

	calendar := app.New(slog.With("service", "calendar"), storage)

	gateway := internalhttp.NewServer(
		slog.With("service", "http_server"),
		calendar,
		config.GRPCGateWay,
		config.GRPCServer,
	)

	grpcServer := internalgrpc.NewServer(
		slog.With("service", "grpc_server"), calendar, config.GRPCServer,
	)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := gateway.Stop(ctx); err != nil {
			slog.Error("failed to stop http server: " + err.Error())
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := gateway.Start(ctx); err != nil {
			slog.Info("gRPC gateway: " + err.Error())
			cancel()
		}
	}()

	go func() {
		defer grpcServer.Stop()
		defer wg.Done()
		if err := grpcServer.Run(ctx); err != nil {
			slog.Info("grpc server: " + err.Error())
			cancel()
		}
		slog.Info("grpc server: 123")
	}()

	wg.Wait()
}
