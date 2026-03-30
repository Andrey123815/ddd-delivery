package ports

import (
	"context"

	"delivery/internal/generated/queues/ordereventspb"
)

type OrderCompletedProducer interface {
	Publish(ctx context.Context, event *ordereventspb.OrderCompletedIntegrationEvent) error
	Close() error
}
