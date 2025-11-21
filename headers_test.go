package middles_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hoppermq/middles"
	"github.com/stretchr/testify/assert"
)

func TestHeaderGeneration(t *testing.T) {
	t.Parallel()

	type args struct {
		handlerFunc    http.HandlerFunc
		serviceName    string
		serviceVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "TestGenerateBasicHeader",
			args: args{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
				serviceName:    "test",
				serviceVersion: "v1",
			},
			wantErr: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				r := i2[0].(*http.Request)

				ctx := r.Context()
				ctx = context.WithValue(ctx, "service_name", i2[1].(string))
				ctx = context.WithValue(ctx, "service_version", i2[2].(string))
				r = r.WithContext(ctx)

				testHandler := middles.SetHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				w2 := httptest.NewRecorder()
				testHandler.ServeHTTP(w2, r)

				return assert.Equal(t, http.StatusOK, w2.Code) &&
					assert.Equal(t, "max-age=31536000; includeSubDomains", w2.Header().Get("Strict-Transport-Security")) &&
					assert.Equal(t, "no-cache, no-store, must-revalidate", w2.Header().Get("Cache-Control")) &&
					assert.Equal(t, "nosniff", w2.Header().Get("X-Content-Type-Options")) &&
					assert.Equal(t, "1; mode=block", w2.Header().Get("X-XSS-Protection")) &&
					assert.Equal(t, "DENY", w2.Header().Get("X-Frame-Options")) &&
					assert.Equal(t, "test", w2.Header().Get("X-Service-Name")) &&
					assert.Equal(t, "v1", w2.Header().Get("X-Service-Version"))
			},
		},
		{
			name: "TestGenerateHeaderWithMissingInfo",
			args: args{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
				serviceName:    "",
				serviceVersion: "",
			},
			wantErr: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				r := i2[0].(*http.Request)

				testHandler := middles.SetHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				w2 := httptest.NewRecorder()
				testHandler.ServeHTTP(w2, r)

				return assert.Equal(t, http.StatusOK, w2.Code) &&
					assert.Equal(t, "max-age=31536000; includeSubDomains", w2.Header().Get("Strict-Transport-Security")) &&
					assert.Equal(t, "no-cache, no-store, must-revalidate", w2.Header().Get("Cache-Control")) &&
					assert.Equal(t, "nosniff", w2.Header().Get("X-Content-Type-Options")) &&
					assert.Equal(t, "1; mode=block", w2.Header().Get("X-XSS-Protection")) &&
					assert.Equal(t, "DENY", w2.Header().Get("X-Frame-Options")) &&
					assert.Empty(t, w2.Header().Get("X-Service-Name")) &&
					assert.Empty(t, w2.Header().Get("X-Service-Version")) &&
					assert.Empty(t, w2.Header().Get("X-Request-ID"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testHandler := middles.SetHeader(tt.args.handlerFunc)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			testHandler.ServeHTTP(w, r)
			tt.wantErr(t, w, r, tt.args.serviceName, tt.args.serviceVersion)
		})
	}
}
