package courier

import (
	"testing"

	"github.com/google/uuid"
)

func Test_NewStoragePlaceValidAtValidParams(t *testing.T) {
	place, err := NewStoragePlace("Склад A", 100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if place.Name() != "Склад A" {
		t.Errorf("Name() = %q, want %q", place.Name(), "Склад A")
	}

	if place.TotalVolume() != 100 {
		t.Errorf("TotalVolume() = %d, want 100", place.TotalVolume())
	}

	if place.Id() == uuid.Nil {
		t.Error("Id() should not be Nil")
	}

	if place.OrderId() != uuid.Nil {
		t.Error("new place should have no order")
	}
}

func Test_NewStoragePlaceCreateErrorAtInvalidName(t *testing.T) {
	_, err := NewStoragePlace("", 100)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func Test_NewStoragePlaceCreateErrorAtInvalidVolume(t *testing.T) {
	_, err := NewStoragePlace("Склад", 0)
	if err == nil {
		t.Error("expected error for volume 0")
	}

	_, err = NewStoragePlace("Склад", -5)
	if err == nil {
		t.Error("expected error for negative volume")
	}
}

func Test_StoragePlaceEqualsReturnsCorrectResult(t *testing.T) {
	a, _ := NewStoragePlace("A", 10)
	if !a.Equals(a) {
		t.Error("place should equal itself")
	}

	b, _ := NewStoragePlace("B", 10)
	if a.Equals(b) {
		t.Error("different places should not be equal")
	}
}

func Test_StoragePlaceCanStoreReturnsErrorAtInvalidVolume(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)

	_, err := place.CanStore(0)
	if err == nil {
		t.Error("expected error for volume 0")
	}

	_, err = place.CanStore(-1)
	if err == nil {
		t.Error("expected error for negative volume")
	}
}

func Test_StoragePlaceCanStoreReturnsTrueWhenFreeAndEnoughVolume(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)

	ok, err := place.CanStore(30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !ok {
		t.Error("free place with enough volume should allow store")
	}

	ok, _ = place.CanStore(50)
	if !ok {
		t.Error("free place with exact volume should allow store")
	}
	
	ok, _ = place.CanStore(51)
	if ok {
		t.Error("volume > totalVolume should not allow store")
	}
}

func Test_StoragePlaceCanStoreReturnsFalseWhenOccupied(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)
	orderID := uuid.New()

	if err := place.Store(orderID, 20); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	ok, err := place.CanStore(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ok {
		t.Error("occupied place should not allow store")
	}
}

func Test_StoragePlaceStoreSucceedsAtValidParams(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)
	orderID := uuid.New()

	err := place.Store(orderID, 30)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	if place.OrderId() != orderID {
		t.Errorf("OrderId() = %v, want %v", place.OrderId(), orderID)
	}
}

func Test_StoragePlaceStoreReturnsErrorAtInvalidVolume(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)

	err := place.Store(uuid.New(), 0)
	if err == nil {
		t.Error("expected error for volume 0")
	}
}

func Test_StoragePlaceStoreReturnsErrorWhenNoCapacity(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)

	err := place.Store(uuid.New(), 60)
	if err == nil {
		t.Error("expected error when volume exceeds capacity")
	}
}

func Test_StoragePlaceStoreReturnsErrorWhenAlreadyOccupied(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)
	_ = place.Store(uuid.New(), 20)

	err := place.Store(uuid.New(), 10)
	if err == nil {
		t.Error("expected error when place already occupied")
	}
}

func Test_StoragePlaceClearSucceedsWhenOccupiedWithMatchingOrderId(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)
	orderID := uuid.New()
	_ = place.Store(orderID, 20)

	err := place.Clear(orderID)
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if place.OrderId() != uuid.Nil {
		t.Error("after Clear, OrderId should be Nil")
	}
}

func Test_StoragePlaceClearReturnsErrorWhenEmpty(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)

	err := place.Clear(uuid.New())
	if err == nil {
		t.Error("expected error when clearing empty place")
	}
}

func Test_StoragePlaceClearReturnsErrorAtWrongOrderId(t *testing.T) {
	place, _ := NewStoragePlace("Склад", 50)
	_ = place.Store(uuid.New(), 20)
	
	err := place.Clear(uuid.New()) // другой ID
	if err == nil {
		t.Error("expected error when clearing with wrong order id")
	}
}
