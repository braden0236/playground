package order

import (
	"context"
	"log"
	"os"

	orderpb "github.com/braden0236/playground/pkg/go-grpc/order"
)

type Service struct {
	orderpb.UnimplementedOrderServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	// log.Printf("GetOrder: %s", req.OrderId)
	hostname, _ := os.Hostname()

	return &orderpb.OrderResponse{
		OrderId: req.OrderId,
		Status:  "SHIPPED",
		Amount:  100.50,
		Description: hostname,
	}, nil
}

func (s *Service) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Printf("CreateOrder: %s, %.2f", req.OrderId, req.Amount)

	return &orderpb.CreateOrderResponse{
		Success: true,
		Message: "Order created successfully",
	}, nil
}

func (s *Service) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	log.Printf("UpdateOrder: %s -> status=%s, amount=%.2f", req.OrderId, req.Status, req.Amount)

	return &orderpb.UpdateOrderResponse{
		Success: true,
		Message: "Order updated successfully",
	}, nil
}

func (s *Service) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	log.Printf("DeleteOrder: %s", req.OrderId)

	return &orderpb.DeleteOrderResponse{
		Success: true,
		Message: "Order deleted successfully",
	}, nil
}

func (s *Service) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	log.Printf("ListOrders: page=%d, size=%d", req.Page, req.PageSize)

	orders := []*orderpb.OrderResponse{
		{OrderId: "order-001", Status: "CREATED", Amount: 123.45},
		{OrderId: "order-002", Status: "DELIVERED", Amount: 456.78},
	}

	return &orderpb.ListOrdersResponse{
		Orders: orders,
		Total:  int32(len(orders)),
	}, nil
}
