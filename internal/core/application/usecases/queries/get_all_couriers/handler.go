package getAllCouriers

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type GetAllCouriersHandler interface {
	Handle(ctx context.Context, query *GetAllCouriersQuery) ([]GetAllCouriersResponse, error)
}

var _ GetAllCouriersHandler = &getAllCouriersHandler{}

type getAllCouriersHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewGetAllCouriersHandler(uowFactory ports.UnitOfWorkFactory) (GetAllCouriersHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequired("uowFactory")
	}
	return &getAllCouriersHandler{uowFactory: uowFactory}, nil
}

func (h *getAllCouriersHandler) Handle(ctx context.Context, query *GetAllCouriersQuery) ([]GetAllCouriersResponse, error) {
	if query == nil {
		return nil, errs.NewValueIsRequired("query")
	}

	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return nil, err
	}

	uow.Begin(ctx)

	couriers, err := uow.CourierRepository().GetAllAvailableCouriers(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]GetAllCouriersResponse, 0, len(couriers))
	for _, c := range couriers {
		responses = append(responses, ToResponse(c))
	}

	return responses, nil
}
