package grpc_server

import (
	"context"
	"fmt"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/storage/schemas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"net"
	"time"
)

type Application interface {
	ListEvents(ctx context.Context) (events []schemas.Event, err error)
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

func (s *Server) List(ctx context.Context, req *pb.ListEventRequest) (*pb.ListEventResponse, error) {
	events, err := s.app.ListEvents(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	eventList := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		var hour, min, sec int
		_, err := fmt.Sscanf(event.Duration, "%d:%d:%d", &hour, &min, &sec)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error parsing duration: %s", err.Error()))
			return nil, err
		}
		duration := time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
		_, err = fmt.Sscanf(event.NotificationTime, "%d:%d:%d", &hour, &min, &sec)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error parsing notification time: %s", err.Error()))
			return nil, err
		}
		notificationTime := time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
		eventList = append(eventList, &pb.Event{
			Id:               event.ID,
			Title:            event.Title,
			Date:             timestamppb.New(event.Date),
			Duration:         durationpb.New(duration),
			Description:      event.Description,
			UserId:           event.UserID,
			NotificationTime: durationpb.New(notificationTime),
		})
	}
	return &pb.ListEventResponse{EventList: eventList}, nil
}

func (s *Server) Create(ctx context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	duration := req.GetDuration().AsDuration()
	// default duration time
	if duration == time.Duration(0) {
		duration = time.Duration(time.Minute * 30)
	}
	notificationTime := req.GetNotificationTime().AsDuration()
	// default notification time
	if notificationTime == time.Duration(0) {
		notificationTime = time.Duration(time.Minute * 15)
	}
	event := schemas.Event{
		Title:            req.GetTitle(),
		Description:      req.GetDescription(),
		Date:             req.GetDate().AsTime(),
		Duration:         duration.String(),
		NotificationTime: notificationTime.String(),
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
		NotificationTime: req.GetNotificationTime().AsDuration().String(),
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
	s.logger.Info("Server started")

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
