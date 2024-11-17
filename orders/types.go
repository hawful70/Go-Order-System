package main

import "context"

type OrdersService interface {
	CreateOrder(ctx context.Context) error
}

type OrdersStore interface {
	Create(ctx context.Context) error
}
