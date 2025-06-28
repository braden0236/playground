package server

import (
	"context"
	"log"
	"net"

	"github.com/braden0236/playground/internal/go-grpc/config"
	"github.com/braden0236/playground/internal/go-grpc/healthz"
	"github.com/braden0236/playground/internal/go-grpc/tls"
	"github.com/braden0236/playground/internal/go-grpc/server/order"
	orderpb "github.com/braden0236/playground/pkg/go-grpc/order"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer   *grpc.Server
	listener     net.Listener
	healthServer *healthz.Server
}

func NewGRPCServer(cfg config.Server) (*Server, error) {
	srvMetrics := grpcprom.NewServerMetrics(grpcprom.WithServerHandlingTimeHistogram())
	prometheus.MustRegister(srvMetrics)

	opts := []grpc.ServerOption{}

	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	if cfg.Metrics.Enabled {
		unaryInterceptors = append(unaryInterceptors, srvMetrics.UnaryServerInterceptor())
		streamInterceptors = append(streamInterceptors, srvMetrics.StreamServerInterceptor())
	}

	opts = append(opts,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	if cfg.UseTLS {
		tlsConfig, err := tls.BuildServerConfig(cfg, cfg.ClientCertAuth)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	grpcSrv := grpc.NewServer(opts...)
	orderpb.RegisterOrderServiceServer(grpcSrv, order.NewService())

	healthSrv := healthz.New()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)
	reflection.Register(grpcSrv)

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, err
	}

	return &Server{
		grpcServer:   grpcSrv,
		listener:     lis,
		healthServer: healthSrv,
	}, nil
}

func (s *Server) Run() error {
	log.Printf("gRPC server listening on %s", s.listener.Addr())
    return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down gRPC server gracefully")
    s.grpcServer.GracefulStop()
    return nil
}

func (s *Server) RunFunc() (func() error, func(error)) {
    return func() error {
            return s.Run()
        }, func(err error) {
            _ = s.Stop(context.Background())
        }
}
