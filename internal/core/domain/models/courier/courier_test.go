package courier

import (
	"testing"

	"delivery/internal/core/domain/models/kernel"

	orderModel "delivery/internal/core/domain/models/order"

	"github.com/google/uuid"
)

func mustNewSpeed(value int) Speed {
	speed, err := NewSpeed(value)
	if err != nil {
		return Speed{}
	}
	return speed
}

func mustNewVolume(value int) Volume {
	volume, err := NewVolume(value)
	if err != nil {
		return Volume{}
	}
	return volume
}

func Test_NewCourierCreatedAtValidParams(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, err := NewCourier("Курьер", mustNewSpeed(5), loc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.Name() != "Курьер" {
		t.Errorf("Name() = %q, want %q", c.Name(), "Курьер")
	}
	if c.Speed().Value() != 5 {
		t.Errorf("Speed() = %d, want 5", c.Speed().Value())
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
	_, err := NewCourier("", mustNewSpeed(5), loc)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func Test_NewCourierReturnsErrorAtInvalidSpeed(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	_, err := NewCourier("Курьер", mustNewSpeed(0), loc)
	if err == nil {
		t.Error("expected error for speed 0")
	}
	_, err = NewCourier("Курьер", mustNewSpeed(-1), loc)
	if err == nil {
		t.Error("expected error for negative speed")
	}
}

func Test_NewCourierReturnsErrorAtEmptyLocation(t *testing.T) {
	_, err := NewCourier("Курьер", mustNewSpeed(5), kernel.Location{})
	if err == nil {
		t.Error("expected error for empty location")
	}
}

func Test_CourierEqualsReturnsCorrectResult(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	a, _ := NewCourier("A", mustNewSpeed(5), loc)
	if !a.Equals(*a) {
		t.Error("courier should equal itself")
	}
	b, _ := NewCourier("B", mustNewSpeed(5), loc)
	if a.Equals(*b) {
		t.Error("different couriers should not be equal")
	}
}

func Test_CourierAddStoragePlaceSucceeds(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	err := c.AddStoragePlace("Склад", mustNewVolume(100))
	if err != nil {
		t.Fatalf("AddStoragePlace failed: %v", err)
	}
	places := c.StoragePlaces()
	if len(places) != 1 {
		t.Fatalf("len(StoragePlaces()) = %d, want 1", len(places))
	}
	if places[0].Name() != "Склад" || places[0].TotalVolume().Value() != 100 {
		t.Errorf("storage place: name=%q volume=%d", places[0].Name(), places[0].TotalVolume().Value())
	}
}

func Test_CourierCanTakeOrderReturnsFalseWhenNoStoragePlaces(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
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
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, err := c.CanTakeOrder(ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true when storage place has enough capacity")
	}
}

func Test_CourierCanTakeOrderReturnsFalseWhenVolumeExceedsCapacity(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(5))
	ord, _ := orderModel.NewOrder(loc, 10)
	ok, err := c.CanTakeOrder(ord)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when order volume exceeds storage capacity")
	}
}

func Test_CourierTakeOrderSucceedsWhenPlaceAvailable(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	ord, _ := orderModel.NewOrder(loc, 10)
	
	err := c.TakeOrder(ord)
	if err != nil {
		t.Fatalf("TakeOrder failed: %v", err)
	}
	
	place, err := c.findStoragePlaceByOrderId(ord.Id())
	if err != nil {
		t.Fatalf("order should be stored in storage place")
	}
	if place.OrderId() != ord.Id() {
		t.Error("storage place should have order ID set")
	}
}

func Test_CourierTakeOrderReturnsErrorWhenNoPlace(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	ord, _ := orderModel.NewOrder(loc, 10)
	
	err := c.TakeOrder(ord)
	if err == nil {
		t.Error("expected error when courier has no storage places")
	}
}

func Test_CourierTakeOrderReturnsErrorWhenVolumeExceeds(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(5))
	ord, _ := orderModel.NewOrder(loc, 10)
	
	err := c.TakeOrder(ord)
	if err == nil {
		t.Error("expected error when order volume exceeds capacity")
	}
}

func Test_CourierCompleteOrderReturnsErrorWhenOrderNil(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	err := c.CompleteOrder(uuid.Nil)
	if err == nil {
		t.Error("expected error when orderId is nil")
	}
}

func Test_CourierCompleteOrderReturnsErrorWhenOrderNotInStorage(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	orderId := uuid.New()
	err := c.CompleteOrder(orderId)
	if err == nil {
		t.Error("expected error when order not in courier storage")
	}
}

func Test_CourierCompleteOrderSucceedsAndClearsStorage(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	ord, _ := orderModel.NewOrder(loc, 10)

	err := c.TakeOrder(ord)
	if err != nil {
		t.Fatalf("TakeOrder failed: %v", err)
	}

	err = c.CompleteOrder(ord.Id())
	if err != nil {
		t.Fatalf("CompleteOrder failed: %v", err)
	}
	
	places := c.StoragePlaces()
	if places[0].OrderId() != uuid.Nil {
		t.Error("storage place should be cleared after CompleteOrder")
	}
}

func Test_CourierCalculateTimeToLocationReturnsErrorWhenCourierLocationEmpty(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	c.location = kernel.Location{}
	_, err := c.CalculateTimeToLocation(loc)
	if err == nil {
		t.Error("expected error when courier location is empty")
	}
}

func Test_CourierCalculateTimeToLocationReturnsErrorWhenTargetEmpty(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_, err := c.CalculateTimeToLocation(kernel.Location{})
	if err == nil {
		t.Error("expected error when target location is empty")
	}
}

func Test_CourierCalculateTimeToLocationReturnsCorrectTime(t *testing.T) {
	loc1, _ := kernel.NewLocation(1, 1)
	loc2, _ := kernel.NewLocation(1, 2)
	c, _ := NewCourier("Курьер", mustNewSpeed(2), loc1)
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
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	err := c.Move(kernel.Location{})
	if err == nil {
		t.Error("expected error when target is empty")
	}
}

func Test_CourierMoveWithSpeedTowardsTarget(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	target, _ := kernel.NewLocation(10, 10)
	
	err := c.Move(target)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}
	
	if c.Location().X() != 6 || c.Location().Y() != 1 {
		t.Errorf("expected move by speed=5 towards (10,10) to be (6,1), got (%d,%d)", c.Location().X(), c.Location().Y())
	}
}

func Test_CourierMoveDistributesSpeedBetweenAxes(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(3), loc)
	target, _ := kernel.NewLocation(3, 4)
	
	err := c.Move(target)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}
	
	if c.Location().X() != 3 || c.Location().Y() != 2 {
		t.Errorf("expected (3,2) after moving speed=3 steps, got (%d,%d)", c.Location().X(), c.Location().Y())
	}
}

func Test_CourierMoveDoesNothingWhenAlreadyAtTarget(t *testing.T) {
	loc, _ := kernel.NewLocation(3, 3)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	
	err := c.Move(loc)
	if err != nil {
		t.Fatalf("Move failed: %v", err)
	}
	
	if c.Location().X() != 3 || c.Location().Y() != 3 {
		t.Errorf("expected location to stay (3,3), got (%d,%d)", c.Location().X(), c.Location().Y())
	}
}

func Test_CourierMoveMultipleStepsReachesTarget(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	target, _ := kernel.NewLocation(4, 3)
	
	for !c.Location().Equals(target) {
		err := c.Move(target)
		if err != nil {
			t.Fatalf("Move failed: %v", err)
		}
	}
	
	if c.Location().X() != 4 || c.Location().Y() != 3 {
		t.Errorf("expected to reach (4,3), got (%d,%d)", c.Location().X(), c.Location().Y())
	}
}

func Test_CourierFindStoragePlaceByOrderIdReturnsErrorWhenNoPlaces(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_, err := c.findStoragePlaceByOrderId(uuid.New())
	if err == nil {
		t.Error("expected error when courier has no storage places")
	}
}

func Test_CourierFindStoragePlaceByOrderIdReturnsErrorWhenOrderNotFound(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	_, err := c.findStoragePlaceByOrderId(uuid.New())
	if err == nil {
		t.Error("expected error when order id not in storage")
	}
}

func Test_CourierFindStoragePlaceByOrderIdReturnsPlaceWhenFound(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	c, _ := NewCourier("Курьер", mustNewSpeed(5), loc)
	_ = c.AddStoragePlace("Склад", mustNewVolume(100))
	ord, _ := orderModel.NewOrder(loc, 10)

	err := c.TakeOrder(ord)
	if err != nil {
		t.Fatalf("TakeOrder failed: %v", err)
	}

	place, err := c.findStoragePlaceByOrderId(ord.Id())
	if err != nil {
		t.Fatalf("findStoragePlaceByOrderId failed: %v", err)
	}
	if place.OrderId() != ord.Id() {
		t.Error("returned place should have the order id")
	}
}
