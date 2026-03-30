package createOrder

import (
	"context"

	"delivery/internal/core/domain/models/kernel"
	orderModel "delivery/internal/core/domain/models/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrderHandler interface {
	Handle(ctx context.Context, command *CreateOrderCommand) (uuid.UUID, error)
}

var _ CreateOrderHandler = &createOrderHandler{}

type createOrderHandler struct {
	uow ports.UnitOfWork
}

func NewCreateOrderHandler(uow ports.UnitOfWork) (CreateOrderHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequired("uow")
	}
	return &createOrderHandler{uow: uow}, nil
}

func (h *createOrderHandler) Handle(ctx context.Context, command *CreateOrderCommand) (uuid.UUID, error) {
	if command == nil {
		return uuid.Nil, errs.NewValueIsRequired("command")
	}
	if !command.IsValid() {
		return uuid.Nil, errs.NewValueIsInvalid("command")
	}

	locPtr, err := kernel.NewRandomLocation()
	if err != nil {
		return uuid.Nil, err
	}

	orderAggregate, err := orderModel.NewOrder(*locPtr, command.Volume())
	if err != nil {
		return uuid.Nil, err
	}

	h.uow.Begin(ctx)
	defer h.uow.RollbackUnlessCommitted(ctx)

	if err := h.uow.OrderRepository().Add(ctx, orderAggregate); err != nil {
		return uuid.Nil, err
	}

	if err := h.uow.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return orderAggregate.Id(), nil
}
