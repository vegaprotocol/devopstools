package tools

import (
	"net"
	"time"
)

const MaximumDialDuration = 2 * time.Second

func FilterHealthyGRPCEndpoints(endpoints []string) []string {
	var healthy []string
	for _, endpoint := range endpoints {
		conn, err := net.DialTimeout("tcp", endpoint, MaximumDialDuration)
		if err == nil && conn != nil {
			_ = conn.Close()
			healthy = append(healthy, endpoint)
		}
	}
	return healthy
}
