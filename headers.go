package middles

import "net/http"

func loadServiceMetadata(w http.ResponseWriter, r *http.Request) {
	if serviceName := r.Context().Value("service_name"); serviceName != nil {
		w.Header().Set("X-Service-Name", serviceName.(string))
	} else {
		w.Header().Set("X-Service-Name", "")
	}

	if serviceVersion := r.Context().Value("service_version"); serviceVersion != nil {
		w.Header().Set("X-Service-Version", serviceVersion.(string))
	} else {
		w.Header().Set("X-Service-Version", "")
	}
}

func SetHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if reqID := r.Context().Value(RequestID); reqID != nil {
			w.Header().Set("X-Request-ID", reqID.(string))
		}

		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "DENY")

		loadServiceMetadata(w, r)
		next.ServeHTTP(w, r)
	})
}
