package server

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/braden0236/playground/internal/go-grpc/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	mux    *http.ServeMux
	server *http.Server
	cfg    config.Server
	mu     sync.Mutex
}

func NewMetricsServer(cfg config.Server) *MetricsServer {
	mux := http.NewServeMux()
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	s := &MetricsServer{
		mux: mux,
		cfg: cfg,
		server: &http.Server{
			Addr:    cfg.Metrics.Address,
			Handler: wrapMetricsHandler(cfg, mux),
		},
	}
	return s
}

func (s *MetricsServer) Register(pattern string, handler http.Handler) {
	log.Printf("Registering metrics handler for %s", pattern)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mux.Handle(pattern, handler)
}

func (s *MetricsServer) Run() error {
	log.Printf("Metrics server listening at %s", s.cfg.Metrics.Address)
	return s.server.ListenAndServe()
}

func (s *MetricsServer) Stop(ctx context.Context) error {
	log.Println("Shutting down metrics server gracefully")
	return s.server.Shutdown(ctx)
}

func (s *MetricsServer) RunFunc() (func() error, func(error)) {
	return func() error {
			return s.Run()
		}, func(err error) {
			if err != nil {
				log.Printf("MetricsServer interrupted due to error: %v", err)
			}
			_ = s.Stop(context.Background())
		}
}

func wrapMetricsHandler(cfg config.Server, handler http.Handler) http.Handler {
	if cfg.Metrics.Auth.Username == "" || cfg.Metrics.Auth.Password == "" {
		return handler
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != cfg.Metrics.Auth.Username || pass != cfg.Metrics.Auth.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
