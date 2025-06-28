package healthz

import (
	"context"
	"sync"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	grpc_health_v1.UnimplementedHealthServer
	mu    sync.RWMutex
	ready bool
}

func New() *Server {
	return &Server{ready: true}
}

func (s *Server) SetNotReady() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = false
}

func (s *Server) Check(ctx context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ready {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
	}
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}, nil
}
