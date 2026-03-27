package createOrder

type CreateOrderCommand struct {
	volume int

	isValid bool
}

const defaultOrderVolume = 1

func NewCreateOrderCommand() (*CreateOrderCommand, error) {
	return &CreateOrderCommand{
		volume:  defaultOrderVolume,
		isValid: true,
	}, nil
}

func (c *CreateOrderCommand) Volume() int {
	return c.volume
}

func (c *CreateOrderCommand) IsValid() bool {
	return c.isValid
}
