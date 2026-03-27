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
	uow ports.UnitOfWork
}

func NewMoveCouriersHandler(uow ports.UnitOfWork) (MoveCouriersHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequired("uow")
	}

	return &moveCouriersHandler{uow: uow}, nil
}

func (h *moveCouriersHandler) Handle(ctx context.Context) error {
	h.uow.Begin(ctx)
	defer h.uow.RollbackUnlessCommitted(ctx)

	assignedOrders, err := h.uow.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		return err
	}

	for _, assignedOrder := range assignedOrders {
		courier, err := h.uow.CourierRepository().Get(ctx, assignedOrder.CourierId())
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

			err = h.uow.OrderRepository().Update(ctx, assignedOrder)
			if err != nil {
				return err
			}
		}

		err = h.uow.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}
	}

	return h.uow.Commit(ctx)
}

