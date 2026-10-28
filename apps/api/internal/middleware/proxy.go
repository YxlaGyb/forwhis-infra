package middleware

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v3"
)

var privateIPNets []*net.IPNet

func init() {
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateCIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil {
			privateIPNets = append(privateIPNets, ipNet)
		}
	}
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, ipNet := range privateIPNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func ProxyFix() fiber.Handler {
	return func(c fiber.Ctx) error {
		remoteAddr := c.RequestCtx().RemoteAddr().String()
		host, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			host = remoteAddr
		}
		clientIP := host

		if isPrivateIP(host) {
			if cfIP := c.Get("CF-Connecting-IP"); cfIP != "" {
				cfIP = strings.TrimSpace(cfIP)
				if cfIP != "" && !isPrivateIP(cfIP) {
					c.Locals("clientIP", cfIP)
					return c.Next()
				}
			}

			if xff := c.Get("X-Forwarded-For"); xff != "" {
				ips := strings.Split(xff, ",")
				if len(ips) > 0 {
					forwardedIP := strings.TrimSpace(ips[0])
					if forwardedIP != "" {
						clientIP = forwardedIP
					}
				}
			}
		}

		c.Locals("clientIP", clientIP)
		return c.Next()
	}
}

func ClientIP(c fiber.Ctx) string {
	if ip, ok := c.Locals("clientIP").(string); ok && ip != "" {
		return ip
	}
	return c.IP()
}
