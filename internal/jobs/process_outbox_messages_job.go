package jobs

import (
	"context"
	"delivery/internal/adapters/out/outbox"
	"delivery/internal/core/application/usecases/events"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &ProcessOutboxMessagesJob{}

type ProcessOutboxMessagesJob struct {
	outboxRepository ports.OutboxRepository
	eventRegistry    outbox.EventRegistry
}

func NewProcessOutboxMessagesJob(outboxRepository ports.OutboxRepository, eventRegistry outbox.EventRegistry) (cron.Job, error) {
	if outboxRepository == nil {
		return nil, errs.NewValueIsRequired("outboxRepository")
	}
	if eventRegistry == nil {
		return nil, errs.NewValueIsRequired("eventRegistry")
	}

	return &ProcessOutboxMessagesJob{outboxRepository: outboxRepository, eventRegistry: eventRegistry}, nil
}

func (j *ProcessOutboxMessagesJob) Run() {
	ctx := context.Background()

	outboxMessages, err := j.outboxRepository.FindAllUnprocessed(ctx)
	if err != nil {
		log.Errorf("ProcessOutboxMessagesJob: %v", err)
		return
	}

	for _, message := range outboxMessages {
		domainEvent, err := j.eventRegistry.DecodeDomainEvent(message)
		if err != nil {
			log.Errorf("ProcessOutboxMessagesJob: %v", err)
			continue
		}
		if domainEvent == nil {
			continue
		}

		switch domainEvent.(type) {
		case *events.OrderCompletedDomainEvent:
			err = mediatr.Publish[*events.OrderCompletedDomainEvent](ctx, domainEvent.(*events.OrderCompletedDomainEvent))
			if err != nil {
				log.Errorf("ProcessOutboxMessagesJob: %v", err)
				continue
			}
		case *events.OrderAssignedDomainEvent:
			err = mediatr.Publish[*events.OrderAssignedDomainEvent](ctx, domainEvent.(*events.OrderAssignedDomainEvent))
			if err != nil {
				log.Errorf("ProcessOutboxMessagesJob: %v", err)
				continue
			}
		default:
			log.Errorf("ProcessOutboxMessagesJob: unknown domain event type: %T", domainEvent)
			continue
		}

		err = j.outboxRepository.MarkAsProcessed(ctx, []*outbox.Message{message})
		if err != nil {
			log.Errorf("ProcessOutboxMessagesJob: %v", err)
			continue
		}
	}
}
