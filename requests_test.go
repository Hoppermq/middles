package middles

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hoppermq/middles/pkg"
	"github.com/stretchr/testify/assert"
)

type mockUUIDGenerator struct {
	value pkg.UUID
	err   error
}

func stringToUUID(s string) pkg.UUID {
	var uuid pkg.UUID
	cleaned := strings.ReplaceAll(s, "-", "")

	if len(cleaned) <= 16 {
		cleaned = cleaned + strings.Repeat("0", 16-len(cleaned))
	} else if len(cleaned) > 16 {
		cleaned = cleaned[:16]
	}

	copy(uuid[:], cleaned)
	return uuid

}

func (m *mockUUIDGenerator) Generate() (pkg.UUID, error) {
	return m.value, m.err
}

func TestGenerateUUID(t *testing.T) {
	t.Parallel()
	type args struct {
		generator mockUUIDGenerator
	}
	tests := []struct {
		name string
		args args
		want pkg.UUID
	}{
		{
			name: "GenerateRequestID",
			args: args{
				generator: mockUUIDGenerator{value: stringToUUID("test-value-uuid-1234")}, // will prob core dump due to exceeded size
			},
			want: stringToUUID("test-value-uuid-1234"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uuid, err := tt.args.generator.Generate()
			assert.Equal(t, nil, err)
			assert.Equal(t, tt.want, uuid)
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	t.Parallel()
	type args struct {
		r         *http.Request
		handler   *http.Handler
		generator mockUUIDGenerator
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "GenerateRequestID",
			args: args{
				r:         httptest.NewRequest(http.MethodGet, "/", nil),
				handler:   nil,
				generator: mockUUIDGenerator{value: stringToUUID("test-value-uuid-1234"), err: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqID := r.Context().Value("request_id")
				assert.NotNil(t, reqID)
				assert.Equal(t, tt.args.generator.value, reqID)
			})

			middleware := GenerateRequestID(&tt.args.generator, next)
			w := httptest.NewRecorder()
			middleware.ServeHTTP(w, tt.args.r)

		})
	}
}
