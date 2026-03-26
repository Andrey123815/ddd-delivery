package createOrder

import (
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	country string
	city string
	street string
	house string
	apartment string
	volume int

	isValid bool
}

func NewCreateOrderCommand(orderID uuid.UUID, country string, city string, street string, house string, apartment string, 	volume int) (*CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequired("orderID")
	}
	if country == "" {
		return nil, errs.NewValueIsRequired("country")
	}
	if city == "" {
		return nil, errs.NewValueIsRequired("city")
	}
	if street == "" {
		return nil, errs.NewValueIsRequired("street")
	}
	if house == "" {
		return nil, errs.NewValueIsRequired("house")
	}
	if apartment == "" {
		return nil, errs.NewValueIsRequired("apartment")
	}

	return &CreateOrderCommand{
		orderID: orderID,
		country: country,
		city: city,
		street: street,
		house: house,
		apartment: apartment,
		volume: volume,
		isValid: true,
	}, nil
}

func (c *CreateOrderCommand) OrderId() uuid.UUID {
	return c.orderID;
}

func (c *CreateOrderCommand) Country() string {
	return c.country;
}

func (c *CreateOrderCommand) City() string {
	return c.city;
}

func (c *CreateOrderCommand) Street() string {
	return c.street;
}

func (c *CreateOrderCommand) House() string {
	return c.house;
}

func (c *CreateOrderCommand) Apartment() string {
	return c.apartment;
}

func (c *CreateOrderCommand) Volume() int {
	return c.volume;
}

func (c *CreateOrderCommand) IsValid() bool {
	return c.isValid;
}
