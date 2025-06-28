package main

import (
	"log"
	"net/http"
	"os"

	"github.com/braden0236/playground/pkg/go-gin/config"
	"github.com/braden0236/playground/pkg/go-gin/metric"
	"github.com/braden0236/playground/pkg/go-gin/server"
)

func main() {

	conf, err := config.Init()
	if err != nil {
		log.Printf("Failed to initialize configuration: %v", err)
		os.Exit(100)
	}

	metrics := metric.NewMetrics(
		metric.WithIgnoredMethods(http.MethodOptions, http.MethodHead),
	)

	server := server.NewServer(metrics)

	go func() {
		if err := server.Run(conf.Server); err != nil{
			log.Printf("Server start failed: %v", err)
			os.Exit(101)
		}
	}()

	if err := server.WaitForShutdown(nil); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

}
