# Architecture Overview (M0)

This document defines the initial architecture and contracts for an event‑driven Order Management System to be deployed on Kubernetes.

## Services (Bounded Contexts)

- Gateway: Public HTTP API. Validates input, calls Orders service for create/get. No direct DB access.
- Orders: System of record for orders. Persists orders, emits domain events, reacts to inventory/payment events to advance state.
- Inventory: Manages stock/availability. Reserves/reverts inventory based on order events.
- Payments (later): Authorizes/captures payments; emits payment outcomes.
- Fulfillment (later): Prepares and delivers orders based on paid events.

Each service owns its database. Cross‑service communication is event‑driven via Kafka; occasional synchronous calls use gRPC.

### Visual: System Context

```mermaid
graph LR
  %% Clients and Gateway
  C[Client Apps] -->|HTTP| GW[Gateway]

  %% Synchronous path
  GW -->|gRPC| ORD[Orders Service]

  %% Datastores
  ORDDB[(Postgres: Orders)]
  INVDB[(Postgres: Inventory)]
  PAYDB[(Postgres: Payments)]
  FULDB[(DB: Fulfillment)]

  %% Kafka bus
  K[((Kafka Cluster))]

  %% Topics (labels on edges)
  ORD -->|publish oms.order.v1| K
  K -->|consume oms.order.v1| INV[Inventory Service]
  K -->|consume oms.order.v1| PAY[Payments Service]

  INV -->|publish oms.inventory.v1| K
  PAY -->|publish oms.payment.v1| K
  K -->|consume oms.inventory.v1| ORD
  K -->|consume oms.payment.v1| ORD

  %% Fulfillment later
  K -->|consume oms.order.v1 (paid)| FUL[Fulfillment Service]

  %% DB ownership
  ORD --- ORDDB
  INV --- INVDB
  PAY --- PAYDB
  FUL --- FULDB
```

## Interaction Model

- External client → Gateway (HTTP):
  - POST /orders: Create order (returns 201 + orderId); gateway calls Orders synchronously to persist PENDING and trigger event emission.
  - GET /orders/{orderId}: Query status (gateway calls Orders synchronously).
- Internal services → Kafka:
  - Orders publishes order.v1 events (Created, Validated, Paid, Cancelled).
  - Inventory publishes inventory.v1 events (Reserved, Rejected).
  - Payments publishes payment.v1 events (Authorized, Failed).
  - Orders consumes inventory/payment events to drive the order state machine.

### Visual: Create Order Saga (Happy and Failure Paths)

```mermaid
sequenceDiagram
  autonumber
  participant Client
  participant GW as Gateway
  participant ORD as Orders
  participant DB as Orders DB
  participant K as Kafka
  participant INV as Inventory
  participant PAY as Payments

  Client->>GW: POST /orders (items, Idempotency-Key)
  GW->>ORD: CreateOrder(customerId, items)
  ORD->>DB: Insert order (PENDING)
  DB-->>ORD: ok (orderId)
  ORD->>K: publish OrderCreated (oms.order.v1)
  ORD-->>GW: 201 {orderId}
  GW-->>Client: 201 {orderId}

  par Inventory path
    K-->>INV: OrderCreated
    INV->>INV: Reserve stock
    alt Reservation OK
      INV->>K: InventoryReserved (oms.inventory.v1)
    else Out of stock
      INV->>K: InventoryRejected (oms.inventory.v1)
    end
  and Orders reacts
    K-->>ORD: InventoryReserved | InventoryRejected
    alt Reserved
      ORD->>DB: Update status VALIDATED
      ORD->>K: OrderValidated
      K-->>PAY: OrderValidated
      PAY->>PAY: Authorize payment
      alt Payment OK
        PAY->>K: PaymentAuthorized (oms.payment.v1)
        K-->>ORD: PaymentAuthorized
        ORD->>DB: Update status PAID
        ORD->>K: OrderPaid
      else Payment Failed
        PAY->>K: PaymentFailed (oms.payment.v1)
        K-->>ORD: PaymentFailed
        ORD->>DB: Update status CANCELLED
        ORD->>K: OrderCancelled
      end
    else Rejected
      ORD->>DB: Update status CANCELLED
      ORD->>K: OrderCancelled
    end
  end
```

## Order State Machine (v1)

- PENDING: Created by Orders when a request arrives. Emits OrderCreated.
- VALIDATED: Orders updates to VALIDATED when InventoryReserved received.
- PAID: Orders updates to PAID when PaymentAuthorized received.
- CANCELLED: If InventoryRejected or PaymentFailed, Orders moves to CANCELLED with reason.

Optional later: PREPARING, COMPLETED (by Fulfillment).

### Visual: Order State Machine

```mermaid
stateDiagram-v2
  [*] --> PENDING
  PENDING --> VALIDATED: InventoryReserved
  PENDING --> CANCELLED: InventoryRejected
  VALIDATED --> PAID: PaymentAuthorized
  VALIDATED --> CANCELLED: PaymentFailed
  PAID --> [*]
  CANCELLED --> [*]
```

## Messaging

- Broker: Kafka (Strimzi on K8s). At‑least‑once delivery, consumer groups per service.
- Topics (keyed by orderId for partition affinity):
  - oms.order.v1
  - oms.inventory.v1
  - oms.payment.v1
- Events are Protobuf messages. Each includes a standard Meta envelope for traceability.
- Retention: dev default; prod tuned. DLQ topics (e.g., oms.order.v1.dlq) added later.

### Visual: Topics and Event Flow

```mermaid
graph LR
  K[((Kafka))]

  subgraph Topics
    O[oms.order.v1]
    I[oms.inventory.v1]
    P[oms.payment.v1]
  end

  O:::topic -->|OrderCreated, OrderValidated, OrderPaid, OrderCancelled| Consumers((Consumers))
  I:::topic -->|InventoryReserved, InventoryRejected| Consumers
  P:::topic -->|PaymentAuthorized, PaymentFailed| Consumers

  classDef topic fill:#eef,stroke:#88a,stroke-width:1px;
```

## Data Stores (initial)

- Orders: Postgres (schema migrations via migrate/goose). Tables: orders, order_items, idempotency (for POST create).
- Inventory: Postgres (items, reservations).
- Payments: Postgres (payments).

## Reliability & Observability (foundations)

- Idempotency: 
  - Create Order uses Idempotency‑Key header; stored in Orders DB.
  - Consumers are idempotent using orderId key and state checks.
- Outbox pattern: planned for M7; begin with simple producer + transactional semantics where possible.
- Tracing: OpenTelemetry across HTTP → gRPC → Kafka; propagate trace context in Meta.
- Health/Readiness: HTTP probes per service; graceful shutdown.

## K8s Deployment Approach (learning‑first)

- Start with `kind` cluster, Kustomize overlays (dev/prod).
- Deploy Gateway first; then Orders + Postgres; then Kafka via Strimzi; then Inventory.

### Visual: Kubernetes (dev) Layout

```mermaid
graph TB
  subgraph kind-cluster
    subgraph ns:oms-dev
      IGW[Ingress] --> SVCGW[Service: gateway]
      SVCGW --> PODGW[(Deployment: gateway pods)]

      SVCORD[Service: orders] --> PODORD[(Deployment: orders pods)]
      STSDB[StatefulSet: postgres] --- PVCDB[(PVC)]

      subgraph Strimzi
        OP[Operator]
        BRK[(Kafka broker)]
      end

      %% Networking
      PODGW -->|gRPC| SVCORD
      PODORD <-->|produce/consume| BRK
      PODORD --> STSDB
    end
  end
```
