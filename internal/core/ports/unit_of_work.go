package ports

import (
	"context"
	"delivery/internal/pkg/ddd"

	"gorm.io/gorm"
)

type UnitOfWork interface {
	Tx() *gorm.DB
	Db() *gorm.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)

	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	RollbackUnlessCommitted(ctx context.Context)

	CourierRepository() CourierRepository
	OrderRepository() OrderRepository
	OutboxRepository() OutboxRepository
}
