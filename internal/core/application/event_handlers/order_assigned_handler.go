package eventhandlers

import (
	"context"

	"delivery/internal/core/application/usecases/events"
	"delivery/internal/core/ports"
	"delivery/internal/generated/queues/ordereventspb"
	"delivery/internal/pkg/errs"
)

type OrderAssignedHandler struct {
	orderAssignedProducer ports.OrderAssignedProducer
}

func NewOrderAssignedHandler(orderAssignedProducer ports.OrderAssignedProducer) (*OrderAssignedHandler, error) {
	if orderAssignedProducer == nil {
		return nil, errs.NewValueIsRequired("orderAssignedProducer")
	}

	return &OrderAssignedHandler{orderAssignedProducer: orderAssignedProducer}, nil
}

func (h *OrderAssignedHandler) Handle(ctx context.Context, domainEvent *events.OrderAssignedDomainEvent) error {
	if domainEvent == nil {
		return errs.NewValueIsRequired("event")
	}

	msg := &ordereventspb.OrderAssignedIntegrationEvent{
		OrderId: domainEvent.GetOrderId().String(),
	}
	
	return h.orderAssignedProducer.Publish(ctx, msg)
}
