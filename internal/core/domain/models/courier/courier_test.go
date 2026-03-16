package courier

import (
	"testing"

	"delivery/internal/core/domain/models/kernel"

	orderModel "delivery/internal/core/domain/models/order"

	"github.com/google/uuid"
)

func Test_NewCourierCreatedAtValidParams(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, err := NewCourier("Курьер", 5, loc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.Name() != "Курьер" {
		t.Errorf("Name() = %q, want %q", c.Name(), "Курьер")
	}
	if c.Speed() != 5 {
		t.Errorf("Speed() = %d, want 5", c.Speed())
	}
	if c.Location() != loc {
		t.Error("Location() should equal passed location")
	}
	if c.StoragePlaces() != nil {
		t.Error("new courier should have nil storage places")
	}
}

func Test_NewCourierReturnsErrorAtEmptyName(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	_, err := NewCourier("", 5, loc)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func Test_NewCourierReturnsErrorAtInvalidSpeed(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	_, err := NewCourier("Курьер", 0, loc)
	if err == nil {
		t.Error("expected error for speed 0")
	}
	_, err = NewCourier("Курьер", -1, loc)
	if err == nil {
		t.Error("expected error for negative speed")
	}
}

func Test_NewCourierReturnsErrorAtEmptyLocation(t *testing.T) {
	_, err := NewCourier("Курьер", 5, kernel.Location{})
	if err == nil {
		t.Error("expected error for empty location")
	}
}

func Test_CourierEqualsReturnsCorrectResult(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	a, _ := NewCourier("A", 5, loc)
	if !a.Equals(*a) {
		t.Error("courier should equal itself")
	}
	b, _ := NewCourier("B", 5, loc)
	if a.Equals(*b) {
		t.Error("different couriers should not be equal")
	}
}

func Test_CourierAddStoragePlaceSucceeds(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	err := c.AddStoragePlace("Склад", 100)
	if err != nil {
		t.Fatalf("AddStoragePlace failed: %v", err)
	}
	places := c.StoragePlaces()
	if len(places) != 1 {
		t.Fatalf("len(StoragePlaces()) = %d, want 1", len(places))
	}
	if places[0].Name() != "Склад" || places[0].TotalVolume() != 100 {
		t.Errorf("storage place: name=%q volume=%d", places[0].Name(), places[0].TotalVolume())
	}
}

func Test_CourierCanTakeOrderReturnsFalseWhenNoStoragePlaces(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, err := c.CanTakeOrder(ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when courier has no storage places")
	}
}

func Test_CourierCanTakeOrderReturnsTrueWhenPlaceHasCapacity(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_ = c.AddStoragePlace("Склад", 100)
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, err := c.CanTakeOrder(ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true when storage place has capacity")
	}
	places := c.StoragePlaces()
	if places[0].OrderId() != ord.Id() {
		t.Error("storage place should contain order after CanTakeOrder")
	}
}

func Test_CourierCanTakeOrderReturnsFalseWhenVolumeExceedsCapacity(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_ = c.AddStoragePlace("Склад", 5)
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, err := c.CanTakeOrder(ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when order volume exceeds storage capacity")
	}
}

func Test_CourierCompleteOrderReturnsErrorWhenOrderNil(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	err := c.CompleteOrder(nil)
	if err == nil {
		t.Error("expected error when order is nil")
	}
}

func Test_CourierCompleteOrderReturnsErrorWhenOrderNotInStorage(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_ = c.AddStoragePlace("Склад", 100)
	ord, _ := orderModel.NewOrder(loc, 10)
	err := c.CompleteOrder(ord)
	if err == nil {
		t.Error("expected error when order not in courier storage")
	}
}

func Test_CourierCompleteOrderSucceedsAndUpdatesOrder(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_ = c.AddStoragePlace("Склад", 100)
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, _ := c.CanTakeOrder(ord)
	if !ok {
		t.Fatal("courier should be able to take order")
	}
	_ = ord.Assign(c.ID())
	err := c.CompleteOrder(ord)
	if err != nil {
		t.Fatalf("CompleteOrder failed: %v", err)
	}
	if ord.Status() != orderModel.OrderStatusCompleted {
		t.Errorf("order status = %s, want Completed", ord.Status())
	}
	places := c.StoragePlaces()
	if places[0].OrderId() != uuid.Nil {
		t.Error("storage place should be cleared after CompleteOrder")
	}
}

func Test_CourierCalculateTimeToLocationReturnsErrorWhenCourierLocationEmpty(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	c.location = kernel.Location{}
	_, err := c.CalculateTimeToLocation(loc)
	if err == nil {
		t.Error("expected error when courier location is empty")
	}
}

func Test_CourierCalculateTimeToLocationReturnsErrorWhenTargetEmpty(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_, err := c.CalculateTimeToLocation(kernel.Location{})
	if err == nil {
		t.Error("expected error when target location is empty")
	}
}

func Test_CourierCalculateTimeToLocationReturnsCorrectTime(t *testing.T) {
	loc1, _ := kernel.NewLocation(1, 1)
	loc2, _ := kernel.NewLocation(1, 2)
	c, _ := NewCourier("Курьер", 2, loc1)
	time, err := c.CalculateTimeToLocation(loc2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if time != 0.5 {
		t.Errorf("CalculateTimeToLocation() = %v, want 0.5 (distance=1, speed=2)", time)
	}
}

func Test_CourierMoveReturnsErrorWhenTargetEmpty(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	err := c.Move(kernel.Location{})
	if err == nil {
		t.Error("expected error when target is empty")
	}
}

func Test_CourierMoveUpdatesLocation(t *testing.T) {
	loc1, _ := kernel.NewLocation(1, 1)
	loc2, _ := kernel.NewLocation(5, 5)
	c, _ := NewCourier("Курьер", 5, loc1)
	err := c.Move(loc2)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}
	if c.Location() != loc2 {
		t.Error("Location() should equal target after Move")
	}
}

func Test_CourierFindStoragePlaceByOrderIdReturnsErrorWhenNoPlaces(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_, err := c.findStoragePlaceByOrderId(uuid.New())
	if err == nil {
		t.Error("expected error when courier has no storage places")
	}
}

func Test_CourierFindStoragePlaceByOrderIdReturnsErrorWhenOrderNotFound(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_ = c.AddStoragePlace("Склад", 100)
	_, err := c.findStoragePlaceByOrderId(uuid.New())
	if err == nil {
		t.Error("expected error when order id not in storage")
	}
}

func Test_CourierFindStoragePlaceByOrderIdReturnsPlaceWhenFound(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", 5, loc)
	_ = c.AddStoragePlace("Склад", 100)
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, _ := c.CanTakeOrder(ord)
	if !ok {
		t.Fatal("courier should take order")
	}
	place, err := c.findStoragePlaceByOrderId(ord.Id())
	if err != nil {
		t.Fatalf("findStoragePlaceByOrderId failed: %v", err)
	}
	if place.OrderId() != ord.Id() {
		t.Error("returned place should have the order id")
	}
}
