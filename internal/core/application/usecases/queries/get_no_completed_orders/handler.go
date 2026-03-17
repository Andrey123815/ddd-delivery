package getNoCompletedOrders

import (
	"context"
	"delivery/internal/core/domain/models/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type GetNoCompletedOrdersHandler interface {
	Handle(ctx context.Context, query *GetNoCompletedOrdersQuery) ([]GetNoCompletedOrdersResponse, error)
}

var _ GetNoCompletedOrdersHandler = &getNoCompletedOrdersHandler{}

type getNoCompletedOrdersHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewGetNoCompletedOrdersHandler(uowFactory ports.UnitOfWorkFactory) (GetNoCompletedOrdersHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequired("uowFactory")
	}
	return &getNoCompletedOrdersHandler{uowFactory: uowFactory}, nil
}

func (h *getNoCompletedOrdersHandler) Handle(ctx context.Context, query *GetNoCompletedOrdersQuery) ([]GetNoCompletedOrdersResponse, error) {
	if query == nil {
		return nil, errs.NewValueIsRequired("query")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}

	uow.Begin(ctx)

	// Получаем заказы в статусе Created
	createdOrders, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}

	// Получаем заказы в статусе Assigned
	assignedOrders, err := uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		return nil, err
	}

	// Объединяем результаты
	var allOrders []*order.Order
	if createdOrders != nil {
		allOrders = append(allOrders, createdOrders)
	}
	allOrders = append(allOrders, assignedOrders...)

	// Маппим в response
	responses := make([]GetNoCompletedOrdersResponse, 0, len(allOrders))
	for _, o := range allOrders {
		responses = append(responses, ToResponse(o))
	}

	return responses, nil
}
