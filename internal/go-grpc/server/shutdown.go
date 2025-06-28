package server

import (
	"context"
	"log"
)

func WaitForShutdown(ctx context.Context, stop context.CancelFunc) (func() error, func(error)) {
    return func() error {
            <-ctx.Done()
            log.Println("WaitForShutdown: received shutdown signal, exiting gracefully")
            return nil
        }, func(error) {
            stop()
        }
}
