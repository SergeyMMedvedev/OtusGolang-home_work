package internalhttp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/SergeyMMedvedev/OtusGolang-home_work/hw12_13_14_15_calendar/internal/config"
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
	config config.ServerConf
	logger *slog.Logger
	srv    *http.Server
	app    Application
}

type Application interface { // TODO
}

func NewServer(logger *slog.Logger, app Application, config config.ServerConf) *Server {
	return &Server{
		config: config,
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("Start server", "host", s.config.Host, "port", s.config.Port)
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	srv.ListenAndServe()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stop server...")
	return s.srv.Shutdown(ctx)
}
