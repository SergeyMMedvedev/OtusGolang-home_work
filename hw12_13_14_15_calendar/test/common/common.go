package common

import (
	"context"
	"log/slog"
	"net"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/app"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/server/grpcserver"
	s "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func Client(_ context.Context) (pb.EventServiceClient, error) {
	conn, err := grpc.NewClient(
		"calendar:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewEventServiceClient(conn)
	return client, nil
}

func Server(_ context.Context, storage s.Storage) (pb.EventServiceClient, func()) {
	buffer := 10 * 1024 * 1024
	lis := bufconn.Listen(buffer)
	baseServer := grpc.NewServer()
	calendar := app.New(slog.With("service", "calendar"), storage)
	pb.RegisterEventServiceServer(baseServer, grpcserver.NewServer(
		slog.With("service", "grpc_server"), calendar, c.GRPCServerConf{
			Host: "localhost",
			Port: 50051,
		},
	))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			slog.Error("error serving server: " + err.Error())
		}
	}()
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("error connecting to server: " + err.Error())
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			slog.Error("error closing listener: " + err.Error())
		}
		baseServer.Stop()
	}

	client := pb.NewEventServiceClient(conn)

	return client, closer
}
