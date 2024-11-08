package endpoints_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/mocks"
)

func TestPing(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	storager := mocks.NewMockMetricStorager(ctrl)
	router := endpoints.NewAPIRouter(storager)

	tests := []struct {
		name string
		err  error
		code int
	}{
		{"success", nil, http.StatusOK},
		{"failed", errors.New("failed"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storager.EXPECT().Ping(gomock.Any()).Return(tt.err)

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.code, w.Code)
		})
	}
}
