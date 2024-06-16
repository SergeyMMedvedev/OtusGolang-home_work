package internalhttp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var mux *http.ServeMux

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s, this is calendar app!\n", r.RemoteAddr)
}

func init() {
	mux = http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/hello", helloHandler)
}

type Server struct {
	configGateway    config.GRPCGateWayConf
	configGRPCServer config.GRPCServerConf
	logger           *slog.Logger
	srv              *http.Server
	app              Application
}

type Application interface { // TODO
}

func NewServer(
	logger *slog.Logger,
	app Application,
	configGateway config.GRPCGateWayConf,
	configGRPCServer config.GRPCServerConf,
) *Server {
	return &Server{
		configGateway:    configGateway,
		configGRPCServer: configGRPCServer,
		logger:           logger,
		app:              app,
	}
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("Start gRPC gateway", "host", s.configGateway.Host, "port", s.configGateway.Port)

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", s.configGateway.Host, s.configGRPCServer.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to create gRPC client: " + err.Error())
	}
	gwmux := runtime.NewServeMux()
	err = pb.RegisterEventServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		slog.Error("failed to register gRPC gateway: " + err.Error())
	}
	s.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.configGateway.Port),
		Handler:           loggingMiddleware(gwmux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stop server...")
	return s.srv.Shutdown(ctx)
}

// gRPC GateWay
// conn, err := grpc.NewClient(
// 	"0.0.0.0:50051",
// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
// )
// if err != nil {
// 	slog.Error("failed to create gRPC client: " + err.Error())
// }
// gwmux := runtime.NewServeMux()
// err = pb.RegisterEventServiceHandler(context.Background(), gwmux, conn)
// if err != nil {
// 	slog.Error("failed to register gRPC gateway: " + err.Error())
// }
// gwServer := &http.Server{
// 	Addr:              ":50052",
// 	Handler:           gwmux,
// 	ReadHeaderTimeout: 30,
// }
