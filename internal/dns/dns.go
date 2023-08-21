package dns

import (
	"context"
	"errors"
	"net"
)

var (
	ResolverNotFoundError = errors.New("resolver not found")

	resolvers = map[string]*Resolver{
		"local":           DefaultResolver,
		"google":          NewResolver("8.8.8.8:53"),
		"google-ipv6":     NewResolver("[2001:4860:4860::8888]:53"),
		"cloudflare":      NewResolver("1.1.1.1:53"),
		"cloudflare-ipv6": NewResolver("[2606:4700:4700::1111]:53"),
		"opendns":         NewResolver("208.67.222.222:53"),
		"adguard":         NewResolver("94.140.14.14:53"),
		"authoritative":   nil,
	}
)

type Record struct {
	Source   string    `json:"source"`
	SourceNS string    `json:"source_ns"`
	A        []net.IP  `json:"a"`
	AAAA     []net.IP  `json:"aaaa"`
	MX       []*net.MX `json:"mx"`
	NS       []*net.NS `json:"ns"`
	TXT      []string  `json:"txt"`
	CNAME    string    `json:"cname"`
	Errors   []error   `json:"errors"`
}

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) Resolve(ctx context.Context, host, resolver string) (Record, error) {
	r, ok := resolvers[resolver]
	if !ok {
		return Record{}, ResolverNotFoundError
	}

	if r == nil {
		r = authoritativeResolver(ctx, host)
		if r == nil {
			return Record{}, ResolverNotFoundError
		}
	}

	a, errIP4 := r.resolver.LookupIP(ctx, "ip4", host)
	aaaa, errIP6 := r.resolver.LookupIP(ctx, "ip6", host)
	mx, errMX := r.resolver.LookupMX(ctx, host)
	ns, errNS := r.resolver.LookupNS(ctx, host)
	txt, errTXT := r.resolver.LookupTXT(ctx, host)
	cname, errCNAME := r.resolver.LookupCNAME(ctx, host)

	return Record{
		Source:   resolver,
		SourceNS: r.Address,
		A:        a,
		AAAA:     aaaa,
		MX:       mx,
		NS:       ns,
		TXT:      txt,
		CNAME:    cname,
		Errors: []error{
			errIP4, errIP6, errMX, errNS, errTXT, errCNAME,
		},
	}, nil
}

func authoritativeResolver(ctx context.Context, host string) *Resolver {
	ns, err := DefaultResolver.resolver.LookupNS(ctx, host)
	if err != nil {
		return nil
	}

	if len(ns) == 0 {
		return nil
	}

	return NewResolver(ns[0].Host + ":53")
}
