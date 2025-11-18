package middles

import (
	"context"
	"errors"
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

		if errors.Is(r.Context().Err(), context.DeadlineExceeded) {
			logger.Warn("request timed out",
				"request_id", reqID,
				"duration", end.Milliseconds(),
			)
		}

		switch {
		case ww.statusCode >= http.StatusOK && ww.statusCode < http.StatusMultipleChoices:
			logger.Info(
				"request handled",
				"request_id", reqID,
				"status_code", ww.statusCode,
				"duration", end.Milliseconds(),
			)
		case ww.statusCode >= http.StatusMultipleChoices && ww.statusCode < http.StatusBadRequest:
			logger.Info(
				"request redirected",
				"request_id", reqID,
				"status_code", ww.statusCode,
				"duration", end.Milliseconds(),
			)
		case ww.statusCode >= http.StatusBadRequest && ww.statusCode < http.StatusInternalServerError:
			logger.Warn(
				"request failed",
				"request_id", reqID,
				"status_code", ww.statusCode,
				"duration", end.Milliseconds(),
			)
		case ww.statusCode >= http.StatusInternalServerError:
			logger.Error(
				"request error",
				"request_id", reqID,
				"status_code", ww.statusCode,
				"duration", end.Milliseconds(),
			)
		}
	})
}
