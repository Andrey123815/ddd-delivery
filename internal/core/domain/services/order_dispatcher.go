package services

import (
	courierModel "delivery/internal/core/domain/models/courier"
	orderModel "delivery/internal/core/domain/models/order"
	"errors"
	"fmt"
)

type OrderDispatcherService interface {
	ReserveCourierForOrder(order *orderModel.Order, couriers []*courierModel.Courier) (*courierModel.Courier, error)
}

var _ OrderDispatcherService = &orderDispatcherService{}

type orderDispatcherService struct{}

func NewOrderDispatcher() OrderDispatcherService {
	return &orderDispatcherService{}
}

func (o *orderDispatcherService) ReserveCourierForOrder(order *orderModel.Order, couriers []*courierModel.Courier) (*courierModel.Courier, error)  {
	if order == nil {
		return nil, errors.New("Заказ для прикрепления к курьеру не передан")
	}

	if len(couriers) == 0 {
		return nil, errors.New("Курьеры для резервирования не переданы")
	}

	if order.Status() != orderModel.OrderStatusCreated {
		return nil, fmt.Errorf("Заказ в статусе %s не может быть назначен курьеру", order.Status())
	}

	var bestCourier *courierModel.Courier
	var bestTime float64

	orderLocation := order.Location()

	for _, courier := range couriers {
		ok, err := courier.CanTakeOrder(order)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		courierTime, err := courier.CalculateTimeToLocation(orderLocation)
		if err != nil {
			return nil, err
		}

		if bestCourier == nil || courierTime < bestTime {
			bestCourier = courier
			bestTime = courierTime
		}
	}

	if bestCourier == nil {
		return nil, errors.New("не найден подходящий курьер для заказа")
	}

	return bestCourier, nil
}
