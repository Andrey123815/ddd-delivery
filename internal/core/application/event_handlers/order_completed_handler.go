package eventhandlers

import (
	"context"

	"delivery/internal/core/application/usecases/events"
	"delivery/internal/core/ports"
	"delivery/internal/generated/queues/ordereventspb"
	"delivery/internal/pkg/errs"
)

type OrderCompletedHandler struct {
	orderCompletedProducer ports.OrderCompletedProducer
}

func NewOrderCompletedHandler(orderCompletedProducer ports.OrderCompletedProducer) (*OrderCompletedHandler, error) {
	if orderCompletedProducer == nil {
		return nil, errs.NewValueIsRequired("orderCompletedProducer")
	}

	return &OrderCompletedHandler{orderCompletedProducer: orderCompletedProducer}, nil
}

func (h *OrderCompletedHandler) Handle(ctx context.Context, domainEvent *events.OrderCompletedDomainEvent) error {
	if domainEvent == nil {
		return errs.NewValueIsRequired("event")
	}

	msg := &ordereventspb.OrderCompletedIntegrationEvent{
		OrderId: domainEvent.GetOrderId().String(),
	}
	
	return h.orderCompletedProducer.Publish(ctx, msg)
}
