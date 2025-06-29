package dns

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc/resolver"
)

const (
	schemeName = "dns5"
)

type dnsResolverBuilder struct{}

func (*dnsResolverBuilder) Scheme() string {
	return schemeName
}

func (*dnsResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &dnsResolver{
		target: target,
		cc:     cc,
		ticker: time.NewTicker(5 * time.Second),
		ctx:    context.Background(),
	}
	ctx, cancel := context.WithCancel(r.ctx)
	r.cancel = cancel
	go r.start(ctx)
	return r, nil
}

type dnsResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
}

func (r *dnsResolver) start(ctx context.Context) {
	for {
		select {
		case <-r.ticker.C:
			r.resolve()
		case <-ctx.Done():
			return
		}
	}
}

func (r *dnsResolver) resolve() {
	host := r.target.Endpoint
	h := host()
	port := ""

	if strings.Contains(h, ":") {
		parts := strings.Split(h, ":")
		h = parts[0]
		port = parts[1]
	}

	addrs, err := net.LookupHost(h)
	if err != nil {
		log.Printf("[resolver] DNS lookup failed: %v", err)
		return
	}

	var resolvedAddrs []resolver.Address
	for _, addr := range addrs {
		fullAddr := addr
		if port != "" {
			fullAddr = net.JoinHostPort(addr, port)
		}
		resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: fullAddr})
	}

	if len(resolvedAddrs) == 0 {
		log.Printf("[resolver] no available IP")
		return
	}

	err = r.cc.UpdateState(resolver.State{Addresses: resolvedAddrs})
	if err != nil {
		log.Printf("[resolver] UpdateState failed: %v", err)
	}
}

func (r *dnsResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.resolve()
}

func (r *dnsResolver) Close() {
	r.cancel()
	r.ticker.Stop()
}

func init() {
	resolver.Register(&dnsResolverBuilder{})
}
