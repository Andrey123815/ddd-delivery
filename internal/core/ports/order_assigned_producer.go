package ports

import (
	"context"

	"delivery/internal/generated/queues/ordereventspb"
)

type OrderAssignedProducer interface {
	Publish(ctx context.Context, event *ordereventspb.OrderAssignedIntegrationEvent) error
	Close() error
}
