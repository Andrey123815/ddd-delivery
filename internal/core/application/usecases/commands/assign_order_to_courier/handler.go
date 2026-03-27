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
	orderRepo ports.OrderRepository
	courierRepo ports.CourierRepository
	orderDispatcherService services.OrderDispatcherService
	uow ports.UnitOfWork
}	

func NewAssignCourierHandler(orderRepo ports.OrderRepository, courierRepo ports.CourierRepository, orderDispatcherService services.OrderDispatcherService, uow ports.UnitOfWork) (AssignCourierHandler, error) {
	if orderRepo == nil {
		return nil, errs.NewValueIsRequired("orderRepo")
	}
	if courierRepo == nil {
		return nil, errs.NewValueIsRequired("courierRepo")
	}
	if orderDispatcherService == nil {
		return nil, errs.NewValueIsRequired("orderDispatcherService")
	}
	
	return &assignCourierHandler{
		orderRepo: orderRepo,
		courierRepo: courierRepo,
		orderDispatcherService: orderDispatcherService,
		uow: uow,
	}, nil
}

func (h *assignCourierHandler) Handle(ctx context.Context) error {
	h.uow.Begin(ctx)
	defer h.uow.RollbackUnlessCommitted(ctx)

	noAssignedOrder, err := h.orderRepo.GetFirstInCreatedStatus(ctx)
	if err != nil {
		return err
	}
	if noAssignedOrder == nil {
		return nil
	}

	availableCouriers, err := h.courierRepo.GetAllAvailableCouriers(ctx)
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

	if err := h.uow.OrderRepository().Update(ctx, noAssignedOrder); err != nil {
		return err
	}

	if err := h.uow.CourierRepository().Update(ctx, courier); err != nil {
		return err
	}

	return h.uow.Commit(ctx)
}
