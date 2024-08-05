package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	_ "log/slog"
	"os/exec"
	_ "strings"
	_ "sync"
	"testing"
	"time"

	_ "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/consumer"
	_ "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/producer"
	brSchemas "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/broker/schemas"
	c "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	_ "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	common "github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/test/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	title       = "Event TEST title"
	dur         = "30:00:00"
	description = "Event description"
)

func getNewEventWithParams(date time.Time, ntTime int32) schemas.Event {
	event := schemas.NewEvent()

	event.Title = title
	event.Date = date
	event.Duration = dur
	event.Description = description
	event.UserID = uuid.New().String()
	event.NotificationTime = ntTime
	return event
}

func TestCreateEventProduceAndConsume(t *testing.T) {
	var err error
	ctx := context.Background()
	client, err := common.Client(ctx)
	require.NoError(t, err)

	fmt.Println("Clear previous test events.")
	listResponse, err := client.List(ctx, &pb.ListEventRequest{})
	require.NoError(t, err)
	for _, event := range listResponse.GetEventList() {
		if event.Title == title {
			_, err = client.Delete(ctx, &pb.DeleteEventRequest{Id: event.GetId()})
			require.NoError(t, err)
		}
	}

	fmt.Println("Prepare test event.")
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

	fmt.Println("Create test event.")
	_, err = client.Create(ctx, req)
	require.NoError(t, err)

	fmt.Println("Check created event.")
	listResponse, err = client.List(ctx, &pb.ListEventRequest{})
	require.NoError(t, err)
	event.ID = listResponse.GetEventList()[len(listResponse.GetEventList())-1].GetId()

	fmt.Println("Check Day events.")
	listDayResponse, err := client.ListDayEvents(ctx, &pb.ListDayEventsRequest{
		Date: &pb.Date{
			Year:  int32(date.Year()),
			Month: int32(date.Month()),
			Day:   int32(date.Day()),
		},
	})
	require.NoError(t, err)
	dayEvent := listDayResponse.GetEventList()[0]
	require.Equal(t, event.ID, dayEvent.GetId())

	fmt.Println("Check week events.")
	listWeekEvents, err := client.ListWeekEvents(ctx, &pb.ListWeekEventsRequest{
		Date: &pb.Date{
			Year:  int32(date.Year()),
			Month: int32(date.Month()),
			Day:   int32(date.Day()),
		},
	})
	require.NoError(t, err)
	weekEvent := listWeekEvents.GetEventList()[0]
	require.Equal(t, event.ID, weekEvent.GetId())

	fmt.Println("Check month events.")
	listMonthResponse, err := client.ListMonthEvents(ctx, &pb.ListMonthEventsRequest{
		Month: &pb.Month{
			Year:  int32(date.Year()),
			Month: int32(date.Month()),
		},
	})
	require.NoError(t, err)
	monthEvent := listMonthResponse.GetEventList()[0]
	require.Equal(t, event.ID, monthEvent.GetId())

	fmt.Println("Wait for event notification.")
	<-time.After(4 * time.Second)
	require.NoError(t, err)
	expectedNotification := brSchemas.Notification{
		EventID:    event.ID,
		EventTitle: event.Title,
		EventDate:  event.Date.UTC(),
		UserID:     event.UserID,
	}

	fmt.Println("Check consumer received event notification.")
	cmd := exec.Command("docker", "compose", "exec", "sender", "cat", "/tmp/delivered.txt")
	stdout, err := cmd.Output()
	require.NoError(t, err)
	receivedEvent := brSchemas.Notification{}
	err = json.Unmarshal(stdout, &receivedEvent)
	require.NoError(t, err)

	fmt.Println("Compare sended and received event")
	fmt.Println("Sended event: ", expectedNotification)
	fmt.Println("Received event: ", receivedEvent)
	require.Equal(t, expectedNotification.EventID, receivedEvent.EventID)

	fmt.Println("Delete test event.")
	_, err = client.Delete(ctx, &pb.DeleteEventRequest{Id: event.ID})
	require.NoError(t, err)

	fmt.Println("Clear test event notification from file.")
	cmd = exec.Command("docker", "compose", "exec", "sender", "rm", "/tmp/delivered.txt")
	_, err = cmd.Output()
	require.NoError(t, err)
}
