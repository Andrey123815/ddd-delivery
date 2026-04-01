package utils

import (
	"delivery/cmd"
	"log"
)

func StartKafkaProducers(cr *cmd.CompositionRoot) {
	orderCompletedProducer, err := cr.NewOrderCompletedProducer()
	if err != nil {
		log.Fatalf("cannot create OrderCompletedProducer: %v", err)
	}
	cr.RegisterCloser(orderCompletedProducer)

	orderAssignedProducer, err := cr.NewOrderAssignedProducer()
	if err != nil {
		log.Fatalf("cannot create OrderAssignedProducer: %v", err)
	}
	cr.RegisterCloser(orderAssignedProducer)

	if err := cmd.RegisterMediatrNotificationHandlers(orderCompletedProducer, orderAssignedProducer); err != nil {
		log.Fatalf("cannot register mediatr handlers: %v", err)
	}
}
