package middles

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/hoppermq/middles/pkg"
)

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := pkg.WrappedWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
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
		case ww.StatusCode >= http.StatusOK && ww.StatusCode < http.StatusMultipleChoices:
			logger.Info(
				"request handled",
				"request_id", reqID,
				"status_code", ww.StatusCode,
				"duration", end.Milliseconds(),
			)
		case ww.StatusCode >= http.StatusMultipleChoices && ww.StatusCode < http.StatusBadRequest:
			logger.Info(
				"request redirected",
				"request_id", reqID,
				"status_code", ww.StatusCode,
				"duration", end.Milliseconds(),
			)
		case ww.StatusCode >= http.StatusBadRequest && ww.StatusCode < http.StatusInternalServerError:
			logger.Warn(
				"request failed",
				"request_id", reqID,
				"status_code", ww.StatusCode,
				"duration", end.Milliseconds(),
			)
		case ww.StatusCode >= http.StatusInternalServerError:
			logger.Error(
				"request error",
				"request_id", reqID,
				"status_code", ww.StatusCode,
				"duration", end.Milliseconds(),
			)
		}
	})
}
