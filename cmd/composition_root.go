package cmd

import (
	"delivery/internal/adapters/out/outbox"
	"delivery/internal/core/application/usecases/events"
	"log"
	"reflect"
)

type CompositionRoot struct {
	configs Config

	closers []Closer
	EventRegistry outbox.EventRegistry
}

func NewCompositionRoot(configs Config) *CompositionRoot {
	eventRegistry, err := outbox.NewEventRegistry()
	if err != nil {
		log.Fatalf("cannot create EventRegistry: %v", err)
	}

	for _, t := range []reflect.Type{
		reflect.TypeOf(events.OrderAssignedDomainEvent{}),
		reflect.TypeOf(events.OrderCompletedDomainEvent{}),
	} {
		if err := eventRegistry.RegisterDomainEvent(t); err != nil {
			log.Fatalf("cannot register domain event %s: %v", t.Name(), err)
		}
	}

	return &CompositionRoot{
		configs:       configs,
		EventRegistry: eventRegistry,
	}
}

///////////////////////////////////////////////////////////
//////////////////// LIFECYCLE ////////////////////////////
///////////////////////////////////////////////////////////

func (cr *CompositionRoot) RegisterCloser(c Closer) {
	cr.closers = append(cr.closers, c)
}

func (cr *CompositionRoot) CloseAll() {
	for _, c := range cr.closers {
		if err := c.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
	}
}
