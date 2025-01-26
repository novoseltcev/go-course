package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "go.uber.org/mock/gomock"

	"github.com/novoseltcev/go-course/pkg/middlewares"
	"github.com/novoseltcev/go-course/pkg/testutils"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

func TestDecryptSuccess(t *testing.T) {
	t.Parallel()
	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	dec := NewMockdecryptor(ctrl)
	ts := httptest.NewServer(middlewares.Decrypt(dec)(helpers.Webhook(t)))
	defer ts.Close()

	dec.EXPECT().Decrypt(testutils.Bytes).Return([]byte(testutils.STRING), nil)

	resp := helpers.SendRequest(t, ts, bytes.NewBuffer(testutils.Bytes), nil)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, testutils.STRING, string(body))
}

func TestDecryptFailureProcessDecrypt(t *testing.T) {
	t.Parallel()

	ctrl := mock.NewController(t)
	defer ctrl.Finish()

	dec := NewMockdecryptor(ctrl)
	ts := httptest.NewServer(middlewares.Decrypt(dec)(helpers.Webhook(t)))
	defer ts.Close()

	dec.EXPECT().Decrypt(mock.Any()).Return(nil, testutils.Err)

	resp := helpers.SendRequest(t, ts, bytes.NewBuffer(testutils.Bytes), nil)
	defer resp.Body.Close()
	body := helpers.ReadBody(t, resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "failed to decrypt\n", string(body))
}
