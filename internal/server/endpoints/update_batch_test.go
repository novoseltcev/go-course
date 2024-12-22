package endpoints_test

import (
	"net/http"
	"testing"

	"github.com/steinfletcher/apitest"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/internal/schemas"
	"github.com/novoseltcev/go-course/internal/server/endpoints"
	"github.com/novoseltcev/go-course/mocks"
	"github.com/novoseltcev/go-course/pkg/testutils"
)

func TestUpdateBatch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	type got struct {
		body  string
		batch []schemas.Metric
	}

	tests := []struct {
		name string
		got
		err  error
		code int
		body string
	}{
		{
			name: "success",
			got: got{
				body:  `[{"id":"string","type":"gauge","value":10.123},{"id":"string","type":"counter","delta":10}]`,
				batch: []schemas.Metric{testGauge, testCounter},
			},
			code: http.StatusOK,
		},
		{
			name: "skip one",
			got: got{
				body:  `[{"id":"string","type":"gauge","value":10.123},{"id":"string","type":"unknown","value":10.123}]`,
				batch: []schemas.Metric{testGauge},
			},
			code: http.StatusOK,
		},
		{
			name: "skip all",
			got: got{
				body:  `[{"id":"string","type":"unknown","value":10.123}]`,
				batch: nil,
			},
			code: http.StatusOK,
		},
		{
			name: "failed save",
			got: got{
				body:  `[{"id":"string","type":"gauge","value":10.123},{"id":"string","type":"gauge","value":10.123}]`,
				batch: []schemas.Metric{testGauge, testGauge},
			},
			err:  testutils.Err,
			code: http.StatusInternalServerError,
			body: "failed to save metrics\n",
		},
		{
			name: "unmarshalable body",
			got:  got{body: `[{"unknown"}]`},
			code: http.StatusBadRequest,
			body: "parse error: syntax error near offset 11 of '[{\"unknown\"}]'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storager := mocks.NewMockMetricStorager(ctrl)
			if tt.got.batch != nil {
				storager.EXPECT().SaveBatch(gomock.Any(), tt.got.batch).Return(tt.err)
			}

			apitest.New().
				Handler(endpoints.NewAPIRouter(storager)).
				Post("/updates/").
				Body(tt.got.body).
				Expect(t).
				Status(tt.code).
				Body(tt.body).
				End()
		})
	}
}
