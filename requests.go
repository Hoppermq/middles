package middles

import (
	"context"
	"net/http"

	"github.com/hoppermq/middles/pkg"
)

func GenerateRequestID(uuidGenerator pkg.UUIDGenerator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, err := uuidGenerator.Generate()
		if err != nil {
			// fall back to error status here

		}
		ctx = context.WithValue(ctx, "request_id", id)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
