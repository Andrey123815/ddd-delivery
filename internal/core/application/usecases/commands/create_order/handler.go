package createOrder

import (
	"context"

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
	uowFactory ports.UnitOfWorkFactory
	geoClient  ports.GeoClient
}

func NewCreateOrderHandler(uowFactory ports.UnitOfWorkFactory, geoClient ports.GeoClient) (CreateOrderHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequired("uowFactory")
	}

	if geoClient == nil {
		return nil, errs.NewValueIsRequired("geoClient")
	}

	return &createOrderHandler{uowFactory: uowFactory, geoClient: geoClient}, nil
}

func (h *createOrderHandler) Handle(ctx context.Context, command *CreateOrderCommand) (uuid.UUID, error) {
	if command == nil {
		return uuid.Nil, errs.NewValueIsRequired("command")
	}
	if !command.IsValid() {
		return uuid.Nil, errs.NewValueIsInvalid("command")
	}

	location, err := h.geoClient.GetGeolocation(ctx, command.Address().Street)
	if err != nil {
		return uuid.Nil, err
	}

	orderAggregate, err := orderModel.NewOrder(location, command.Volume())
	if err != nil {
		return uuid.Nil, err
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	uow.Begin(ctx)
	defer uow.RollbackUnlessCommitted(ctx)

	if err := uow.OrderRepository().Add(ctx, orderAggregate); err != nil {
		return uuid.Nil, err
	}

	if err := uow.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return orderAggregate.Id(), nil
}
