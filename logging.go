package middles

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := wrappedWriter{
			ResponseWriter: w,
		}

		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			path   = r.URL.Path
			proto  = r.Proto
		)

		reqID := r.Context().Value(RequestID)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "path", path, "proto", proto)

		logger.Info(
			"request received",
			"request_id",
			reqID,
			userAttrs,
			requestAttrs,
		)

		start := time.Now()
		next.ServeHTTP(&ww, r)
		end := time.Since(start)

		if ww.statusCode >= http.StatusOK && ww.statusCode < 301 {
			logger.Info(
				"request handled",
				"request_id",
				reqID,
				"status_code",
				ww.statusCode,
				slog.Int64("duration", end.Milliseconds()),
			)
		}
	})
}
