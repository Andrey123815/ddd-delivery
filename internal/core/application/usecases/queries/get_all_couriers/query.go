package getAllCouriers

import (
	"context"
	"delivery/internal/core/domain/models/courier"
)

type GetAllCouriersQuery struct {}

func NewGetAllCouriersQuery() (*GetAllCouriersQuery, error) {
	return &GetAllCouriersQuery{}, nil
}

func (q *GetAllCouriersQuery) Execute(ctx context.Context) ([]*courier.Courier, error) {
	return nil, nil
}
