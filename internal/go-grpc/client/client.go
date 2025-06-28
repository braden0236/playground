package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/braden0236/playground/internal/go-grpc/config"
	"github.com/braden0236/playground/internal/go-grpc/tls"
	orderpb "github.com/braden0236/playground/pkg/go-grpc/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	serviceConfig = `{"loadBalancingPolicy": "round_robin"}`
	concurrency   = 5
)

var (
	mu      sync.Mutex
	counter = make(map[string]int)
	total   int
)

type Client struct {
	conn   *grpc.ClientConn
	client orderpb.OrderServiceClient
}

func New(cfg config.Client) (*Client, error) {

	dialOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                cfg.KeepAlive,
				Timeout:             cfg.KeepAliveTimeout,
				PermitWithoutStream: true,
			},
		),
		// grpc.WithChainUnaryInterceptor(
		// 	middleware.ClientUserHeaderInterceptor,
		// ),
	}

	if cfg.EnableRoundRobin {
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(serviceConfig))
	}

	if cfg.UseTLS {
		tlsConfig, err := tls.BuildConfig(cfg)
		if err != nil {
			return nil, err
		}

		creds := credentials.NewTLS(tlsConfig)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(cfg.Address, dialOpts...)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:   conn,
		client: orderpb.NewOrderServiceClient(conn),
	}

	return c, nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (*orderpb.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.GetOrder(ctx, &orderpb.OrderRequest{OrderId: orderID})
}

func (c *Client) SendBatchRequests() {
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()
			c.SendSingleRequest(i)
		}(i)
	}

	wg.Wait()
}

func (c *Client) SendSingleRequest(workerID int) {
	orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())

	resp, err := c.GetOrder(context.Background(), orderID)
	if err != nil {
		log.Printf("[Worker %d] GetOrder failed: %v", workerID, err)
		time.Sleep(1 * time.Second)
		return
	}

	c.UpdateStats(resp)
}

func (c *Client) UpdateStats(resp *orderpb.OrderResponse) {
	mu.Lock()
	defer mu.Unlock()

	counter[resp.Description]++
	total++
}

func (c *Client) PrintStats() {
	for {
		time.Sleep(5 * time.Second)

		mu.Lock()
		localCounter := make(map[string]int)
		localTotal := total

		for k, v := range counter {
			localCounter[k] = v
		}
		mu.Unlock()

		if localTotal == 0 {
			log.Println("Stats: No data")
			continue
		}

		for desc, cnt := range localCounter {
			percent := float64(cnt) / float64(localTotal) * 100
			log.Printf("Description: %s, Count: %d, Percent: %.2f%%\n", desc, cnt, percent)
		}
	}
}
