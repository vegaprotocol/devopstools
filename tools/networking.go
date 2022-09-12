package tools

import (
	"fmt"
	"net"
)

func GetIP(host string) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", fmt.Errorf("failed to get IP for %s, %w", host, err)
	}
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return string(ipv4), nil
		}
	}
	return "", fmt.Errorf("no IP assigned to %s", host)
}
