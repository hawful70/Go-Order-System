# Event Catalog (M0)

This catalog lists domain events, their purpose, producers, consumers, topics, and keys. All payloads are Protobuf; schemas live under `proto/`.

- Topic naming: `oms.<domain>.v1`
- Partition key: `orderId` (ensures event ordering per order)
- Common envelope: `oms.common.v1.Meta` (eventId, correlationId, causationId, producer, occurredAt, version)

## Order Domain — topic: `oms.order.v1`

- OrderCreated
  - Producer: Orders
  - Consumers: Inventory, (Orders for self‑processing), Analytics
  - Key: `orderId`
  - Payload: `oms.order.events.v1.OrderCreated`
  - Purpose: Signals a new order was created in PENDING state.

- OrderValidated
  - Producer: Orders
  - Consumers: Payments, Analytics
  - Key: `orderId`
  - Payload: `oms.order.events.v1.OrderValidated`
  - Purpose: Inventory is reserved; safe to proceed to payment.

- OrderPaid
  - Producer: Orders
  - Consumers: Fulfillment, Analytics
  - Key: `orderId`
  - Payload: `oms.order.events.v1.OrderPaid`
  - Purpose: Payment authorized; proceed to fulfillment.

- OrderCancelled
  - Producer: Orders
  - Consumers: Inventory (to release), Analytics
  - Key: `orderId`
  - Payload: `oms.order.events.v1.OrderCancelled`
  - Purpose: Terminal state due to inventory/payment failure.

## Inventory Domain — topic: `oms.inventory.v1`

- InventoryReserved
  - Producer: Inventory
  - Consumers: Orders
  - Key: `orderId`
  - Payload: `oms.inventory.events.v1.InventoryReserved`
  - Purpose: Confirms items are reserved for an order.

- InventoryRejected
  - Producer: Inventory
  - Consumers: Orders
  - Key: `orderId`
  - Payload: `oms.inventory.events.v1.InventoryRejected`
  - Purpose: Reservation failed; includes reason (e.g., OUT_OF_STOCK).

## Payment Domain — topic: `oms.payment.v1`

- PaymentAuthorized
  - Producer: Payments
  - Consumers: Orders, Fulfillment
  - Key: `orderId`
  - Payload: `oms.payment.events.v1.PaymentAuthorized`
  - Purpose: Payment is authorized/captured for the order.

- PaymentFailed
  - Producer: Payments
  - Consumers: Orders
  - Key: `orderId`
  - Payload: `oms.payment.events.v1.PaymentFailed`
  - Purpose: Payment failed; Orders should cancel and emit OrderCancelled.

## DLQ and Retries (later milestones)

- Dead‑letter topics: `<topic>.dlq` per domain.
- Retry policy: exponential backoff with max attempts; poison pill handling via DLQ.

