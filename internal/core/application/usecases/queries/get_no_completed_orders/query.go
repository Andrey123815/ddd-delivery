package getNoCompletedOrders

import (
	"context"
	"delivery/internal/core/domain/models/order"
)

type GetNoCompletedOrdersQuery struct {}

func NewGetNoCompletedOrdersQuery() (*GetNoCompletedOrdersQuery, error) {
	return &GetNoCompletedOrdersQuery{}, nil
}

func (q *GetNoCompletedOrdersQuery) Execute(ctx context.Context) ([]*order.Order, error) {
	return nil, nil
}