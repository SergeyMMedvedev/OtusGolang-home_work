package internalhttp

import (
	"log/slog"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("Request",
			slog.String("client_ip", r.RemoteAddr),
			slog.String("date", r.Header.Get("Date")),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("protocol", r.Proto),
			slog.String("user_agent", r.UserAgent()),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
