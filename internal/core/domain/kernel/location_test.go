package kernel

import (
	"testing"
)

func Test_NewLocationValidAtValidParams(t *testing.T) {
	loc, err := NewLocation(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if loc.X() != 1 || loc.Y() != 1 {
		t.Errorf("expected (1,1), got (%d,%d)", loc.X(), loc.Y())
	}

	loc, err = NewLocation(10, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if loc.X() != 10 || loc.Y() != 10 {
		t.Errorf("expected (10,10), got (%d,%d)", loc.X(), loc.Y())
	}

	loc, err = NewLocation(5, 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if loc.X() != 5 || loc.Y() != 7 {
		t.Errorf("expected (5,7), got (%d,%d)", loc.X(), loc.Y())
	}
}

func Test_NewLocationCreateErrorAtInvalidX(t *testing.T) {
	_, err := NewLocation(0, 5)
	if err == nil {
		t.Error("expected error for x=0")
	}

	_, err = NewLocation(11, 5)
	if err == nil {
		t.Error("expected error for x=11")
	}
}

func Test_NewLocationCreateErrorAtInvalidY(t *testing.T) {
	_, err := NewLocation(5, 0)
	if err == nil {
		t.Error("expected error for y=0")
	}

	_, err = NewLocation(5, 11)
	if err == nil {
		t.Error("expected error for y=11")
	}
}

func Test_NewRandomLocationGeneratedInValidRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		loc, err := NewRandomLocation()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if loc.X() < 1 || loc.X() > 10 {
			t.Errorf("x=%d out of range [1,10]", loc.X())
		}
		if loc.Y() < 1 || loc.Y() > 10 {
			t.Errorf("y=%d out of range [1,10]", loc.Y())
		}
	}
}

func TestLocation_X_Y(t *testing.T) {
	loc, _ := NewLocation(3, 4)
	if loc.X() != 3 {
		t.Errorf("X() = %d, want 3", loc.X())
	}
	if loc.Y() != 4 {
		t.Errorf("Y() = %d, want 4", loc.Y())
	}
}

func Test_LocationEqualsReturnsCorrectResult(t *testing.T) {
	a, _ := NewLocation(2, 3)
	b, _ := NewLocation(2, 3)
	c, _ := NewLocation(2, 4)

	if !a.Equals(*b) {
		t.Error("same coords should be equal")
	}
	if a.Equals(*c) {
		t.Error("different coords should not be equal")
	}
}

func Test_AbsInt(t *testing.T) {
	if got := AbsInt(5); got != 5 {
		t.Errorf("AbsInt(5) = %d, want 5", got)
	}
	if got := AbsInt(-5); got != 5 {
		t.Errorf("AbsInt(-5) = %d, want 5", got)
	}
	if got := AbsInt(0); got != 0 {
		t.Errorf("AbsInt(0) = %d, want 0", got)
	}
}

func Test_LocationDistanceCalculatedCorrectly(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 uint8
		want           uint8
	}{
		{1, 1, 1, 1, 0},
		{1, 1, 1, 2, 1},
		{1, 1, 2, 1, 1},
		{1, 1, 2, 2, 2},
		{2, 2, 1, 1, 2},
		{1, 10, 10, 1, 18},
	}
	for _, tt := range tests {
		a, _ := NewLocation(tt.x1, tt.y1)
		b, _ := NewLocation(tt.x2, tt.y2)
		if got := a.Distance(*b); got != tt.want {
			t.Errorf("Distance((%d,%d), (%d,%d)) = %d, want %d", tt.x1, tt.y1, tt.x2, tt.y2, got, tt.want)
		}
	}
}
