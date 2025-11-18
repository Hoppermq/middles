package middles

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/hoppermq/middles/pkg"
)

func GenerateRequestID(uuidGenerator func() (uuid.UUID, error), next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "request_id", pkg.Generator(uuidGenerator))

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
