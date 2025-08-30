# API Contracts (M0)

Public HTTP API exposed by the Gateway. Gateway calls Orders via gRPC internally.

## Create Order

- Method: POST
- Path: `/api/customers/{customerId}/orders`
- Headers:
  - `Content-Type: application/json`
  - `Idempotency-Key: <uuid>` (recommended)
- Request Body:
```
{
  "items": [
    { "id": "sku-123", "quantity": 2 },
    { "id": "sku-456", "quantity": 1 }
  ]
}
```
- Responses:
  - 201 Created
```
{ "orderId": "ord_01HXYZ..." }
```
  - 400 Bad Request (validation error)
  - 409 Conflict (duplicate Idempotency-Key)
  - 500 Internal Server Error

Behavior: Orders persists a PENDING order, publishes `OrderCreated`, and returns the new `orderId`.

## Get Order

- Method: GET
- Path: `/api/customers/{customerId}/orders/{orderId}`
- Responses:
  - 200 OK
```
{
  "id": "ord_01HXYZ...",
  "customerId": "cust_123",
  "status": "PENDING|VALIDATED|PAID|CANCELLED",
  "items": [ { "id": "sku-123", "quantity": 2 } ],
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```
  - 404 Not Found
  - 500 Internal Server Error

## Internal gRPC (Orders)

- Service: `oms.order.v1.OrderService`
  - `CreateOrder(CreateOrderRequest) returns (Order)`
  - `GetOrder(GetOrderRequest) returns (Order)`

