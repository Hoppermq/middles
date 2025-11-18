package middles_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hoppermq/middles"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	type args struct {
		logger      *slog.Logger
		handlerFunc http.HandlerFunc
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestLoggingSuccessResponse",
			args: args{
				logger: logger,
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			},
			want: http.StatusOK,
		},
		{
			name: "TestLoggingErrorResponse",
			args: args{
				logger: logger,
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				},
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "TestLoggingFailedResponse",
			args: args{
				logger: logger,
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(400)
				},
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testHandler := middles.Logging(tt.args.logger, tt.args.handlerFunc)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			testHandler.ServeHTTP(w, r)
			assert.Equal(t, tt.want, w.Code)
		})
	}
}
