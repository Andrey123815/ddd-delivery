package services

import (
	"testing"

	"delivery/internal/core/domain/models/kernel"

	courierModel "delivery/internal/core/domain/models/courier"
	orderModel "delivery/internal/core/domain/models/order"
)

func mustNewSpeed(value int) courierModel.Speed {
	speed, err := courierModel.NewSpeed(value)
	if err != nil {
		panic(err)
	}
	return speed
}

func mustNewVolume(value int) courierModel.Volume {
	volume, err := courierModel.NewVolume(value)
	if err != nil {
		panic(err)
	}
	return volume
}

func Test_ReserveCourierForOrderReturnsErrorWhenOrderNil(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := courierModel.NewCourier("Курьер", mustNewSpeed(5), loc)
	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(nil, []*courierModel.Courier{c})
	if err == nil {
		t.Error("expected error when order is nil")
	}
	if got != nil {
		t.Error("expected nil courier when error")
	}
}

func Test_ReserveCourierForOrderReturnsErrorWhenCouriersEmpty(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	ord, _ := orderModel.NewOrder(loc, 10)
	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(ord, nil)
	if err == nil {
		t.Error("expected error when couriers is nil")
	}
	if got != nil {
		t.Error("expected nil courier when error")
	}

	got, err = svc.ReserveCourierForOrder(ord, []*courierModel.Courier{})
	if err == nil {
		t.Error("expected error when couriers slice is empty")
	}
	if got != nil {
		t.Error("expected nil courier when error")
	}
}

func Test_ReserveCourierForOrderReturnsErrorWhenOrderNotCreated(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	ord, _ := orderModel.NewOrder(loc, 10)
	c, _ := courierModel.NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = ord.Assign(c.Id())
	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(ord, []*courierModel.Courier{c})
	if err == nil {
		t.Error("expected error when order status is not Created")
	}
	if got != nil {
		t.Error("expected nil courier when error")
	}
}

func Test_ReserveCourierForOrderReturnsCourierWhenSingleCanTake(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	ord, _ := orderModel.NewOrder(loc, 10)
	c, _ := courierModel.NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(ord, []*courierModel.Courier{c})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != c {
		t.Error("expected the only courier that can take the order")
	}
}

func Test_ReserveCourierForOrderReturnsErrorWhenNoCourierCanTake(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	ord, _ := orderModel.NewOrder(loc, 10)
	c, _ := courierModel.NewCourier("Курьер", mustNewSpeed(5), loc)
	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(ord, []*courierModel.Courier{c})
	if err == nil {
		t.Error("expected error when no courier can take the order")
	}
	if got != nil {
		t.Error("expected nil courier when error")
	}
}

func Test_ReserveCourierForOrderReturnsClosestCourierWhenSeveralCanTake(t *testing.T) {
	locOrder, _ := kernel.NewLocation(1, 1)
	locNear, _ := kernel.NewLocation(1, 2) 
	locFar, _ := kernel.NewLocation(10, 10)
	ord, _ := orderModel.NewOrder(locOrder, 10)

	cNear, _ := courierModel.NewCourier("Ближний", mustNewSpeed(10), locNear)
	_ = cNear.AddStoragePlace("Склад", mustNewVolume(100))
	cFar, _ := courierModel.NewCourier("Дальний", mustNewSpeed(1), locFar)
	_ = cFar.AddStoragePlace("Склад", mustNewVolume(100))

	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(ord, []*courierModel.Courier{cNear, cFar})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != cNear {
		t.Errorf("expected nearer courier (time 0.1), got %v", got)
	}

	ord2, _ := orderModel.NewOrder(locOrder, 10)
	cFar2, _ := courierModel.NewCourier("Дальний2", mustNewSpeed(1), locFar)
	_ = cFar2.AddStoragePlace("Склад", mustNewVolume(100))
	cNear2, _ := courierModel.NewCourier("Ближний2", mustNewSpeed(10), locNear)
	_ = cNear2.AddStoragePlace("Склад", mustNewVolume(100))

	got, err = svc.ReserveCourierForOrder(ord2, []*courierModel.Courier{cFar2, cNear2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != cNear2 {
		t.Error("expected nearer courier when sorted by time, not list order")
	}
}

func Test_ReserveCourierForOrderSkipsCouriersThatCannotTake(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	ord, _ := orderModel.NewOrder(loc, 10)
	cNoStorage, _ := courierModel.NewCourier("Без склада", mustNewSpeed(5), loc)
	cWithStorage, _ := courierModel.NewCourier("Со складом", mustNewSpeed(5), loc)
	_ = cWithStorage.AddStoragePlace("Склад", mustNewVolume(100))
	svc := NewOrderDispatcher()

	got, err := svc.ReserveCourierForOrder(ord, []*courierModel.Courier{cNoStorage, cWithStorage})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != cWithStorage {
		t.Error("expected courier that can take order (has storage)")
	}
}
