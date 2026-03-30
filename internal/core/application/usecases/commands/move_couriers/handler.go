package moveCouriers

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type MoveCouriersHandler interface {
	Handle(ctx context.Context) error
}

var _ MoveCouriersHandler = &moveCouriersHandler{}

type moveCouriersHandler struct {
	uowFactory ports.UnitOfWorkFactory
}

func NewMoveCouriersHandler(uowFactory ports.UnitOfWorkFactory) (MoveCouriersHandler, error) {
	if uowFactory == nil {
		return nil, errs.NewValueIsRequired("uowFactory")
	}

	return &moveCouriersHandler{uowFactory: uowFactory}, nil
}

func (h *moveCouriersHandler) Handle(ctx context.Context) error {
	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}

	uow.Begin(ctx)
	defer uow.RollbackUnlessCommitted(ctx)

	assignedOrders, err := uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		return err
	}

	for _, assignedOrder := range assignedOrders {
		courier, err := uow.CourierRepository().Get(ctx, assignedOrder.CourierId())
		if err != nil {
			return err
		}
		if courier == nil {
			return errs.NewValueIsRequired("courier")
		}

		err = courier.Move(assignedOrder.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equals(assignedOrder.Location()) {
			err = courier.CompleteOrder(assignedOrder.Id())
			if err != nil {
				return err
			}

			err = assignedOrder.Complete(courier.Id())
			if err != nil {
				return err
			}

			err = uow.OrderRepository().Update(ctx, assignedOrder)
			if err != nil {
				return err
			}
		}

		err = uow.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}
	}

	return uow.Commit(ctx)
}

