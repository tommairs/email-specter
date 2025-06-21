package util

import (
	"context"
	"net"
	"time"
)

var dnsServers = []string{
	"8.8.8.8:53",
	"1.1.1.1:53",
	"9.9.9.9:53",
	"8.8.4.4:53",
	"1.0.0.1:53",
	"208.67.222.222:53",
}

func getNextDNSServer(failedIndex int) string {
	return dnsServers[failedIndex%len(dnsServers)]
}

func GetDnsResolver() *net.Resolver {

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			var conn net.Conn
			var err error
			for i := 0; i < len(dnsServers); i++ {
				server := getNextDNSServer(i)
				conn, err = d.DialContext(ctx, network, server)
				if err == nil {
					return conn, nil
				}
			}
			return nil, err
		},
	}

	return resolver

}
