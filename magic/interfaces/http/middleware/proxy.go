package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
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

func ProxyFix() gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteAddr := c.Request.RemoteAddr
		host, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			host = remoteAddr
		}

		if isPrivateIP(host) {
			if cfIP := c.GetHeader("CF-Connecting-IP"); cfIP != "" {
				cfIP = strings.TrimSpace(cfIP)
				if cfIP != "" && !isPrivateIP(cfIP) {
					c.Request.RemoteAddr = cfIP + ":0"
					c.Next()
					return
				}
			}

			if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
				ips := strings.Split(xff, ",")
				if len(ips) > 0 {
					clientIP := strings.TrimSpace(ips[0])
					if clientIP != "" {
						c.Request.RemoteAddr = clientIP + ":0"
					}
				}
			}
		}

		c.Next()
	}
}
