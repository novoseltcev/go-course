package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/pkg/middlewares"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func TestDecryptSuccess(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dec := NewMockdecryptor(ctrl)
	ts := httptest.NewServer(middlewares.Decrypt(dec)(helpers.Webhook(t)))
	defer ts.Close()

	got := []byte{1, 2, 3}
	dec.EXPECT().Decrypt(got).Return([]byte(testutils.STRING), nil)

	resp := helpers.SendRequest(t, ts, bytes.NewBuffer(got), nil)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, testutils.STRING, string(body))
}
