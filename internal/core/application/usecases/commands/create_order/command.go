package createOrder

import (
	"delivery/internal/generated/queues/basketeventspb"
	"delivery/internal/pkg/errs"
)

type CreateOrderCommand struct {
	basketId       string
	address        *basketeventspb.Address
	items          []*basketeventspb.Item
	deliveryPeriod *basketeventspb.DeliveryPeriod
	volume         int
	isValid        bool
}

const defaultOrderVolume = 1

func NewCreateOrderCommand(
	basketId string,
	address *basketeventspb.Address,
	items []*basketeventspb.Item,
	deliveryPeriod *basketeventspb.DeliveryPeriod,
	volume int,
) (*CreateOrderCommand, error) {
	if basketId == "" {
		return nil, errs.NewValueIsRequired("basketId")
	}
	if address == nil {
		return nil, errs.NewValueIsRequired("address")
	}
	if len(items) == 0 {
		return nil, errs.NewValueIsRequired("items")
	}
	if deliveryPeriod == nil || deliveryPeriod.From == 0 || deliveryPeriod.To == 0 {
		return nil, errs.NewValueIsRequired("deliveryPeriod")
	}
	if volume <= 0 {
		volume = defaultOrderVolume
	}

	return &CreateOrderCommand{
		basketId:       basketId,
		address:        address,
		items:          items,
		deliveryPeriod: deliveryPeriod,
		volume:         volume,
		isValid:        true,
	}, nil
}

func (c *CreateOrderCommand) IsValid() bool {
	return c.isValid
}

func (c *CreateOrderCommand) BasketId() string {
	return c.basketId
}

func (c *CreateOrderCommand) Address() *basketeventspb.Address {
	return c.address
}

func (c *CreateOrderCommand) Volume() int {
	return c.volume
}

func (c *CreateOrderCommand) Items() []*basketeventspb.Item {
	return c.items
}

func (c *CreateOrderCommand) DeliveryPeriod() *basketeventspb.DeliveryPeriod {
	return c.deliveryPeriod
}
