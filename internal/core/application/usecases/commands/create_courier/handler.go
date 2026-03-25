package createCourier

import (
	"context"
	"delivery/internal/core/domain/models/courier"
	"delivery/internal/core/domain/models/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierHandler interface {
	Handle(ctx context.Context, command *CreateCourierCommand) error
}

var _ CreateCourierHandler = &createCourierHandler{}

type createCourierHandler struct {
	courierRepo ports.CourierRepository
}

func NewCreateCourierHandler(courierRepo ports.CourierRepository) (CreateCourierHandler, error) {
	if courierRepo == nil {
		return nil, errs.NewValueIsRequired("courierRepo")
	}

	return &createCourierHandler{courierRepo: courierRepo}, nil
}

func (h *createCourierHandler) Handle(ctx context.Context, command *CreateCourierCommand) error {
	if command == nil {
		return errs.NewValueIsRequired("command")
	}
	if !command.IsValid() {
		return errs.NewValueIsInvalid("command")
	}

	location, err := kernel.NewRandomLocation()
	if err != nil {
		return err
	}

	speed, err := courier.NewSpeed(command.Speed())
	if err != nil {
		return err
	}

	courierAggregate, err := courier.NewCourier(command.Name(), speed, *location)
	if err != nil {
		return err
	}
	
	return h.courierRepo.Add(ctx, courierAggregate)
}