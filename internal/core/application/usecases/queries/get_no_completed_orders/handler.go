package getNoCompletedOrders

import (
	"context"
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

	noCompletedOrders, err := uow.OrderRepository().GetNotCompleted(ctx)
	if err != nil {
		return nil, err
	}	

	// Маппим в response
	orders := make([]GetNoCompletedOrdersResponse, 0, len(noCompletedOrders))
	for _, o := range noCompletedOrders {
		orders = append(orders, ToResponse(o))
	}

	return orders, nil
}
