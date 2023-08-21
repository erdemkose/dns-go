package dns

import (
	"context"
	"net"
	"time"
)

var (
	DefaultResolver = &Resolver{
		Address:  "",
		resolver: net.DefaultResolver,
	}
)

type Resolver struct {
	Address  string
	resolver *net.Resolver
}

func NewResolver(addr string) *Resolver {
	return &Resolver{
		Address: addr,
		resolver: &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(1000),
				}
				return d.DialContext(ctx, network, addr)
			},
		},
	}
}
