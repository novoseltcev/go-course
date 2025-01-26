package middlewares

import (
	"net"
	"net/http"
)

func TrustedSubnets(subnets []net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := net.ParseIP(r.Header.Get("X-Real-IP"))
			if ip == nil || !inSubnets(ip, subnets) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func inSubnets(ip net.IP, subnets []net.IPNet) bool {
	for _, subnet := range subnets {
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}
