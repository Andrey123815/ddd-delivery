package ports

import (
	"context"

	"delivery/internal/adapters/out/outbox"
)

type OutboxRepository interface {
	Save(ctx context.Context, messages []*outbox.Message) error
	FindAllUnprocessed(ctx context.Context) ([]*outbox.Message, error)
	MarkAsProcessed(ctx context.Context, messages []*outbox.Message) error
}
