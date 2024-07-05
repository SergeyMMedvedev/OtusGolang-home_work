package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Application interface {
	ListEvents(ctx context.Context) (events []schemas.Event, err error)
	ListDayEvents(ctx context.Context, date time.Time) (events []schemas.Event, err error)
	ListWeekEvents(ctx context.Context, date time.Time) (events []schemas.Event, err error)
	ListMonthEvents(ctx context.Context, date time.Time) (events []schemas.Event, err error)
	CreateEvent(ctx context.Context, event schemas.Event) error
	DeleteEvent(ctx context.Context, id string) error
	UpdateEvent(ctx context.Context, event schemas.Event) error
}

type Server struct {
	pb.UnimplementedEventServiceServer
	config config.GRPCServerConf
	logger *slog.Logger
	app    Application
	srv    *grpc.Server
}

func prepareEventsList(events []schemas.Event) ([]*pb.Event, error) {
	eventList := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		duration, err := event.GetDurationPb()
		if err != nil {
			slog.Error(fmt.Sprintf("Error parsing duration: %s", err.Error()))
			return nil, err
		}
		eventList = append(eventList, &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			Date:             timestamppb.New(event.Date),
			Duration:         duration,
			Description:      event.Description,
			UserId:           event.UserID,
			NotificationTime: event.NotificationTime,
		})
	}
	return eventList, nil
}

func (s *Server) List(ctx context.Context, _ *pb.ListEventRequest) (*pb.ListEventResponse, error) {
	events, err := s.app.ListEvents(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	eventList, err := prepareEventsList(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error preparing event list: %s", err.Error()))
		return nil, err
	}
	return &pb.ListEventResponse{EventList: eventList}, nil
}

func (s *Server) ListDayEvents(ctx context.Context, req *pb.ListDayEventsRequest) (*pb.ListDayEventsResponse, error) {
	year, month, day := req.GetDate().GetYear(), req.GetDate().GetMonth(), req.GetDate().GetDay()
	err := ValidateYear(year)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid year: %d", year))
		return nil, err
	}
	err = ValidateMonth(month)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid month: %d", month))
		return nil, err
	}
	err = ValidateDay(day, month, year)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid day: %d", day))
		return nil, err
	}
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	events, err := s.app.ListDayEvents(ctx, date)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	eventList, err := prepareEventsList(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	return &pb.ListDayEventsResponse{EventList: eventList}, nil
}

func (s *Server) ListWeekEvents(
	ctx context.Context, req *pb.ListWeekEventsRequest,
) (*pb.ListWeekEventsResponse, error) {
	year, month, day := req.GetDate().GetYear(), req.GetDate().GetMonth(), req.GetDate().GetDay()
	err := ValidateYear(year)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid year: %d", year))
		return nil, err
	}
	err = ValidateMonth(month)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid month: %d", month))
		return nil, err
	}
	err = ValidateDay(day, month, year)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid day: %d", day))
		return nil, err
	}
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	weekStart := StartOfWeek(date)
	weekEnd := weekStart.AddDate(0, 0, 7)
	slog.Info(fmt.Sprintf("weekStart: %s, weekEnd: %s", weekStart, weekEnd))
	events, err := s.app.ListWeekEvents(ctx, weekStart)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	eventList, err := prepareEventsList(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	return &pb.ListWeekEventsResponse{EventList: eventList}, nil
}

func (s *Server) ListMonthEvents(
	ctx context.Context, req *pb.ListMonthEventsRequest,
) (*pb.ListMonthEventsResponse, error) {
	year, month := req.GetMonth().GetYear(), req.GetMonth().GetMonth()
	err := ValidateYear(year)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid year: %d", year))
		return nil, err
	}
	err = ValidateMonth(month)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid month: %d", month))
		return nil, err
	}
	date := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	events, err := s.app.ListMonthEvents(ctx, date)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	eventList, err := prepareEventsList(events)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error listing events: %s", err.Error()))
		return nil, err
	}
	return &pb.ListMonthEventsResponse{EventList: eventList}, nil
}

func (s *Server) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	duration := req.GetDuration().AsDuration()
	// default duration time
	if duration == time.Duration(0) {
		duration = time.Minute * 30
	}
	event := schemas.Event{
		Title:            req.GetTitle(),
		Description:      req.GetDescription(),
		Date:             req.GetDate().AsTime(),
		Duration:         duration.String(),
		NotificationTime: req.GetNotificationTime(),
		UserID:           req.GetUserId(),
	}
	err := s.app.CreateEvent(ctx, event)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pb.CreateEventResponse{}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := s.app.DeleteEvent(ctx, req.GetId())
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pb.DeleteEventResponse{}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	err := s.app.UpdateEvent(ctx, schemas.Event{
		ID:               req.GetId(),
		Title:            req.GetTitle(),
		Description:      req.GetDescription(),
		Date:             req.GetDate().AsTime(),
		Duration:         req.GetDuration().AsDuration().String(),
		NotificationTime: req.GetNotificationTime(),
		UserID:           req.GetUserId(),
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &pb.UpdateEventResponse{}, nil
}

func (s *Server) Run(_ context.Context) error {
	s.logger.Info("Starting server")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer()
	pb.RegisterEventServiceServer(s.srv, s)
	reflection.Register(s.srv)
	s.logger.Info(fmt.Sprintf("Server started on %v:%v", s.config.Host, s.config.Port))

	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.logger.Info("Stopping server")
	s.srv.GracefulStop()
}

func NewServer(logger *slog.Logger, app Application, config config.GRPCServerConf) *Server {
	return &Server{
		UnimplementedEventServiceServer: pb.UnimplementedEventServiceServer{},
		config:                          config,
		logger:                          logger,
		app:                             app,
	}
}
