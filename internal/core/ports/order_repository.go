package ports

import (
	"context"
	"delivery/internal/core/domain/models/order"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Add(ctx context.Context, aggregate *order.Order) error
	Update(ctx context.Context, aggregate *order.Order) error
	Get(ctx context.Context, Id uuid.UUID) (*order.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error)
	GetNotCompleted(ctx context.Context) ([]*order.Order, error)
}
