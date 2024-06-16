package grpcserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/app"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	s "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func server(_ context.Context) (pb.EventServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	storage, err := s.NewStorage(c.StorageConf{
		Type: "memory",
	})
	if err != nil {
		log.Printf("error creating storage: %v", err)
	}
	baseServer := grpc.NewServer()
	calendar := app.New(slog.With("service", "calendar"), storage)
	pb.RegisterEventServiceServer(baseServer, NewServer(
		slog.With("service", "grpc_server"), calendar, c.GRPCServerConf{
			Host: "localhost",
			Port: 50051,
		},
	))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewEventServiceClient(conn)

	return client, closer
}

func TestCreateListUpdateMethod(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	out, err := client.Create(ctx, &pb.CreateEventRequest{
		Title:       "Test title",
		Description: "Test description",
		UserId:      "f689adda-d179-4b34-9b26-dfc6eeecd352",
	})
	require.NoError(t, err)
	fmt.Println(out)

	expected := &pb.ListEventResponse{EventList: []*pb.Event{
		{
			Title:       "Test title",
			Description: "Test description",
			UserId:      "f689adda-d179-4b34-9b26-dfc6eeecd352",
		},
	}}
	listResponse, err := client.List(ctx, &pb.ListEventRequest{})
	fmt.Println("listResponse", listResponse)
	require.NoError(t, err)
	require.Equal(
		t,
		listResponse.GetEventList()[0].GetTitle(),
		expected.GetEventList()[0].GetTitle(),
	)
	require.Equal(
		t,
		listResponse.GetEventList()[0].GetDescription(),
		expected.GetEventList()[0].GetDescription(),
	)
	require.Equal(
		t,
		listResponse.GetEventList()[0].GetUserId(),
		expected.GetEventList()[0].GetUserId(),
	)
	_, err = client.Update(ctx, &pb.UpdateEventRequest{
		Id:          listResponse.GetEventList()[0].GetId(),
		Title:       "Updated title",
		Description: "Updated description",
		UserId:      "f689adda-d179-4b34-9b26-dfc6eeecd352",
	})
	require.NoError(t, err)
	listResponse, err = client.List(ctx, &pb.ListEventRequest{})
	require.NoError(t, err)
	require.Equal(
		t,
		listResponse.GetEventList()[0].GetTitle(),
		"Updated title",
	)
	require.Equal(
		t,
		listResponse.GetEventList()[0].GetDescription(),
		"Updated description",
	)
	require.Equal(
		t,
		listResponse.GetEventList()[0].GetUserId(),
		"f689adda-d179-4b34-9b26-dfc6eeecd352",
	)
}

func datesGenerator(dateFrom time.Time, dateTo time.Time) []time.Time {
	var dates []time.Time
	for dateFrom.Before(dateTo) {
		dates = append(dates, dateFrom)
		dateFrom = dateFrom.AddDate(0, 0, 1)
	}
	return dates
}

func TestListMonthEventsMethod(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	dates := datesGenerator(
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
	)

	userID := "f689adda-d179-4b34-9b26-dfc6eeecd352"

	for _, date := range dates {
		_, err := client.Create(ctx, &pb.CreateEventRequest{
			Title:       fmt.Sprintf("Test title %v", date.Format("2006-01-02")),
			Description: fmt.Sprintf("Test description %v", date.Format("2006-01-02")),
			Date:        timestamppb.New(date),
			UserId:      userID,
		})
		require.NoError(t, err)
	}

	listJanResponse, err := client.ListMonthEvents(ctx, &pb.ListMonthEventsRequest{
		Month: &pb.Month{
			Month: int32(1),
			Year:  int32(2025),
		},
	})
	require.NoError(t, err)
	require.Len(t, listJanResponse.GetEventList(), 31)

	listFebResponse, err := client.ListMonthEvents(ctx, &pb.ListMonthEventsRequest{
		Month: &pb.Month{
			Month: int32(2),
			Year:  int32(2025),
		},
	})
	require.NoError(t, err)
	require.Len(t, listFebResponse.GetEventList(), 28)
}

func TestWeekEventsMethod(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	dates := datesGenerator(
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
	)

	userID := "f689adda-d179-4b34-9b26-dfc6eeecd352"

	for _, date := range dates {
		_, err := client.Create(ctx, &pb.CreateEventRequest{
			Title:       fmt.Sprintf("Test title %v", date.Format("2006-01-02")),
			Description: fmt.Sprintf("Test description %v", date.Format("2006-01-02")),
			Date:        timestamppb.New(date),
			UserId:      userID,
		})
		require.NoError(t, err)
	}
	weekEventsResponse, err := client.ListWeekEvents(ctx, &pb.ListWeekEventsRequest{
		Date: &pb.Date{
			Year:  2025,
			Month: 1,
			Day:   6,
		},
	})
	require.NoError(t, err)
	require.Len(t, weekEventsResponse.GetEventList(), 7)
}
