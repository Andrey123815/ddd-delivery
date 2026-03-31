package outboxRepo

import (
	"context"
	"time"

	"delivery/internal/adapters/out/outbox"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type Repository struct {
	uow ports.UnitOfWork
}

func NewRepository(uow ports.UnitOfWork) (ports.OutboxRepository, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequired("uow")
	}
	return &Repository{uow: uow}, nil
}

func (r *Repository) conn(ctx context.Context) *gorm.DB {
	if r.uow.InTx() && r.uow.Tx() != nil {
		return r.uow.Tx().WithContext(ctx)
	}
	return r.uow.Db().WithContext(ctx)
}

func (r *Repository) Save(ctx context.Context, messages []*outbox.Message) error {
	for _, m := range messages {
		if err := r.conn(ctx).Create(m).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) FindAllUnprocessed(ctx context.Context) ([]*outbox.Message, error) {
	messages := []*outbox.Message{}
	if err := r.conn(ctx).Where("processed_at_utc IS NULL").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *Repository) MarkAsProcessed(ctx context.Context, messages []*outbox.Message) error {
	now := time.Now().UTC()
	for _, m := range messages {
		if err := r.conn(ctx).Model(m).Update("processed_at_utc", now).Error; err != nil {
			return err
		}
	}
	
	return nil
}
