package ports

import (
	"context"
	"delivery/internal/generated/queues/basketeventspb"
)

type BasketConsumer interface {
	Consume(ctx context.Context, event *basketeventspb.BasketConfirmedIntegrationEvent) error
	Close() error
}
