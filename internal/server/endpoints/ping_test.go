package endpoints_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/mocks"
)

func TestPing(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

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
			storager := mocks.NewMockMetricStorager(ctrl)
			storager.EXPECT().Ping(gomock.Any()).Return(tt.err)

			apitest.New().
				Handler(endpoints.NewAPIRouter(storager)).
				Get("/ping").
				Expect(t).
				Status(tt.code)
		})
	}
}
