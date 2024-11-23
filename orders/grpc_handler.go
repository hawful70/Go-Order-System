package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	pb "github.com/hawful70/common/api"
	"github.com/hawful70/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

// grpcHandler implements the gRPC OrderService interface
type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	channel *amqp.Channel
}

// NewGRPCHandler initializes and registers the gRPC handler
func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

// CreateOrder handles the creation of a new order
func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("Received new order request: %v", req)

	// Create the order object
	order := &pb.Order{
		ID: "43",
	}

	// Publish the order to RabbitMQ
	err := h.publishOrderEvent(ctx, broker.OrderCreatedEvent, order)
	if err != nil {
		log.Printf("Failed to publish order event: %v", err)
		return nil, errors.New("internal server error")
	}

	log.Printf("Order created successfully: %v", order)
	return order, nil
}

// publishOrderEvent publishes an order event to RabbitMQ
func (h *grpcHandler) publishOrderEvent(ctx context.Context, eventName string, order *pb.Order) error {
	// Declare the RabbitMQ queue
	queue, err := h.channel.QueueDeclare(
		eventName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	// Marshal the order to JSON
	orderData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Publish the message
	err = h.channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         orderData,
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		return err
	}

	return nil
}
