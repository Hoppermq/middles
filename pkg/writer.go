package pkg

import "net/http"

type WrappedWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *WrappedWriter) WriteHeader(statusCode int) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}
