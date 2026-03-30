package createOrder

import "delivery/internal/pkg/errs"

type CreateOrderCommand struct {
	volume int
	street string

	isValid bool
}

const defaultOrderVolume = 1

func NewCreateOrderCommand(street string) (*CreateOrderCommand, error) {
	if street == "" {
		return nil, errs.NewValueIsRequired("street")
	}

	return &CreateOrderCommand{
		volume:  defaultOrderVolume,
		street:  street,
		isValid: true,
	}, nil
}

func (c *CreateOrderCommand) Volume() int {
	return c.volume
}

func (c *CreateOrderCommand) IsValid() bool {
	return c.isValid
}

func (c *CreateOrderCommand) Street() string {
	return c.street
}
