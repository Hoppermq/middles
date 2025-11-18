package middles

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			url    = r.URL.String()
			path   = r.URL.Path
			proto  = r.Proto
			start  = time.Now()
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "path", path, "proto", proto)

		logger.Info(
			"request received",
			userAttrs,
			requestAttrs,
			slog.Float64("duration", float64(time.Since(start).Milliseconds())),
		)

		next.ServeHTTP(w, r)

		logger.Info("request handled")

	})
}
