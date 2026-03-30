package assignCourier

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AssignCourierHandler interface {
	Handle(ctx context.Context) error
}

var _ AssignCourierHandler = &assignCourierHandler{}

type assignCourierHandler struct {
	orderDispatcherService services.OrderDispatcherService
	uowFactory             ports.UnitOfWorkFactory
}

func NewAssignCourierHandler(orderDispatcherService services.OrderDispatcherService, uowFactory ports.UnitOfWorkFactory) (AssignCourierHandler, error) {
	if orderDispatcherService == nil {
		return nil, errs.NewValueIsRequired("orderDispatcherService")
	}
	if uowFactory == nil {
		return nil, errs.NewValueIsRequired("uowFactory")
	}

	return &assignCourierHandler{
		orderDispatcherService: orderDispatcherService,
		uowFactory:             uowFactory,
	}, nil
}

func (h *assignCourierHandler) Handle(ctx context.Context) error {
	uow, err := h.uowFactory.New(ctx)
	if err != nil {
		return err
	}

	uow.Begin(ctx)
	defer uow.RollbackUnlessCommitted(ctx)

	noAssignedOrder, err := uow.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		return err
	}
	if noAssignedOrder == nil {
		return nil
	}

	availableCouriers, err := uow.CourierRepository().GetAllAvailableCouriers(ctx)
	if err != nil {
		return err
	}
	if len(availableCouriers) == 0 {
		return nil
	}

	courier, err := h.orderDispatcherService.ReserveCourierForOrder(noAssignedOrder, availableCouriers)
	if err != nil {
		return err
	}
	if courier == nil {
		return errs.NewValueIsRequired("courier")
	}

	if err := uow.OrderRepository().Update(ctx, noAssignedOrder); err != nil {
		return err
	}

	if err := uow.CourierRepository().Update(ctx, courier); err != nil {
		return err
	}

	return uow.Commit(ctx)
}
