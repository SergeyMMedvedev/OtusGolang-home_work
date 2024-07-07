package integration_test

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"

	cnsmr "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/consumer"
	p "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/producer"
	brSchemas "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/schemas"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	s "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	common "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/test/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var loggerConf = c.LoggerConf{
	Level: "INFO",
}

var brokerConf = c.BrokerConf{
	URI: "amqp://guest:guest@localhost:5672/",
}

var exchangeConf = c.ExchangeConf{
	Name:       "calendar",
	Type:       "direct",
	Durable:    true,
	AutoDelete: false,
	Internal:   false,
	NoWait:     false,
	Key:        "scheduler",
	Reliable:   false,
}

var consumerConf = c.ConsumerConf{
	Tag:      "consumer-tag",
	Lifetime: time.Second * 1200,
}

var queueConf = c.QueueConf{
	Name:       "calendar",
	Durable:    true,
	AutoDelete: false,
	Exclusive:  false,
	NoWait:     false,
	Key:        "scheduler",
}

var bindingConf = c.BindingConf{
	QueueName: "calendar",
	Exchange:  "calendar",
	Key:       "scheduler",
	NoWait:    false,
}

var storageConf = c.StorageConf{
	Type: "memory",
}

var schedulerConfig = c.SchedulerConfig{
	Logger:   loggerConf,
	Broker:   brokerConf,
	Exchange: exchangeConf,
	Storage:  storageConf,
}

func getNewEventWithParams(date time.Time, ntTime int32) schemas.Event {
	event := schemas.NewEvent()

	event.Title = "Event title123"
	event.Date = date
	event.Duration = "30:00:00"
	event.Description = "Event description"
	event.UserID = uuid.New().String()
	event.NotificationTime = ntTime
	return event
}

func TestCreateEventProduceAndConsume(t *testing.T) {
	var err error
	ctx := context.Background()
	producer, err := p.NewProducer(
		slog.With("service", "producer"),
		brokerConf.URI,
		exchangeConf,
	)
	require.NoError(t, err)
	storage, err := s.NewStorage(schedulerConfig.Storage)
	require.NoError(t, err)
	err = storage.Connect(context.Background())
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go producer.SearchEventsForNotification(ctx, time.NewTicker(time.Second), storage, wg)
	wg.Add(1)
	go producer.PublishEventsForNotification(wg, exchangeConf)

	consumer, err := cnsmr.NewConsumer(
		slog.With("service", "consumer"),
		brokerConf.URI,
		consumerConf,
		exchangeConf,
		queueConf,
		bindingConf,
	)
	require.NoError(t, err)

	client, closer := common.Server(ctx, storage)
	defer closer()

	date := time.Now().Add(time.Hour*24 + time.Second*3)
	ntTime := int32(1)
	event := getNewEventWithParams(date, ntTime)

	req := &pb.CreateEventRequest{
		Title:            event.Title,
		Description:      event.Description,
		UserId:           event.UserID,
		Date:             timestamppb.New(event.Date),
		NotificationTime: event.NotificationTime,
	}
	_, err = client.Create(ctx, req)
	require.NoError(t, err)
	listResponse, err := client.List(ctx, &pb.ListEventRequest{})

	require.NoError(t, err)
	require.Len(t, listResponse.GetEventList(), 1)
	event.ID = listResponse.GetEventList()[0].GetId()

	<-time.After(4 * time.Second)
	a, err := consumer.Buf.Dequeue()
	require.NoError(t, err)
	byteSlice, ok := a.([]byte)
	require.True(t, ok)
	expectedNotification := brSchemas.Notification{
		EventID:    event.ID,
		EventTitle: event.Title,
		EventDate:  event.Date.UTC(),
		UserID:     event.UserID,
	}
	fmt.Println(string(byteSlice))
	require.Equal(t, expectedNotification.String(), string(byteSlice))
}

func TestRemoveOldEvents(t *testing.T) {
	var err error
	ctx := context.Background()
	producer, err := p.NewProducer(
		slog.With("service", "producer"),
		brokerConf.URI,
		exchangeConf,
	)
	require.NoError(t, err)
	storage, err := s.NewStorage(schedulerConfig.Storage)
	require.NoError(t, err)
	err = storage.Connect(context.Background())
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go producer.SearchAndRemoveOldEvents(ctx, time.NewTicker(time.Second), storage, wg)

	events, err := storage.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 0)
	event := getNewEventWithParams(time.Now().AddDate(-1, 0, -1), 1)
	storage.CreateEvent(ctx, event)
	<-time.After(time.Second * 3)

	events, err = storage.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 0)
}
