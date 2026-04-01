package cmd

import (
	eventhandlers "delivery/internal/core/application/event_handlers"
	"delivery/internal/core/application/usecases/events"
	"delivery/internal/core/ports"

	mediatr "github.com/mehdihadeli/go-mediatr"
)

func RegisterMediatrNotificationHandlers(
	orderCompletedProducer ports.OrderCompletedProducer,
	orderAssignedProducer ports.OrderAssignedProducer,
) error {
	completed, err := eventhandlers.NewOrderCompletedHandler(orderCompletedProducer)
	if err != nil {
		return err
	}
	if err := mediatr.RegisterNotificationHandler[*events.OrderCompletedDomainEvent](completed); err != nil {
		return err
	}

	assigned, err := eventhandlers.NewOrderAssignedHandler(orderAssignedProducer)
	if err != nil {
		return err
	}
	return mediatr.RegisterNotificationHandler[*events.OrderAssignedDomainEvent](assigned)
}
