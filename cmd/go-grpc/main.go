package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braden0236/playground/internal/go-grpc/client"
	"github.com/braden0236/playground/internal/go-grpc/config"
	"github.com/braden0236/playground/internal/go-grpc/server"
	"github.com/oklog/run"
)

var Conf *config.Config

func init() {
	conf, err := config.Init()
	if err != nil {
		log.Printf("Failed to initialize configuration: %v", err)
		os.Exit(100)
	}
	Conf = conf
	log.Printf("Loaded config: %+v", Conf)
}

func main() {
	mode := flag.String("mode", "", "Mode: server | client")
	flag.Parse()

	switch *mode {
	case "server":
		startServer()
	case "client":
		startClient()
	default:
		log.Println("Usage: go run main.go --mode [server|client]")
	}
}

func startServer() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	srv, err := server.NewGRPCServer(Conf.Server)
	if err != nil {
		log.Fatalf("❌ Failed to init gRPC server: %v", err)
	}

	var g run.Group
	g.Add(srv.RunFunc())

	if Conf.Server.Metrics.Enabled {
		metricsSrv := server.NewMetricsServer(Conf.Server)
		metricsSrv.Register("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}))
		g.Add(metricsSrv.RunFunc())
	}

	g.Add(server.WaitForShutdown(ctx, stop))

	if err := g.Run(); err != nil {
		log.Printf("💥 Exited with error: %v", err)
	}
}

func startClient() {

	client, err := client.New(Conf.Client)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	go client.PrintStats()

	for {
		client.SendBatchRequests()
	}

	// resp, err := client.GetOrder(context.Background(), "1234")
	// if err != nil {
	// 	log.Fatalf("Failed to get order: %v", err)
	// }

	// log.Printf("Order received: ID=%s, Status=%s, Amount=%.2f, Description=%s",
	// 	resp.OrderId, resp.Status, resp.Amount, resp.Description)
}
