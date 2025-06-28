package server

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braden0236/playground/internal/go-grpc/config"
	"github.com/braden0236/playground/internal/go-grpc/healthz"
	"github.com/braden0236/playground/internal/go-grpc/server/order"
	orderpb "github.com/braden0236/playground/pkg/go-grpc/order"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
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

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(srvMetrics.StreamServerInterceptor()),
	}

	if cfg.UseTLS {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			log.Fatalf("failed to load server cert/key: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
			VerifyConnection: func(state tls.ConnectionState) error {
				log.Printf("====")
				log.Printf("🔐 VerifyConnection called: negotiated protocol: %s", state.NegotiatedProtocol)
				if len(state.PeerCertificates) > 0 {
					ce := state.PeerCertificates[0]
					log.Printf("🔐 Client CN: %s, sn: %d", ce.Subject.CommonName, ce.SerialNumber)
				}
				return nil
			},
		}

		if cfg.ClientCertAuth {
			tlsConfig.ClientAuth = tls.RequireAnyClientCert
		}

		// tls.NoClientCert
		// tls.RequestClientCert
		// tls.RequireAnyClientCert
		// tls.RequireAndVerifyClientCert

		if cfg.CaFile != "" {
			caCertBytes, err := os.ReadFile(cfg.CaFile)
			if err != nil {
				log.Fatalf("failed to load client CA cert: %v", err)
			}
			clientCAPool := x509.NewCertPool()
			if !clientCAPool.AppendCertsFromPEM(caCertBytes) {
				log.Fatalf("failed to append client CA cert")
			}
			tlsConfig.ClientCAs = clientCAPool
		} else {
			clientCAPool, err := x509.SystemCertPool()
			if err != nil {
				log.Fatalf("failed to load system root CAs: %v", err)
			}
			tlsConfig.ClientCAs = clientCAPool
		}

		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
		// creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		// if err != nil {
		// 	log.Fatalf("Load TLS fail: %v", err)
		// }
		// opts = append(opts, grpc.Creds(creds))
	}

	grpcSrv := grpc.NewServer(opts...)

	orderpb.RegisterOrderServiceServer(grpcSrv, order.NewService())

	healthSrv := healthz.New()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)

	reflection.Register(grpcSrv)

	// srvMetrics.InitializeMetrics(grpcSrv)

	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, err
	}

	go func() {
		log.Println("📈 Prometheus metrics endpoint on :9090/metrics")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9092", nil); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()

	return &Server{
		grpcServer:   grpcSrv,
		listener:     lis,
		healthServer: healthSrv,
	}, nil
}

func (s *Server) Start() error {
	go s.listenForShutdown()

	log.Printf("gRPC server listening on %s", s.listener.Addr().String())
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) listenForShutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("Received signal: %v. Gracefully shutting down...", sig)

	s.healthServer.SetNotReady()

	s.grpcServer.GracefulStop()
	log.Println("gRPC server stopped gracefully.")
}
