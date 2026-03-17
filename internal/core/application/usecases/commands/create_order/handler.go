package createOrder

import (
	"context"
	"delivery/internal/core/domain/models/kernel"
	orderModel "delivery/internal/core/domain/models/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderHandler interface {
	Handle(ctx context.Context, command *CreateOrderCommand) error
}

var _ CreateOrderHandler = &createOrderHandler{}

type createOrderHandler struct {
	orderRepo ports.OrderRepository
}

func NewCreateOrderHandler(orderRepo ports.OrderRepository) (CreateOrderHandler, error) {
	if orderRepo == nil {
		return nil, errs.NewValueIsRequired("orderRepo")
	}
	return &createOrderHandler{orderRepo: orderRepo}, nil
}

func (h *createOrderHandler) Handle(ctx context.Context, command *CreateOrderCommand) error {
	if command == nil {
		return errs.NewValueIsRequired("command")
	}
	if !command.IsValid() {
		return errs.NewValueIsInvalid("command")
	}

	location, err := kernel.NewLocation(uint8(1) ,uint8(1))
	if err != nil {
		return err
	}

	orderAggregate, err := orderModel.NewOrder(location, command.Volume())
	if err != nil {
		return err
	}	

	return h.orderRepo.Add(ctx, orderAggregate)
}