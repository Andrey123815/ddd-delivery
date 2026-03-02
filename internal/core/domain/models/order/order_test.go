package order

import (
	"testing"

	"delivery/internal/core/domain/kernel"
	"github.com/google/uuid"
)

func Test_NewOrderCreatedAtValidParams(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)

	order, err := NewOrder(loc, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.Location() != loc {
		t.Errorf("Location() = %+v, want %+v", order.Location(), loc)
	}
	if order.Volume() != 10 {
		t.Errorf("Volume() = %d, want 10", order.Volume())
	}
	if order.Status() != OrderStatusCreated {
		t.Errorf("Status() = %s, want %s", order.Status(), OrderStatusCreated)
	}
	if order.CourierId() != uuid.Nil {
		t.Errorf("CourierId() = %v, want uuid.Nil", order.CourierId())
	}
}

func Test_NewOrderCreateErrorAtEmptyLocation(t *testing.T) {
	emptyLoc := kernel.Location{}
	_, err := NewOrder(emptyLoc, 10)
	if err == nil {
		t.Error("expected error when creating order with empty location")
	}
}

func Test_NewOrderCreateErrorAtInvalidVolume(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)

	_, err := NewOrder(loc, 0)
	if err == nil {
		t.Error("expected error for volume 0")
	}

	_, err = NewOrder(loc, -5)
	if err == nil {
		t.Error("expected error for negative volume")
	}
}

func Test_OrderEqualsReturnsCorrectResult(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)

	a, _ := NewOrder(loc, 10)
	if !a.Equals(*a) {
		t.Error("order should equal itself")
	}

	b, _ := NewOrder(loc, 10)
	if a.Equals(*b) {
		t.Error("different orders should not be equal")
	}
}

func Test_OrderAssignSetsStatusAndCourierIdAtValidParams(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)
	courierId := uuid.New()

	if err := order.Assign(courierId); err != nil {
		t.Fatalf("Assign failed: %v", err)
	}
	if order.Status() != OrderStatusAssigned {
		t.Errorf("Status() = %s, want %s", order.Status(), OrderStatusAssigned)
	}
	if order.CourierId() != courierId {
		t.Errorf("CourierId() = %v, want %v", order.CourierId(), courierId)
	}
}

func Test_OrderAssignReturnsErrorWhenAlreadyAssigned(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)
	courierId := uuid.New()

	_ = order.Assign(courierId)
	if err := order.Assign(courierId); err == nil {
		t.Error("expected error when assigning already assigned order")
	}
}

func Test_OrderAssignReturnsErrorWhenCompleted(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)
	orderStatusCompleted := OrderStatusCompleted
	order.status = orderStatusCompleted

	if err := order.Assign(uuid.New()); err == nil {
		t.Error("expected error when assigning completed order")
	}
}

func Test_OrderAssignReturnsErrorWhenCourierIdNil(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)

	if err := order.Assign(uuid.Nil); err == nil {
		t.Error("expected error when assigning with Nil courier id")
	}
}

func Test_OrderCompleteReturnsErrorWhenCreated(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)

	if err := order.Complete(uuid.New()); err == nil {
		t.Error("expected error when completing unassigned order")
	}
}

func Test_OrderCompleteReturnsErrorWhenCompleted(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)
	courierId := uuid.New()

	_ = order.Assign(courierId)
	_ = order.Complete(courierId)

	if err := order.Complete(courierId); err == nil {
		t.Error("expected error when completing already completed order")
	}
}

func Test_OrderCompleteReturnsErrorWhenCourierDiffers(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)
	courierId := uuid.New()
	otherCourierId := uuid.New()

	_ = order.Assign(courierId)

	if err := order.Complete(otherCourierId); err == nil {
		t.Error("expected error when completing with different courier id")
	}
}

func Test_OrderCompleteSucceedsAtValidParams(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)
	order, _ := NewOrder(loc, 10)
	courierId := uuid.New()

	if err := order.Assign(courierId); err != nil {
		t.Fatalf("Assign failed: %v", err)
	}
	if err := order.Complete(courierId); err != nil {
		t.Fatalf("Complete failed: %v", err)
	}
	if order.Status() != OrderStatusCompleted {
		t.Errorf("Status() = %s, want %s", order.Status(), OrderStatusCompleted)
	}
}

