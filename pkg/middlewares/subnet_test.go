package middlewares_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/novoseltcev/go-course/pkg/middlewares"
	"github.com/novoseltcev/go-course/pkg/testutils/helpers"
)

var testSubnets = []net.IPNet{
	{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
	{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(12, 32)},
}

func TestTrustedSubnetsSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ip   string
	}{
		{
			name: "first host of 127.0.0.0\\8",
			ip:   "127.0.0.1",
		},
		{
			name: "last host (broadcast) of 127.0.0.0\\8",
			ip:   "127.255.255.255",
		},
		{
			name: "first host of 192.168.0.0\\16",
			ip:   "192.168.0.1",
		},
		{
			name: "last host (broadcast) of 192.168.0.0\\16",
			ip:   "192.168.255.255",
		},
		{
			name: "first host of 10.0.0.0\\12",
			ip:   "10.0.0.1",
		},
		{
			name: "last host (broadcast) of 10.0.0.0\\12",
			ip:   "10.15.255.255",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(middlewares.TrustedSubnets(testSubnets)(helpers.Webhook(t)))
			defer ts.Close()

			resp := helpers.SendRequest(t, ts, nil, map[string]string{
				"X-Real-IP": tt.ip,
			})
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestTrustedSubnetsFail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ip   string
	}{
		{name: "without X-Real-IP header"},
		{
			name: "min border of 127.0.0.0\\8",
			ip:   "126.255.255.255",
		},
		{
			name: "max border of 127.0.0.0\\8",
			ip:   "128.0.0.0",
		},
		{
			name: "min border (multicast) of 192.168.0.0\\16",
			ip:   "192.167.255.255",
		},
		{
			name: "max border of 192.168.0.0\\16",
			ip:   "192.169.0.0",
		},
		{
			name: "min border (multicast) of 10.0.0.0\\12",
			ip:   "9.255.255.255",
		},
		{
			name: "max border of 10.0.0.0\\12",
			ip:   "10.16.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := httptest.NewServer(middlewares.TrustedSubnets(testSubnets)(helpers.Webhook(t)))
			defer ts.Close()

			resp := helpers.SendRequest(t, ts, nil, map[string]string{
				"X-Real-IP": tt.ip,
			})
			defer resp.Body.Close()

			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	}
}
