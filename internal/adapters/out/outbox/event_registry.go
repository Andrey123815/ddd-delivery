package outbox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
)

type EventRegistry interface {
	RegisterDomainEvent(eventType reflect.Type) error
	DecodeDomainEvent(event *Message) (ddd.DomainEvent, error)
}

var _ EventRegistry = &eventRegistry{}

type eventRegistry struct {
	eventRegistry map[string]reflect.Type
}

func NewEventRegistry() (EventRegistry, error) {
	return &eventRegistry{eventRegistry: make(map[string]reflect.Type)}, nil
}

func (r *eventRegistry) RegisterDomainEvent(eventType reflect.Type) error {
	if eventType == nil {
		return errs.NewValueIsRequired("eventType")
	}

	et := eventType
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	if et.Kind() != reflect.Struct {
		return fmt.Errorf("eventType must be struct or pointer to struct, got %v", eventType)
	}

	sample, ok := reflect.New(et).Interface().(ddd.DomainEvent)
	if !ok {
		return fmt.Errorf("%s does not implement DomainEvent", et.String())
	}

	name := sample.GetName()
	if name == "" {
		return fmt.Errorf("domain event %s returned empty GetName()", et.Name())
	}

	r.eventRegistry[name] = et
	return nil
}

func encodeDomainEvent(event ddd.DomainEvent) (Message, error) {
	if event == nil {
		return Message{}, errs.NewValueIsRequired("event")
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return Message{}, fmt.Errorf("failed to marshal event: %w", err)
	}

	return Message{
		ID:             event.GetID(),
		Name:           event.GetName(),
		Payload:        payload,
		OccurredAtUtc:  time.Now().UTC(),
		ProcessedAtUtc: nil,
	}, nil
}

func EncodeDomainEvents(events []ddd.DomainEvent) ([]*Message, error) {
	out := make([]*Message, len(events))
	for i := range events {
		msg, err := encodeDomainEvent(events[i])
		if err != nil {
			return nil, err
		}
		m := msg
		out[i] = &m
	}
	return out, nil
}

func (r *eventRegistry) DecodeDomainEvent(outboxMessage *Message) (ddd.DomainEvent, error) {
	if outboxMessage == nil {
		return nil, errs.NewValueIsRequired("outboxMessage")
	}

	t, ok := r.eventRegistry[outboxMessage.Name]
	if !ok {
		return nil, fmt.Errorf("unknown outboxMessage type: %s", outboxMessage.Name)
	}

	eventPtr := reflect.New(t).Interface()
	if err := json.Unmarshal(outboxMessage.Payload, eventPtr); err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	domainEvent, ok := eventPtr.(ddd.DomainEvent)
	if !ok {
		return nil, fmt.Errorf("decoded value does not implement DomainEvent")
	}

	return domainEvent, nil
}
