package outbox

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, messages []*Message) error
	FindAllUnprocessed(ctx context.Context) ([]*Message, error)
	MarkAsProcessed(ctx context.Context, messages []*Message) error
}

var _ Repository = &RepositoryImpl{}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (Repository, error) {
	return &RepositoryImpl{db: db}, nil
}

func (r *RepositoryImpl) Save(ctx context.Context, messages []*Message) error {
	return r.db.WithContext(ctx).Create(messages).Error
}

func (r *RepositoryImpl) FindAllUnprocessed(ctx context.Context) ([]*Message, error) {
	messages := []*Message{}
	if err := r.db.
		WithContext(ctx).
		Order("occurred_at_utc ASC").
		Limit(100).
		Where("processed_at_utc IS NULL").
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *RepositoryImpl) MarkAsProcessed(ctx context.Context, messages []*Message) error {
	return r.db.WithContext(ctx).Model(messages).Update("processed_at_utc", time.Now().UTC()).Error
}
