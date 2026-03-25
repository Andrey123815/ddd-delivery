package kernel

import "testing"

func Test_DistanceEquals(t *testing.T) {
	loc1, _ := NewLocation(1, 1)
	loc2, _ := NewLocation(1, 1)
	loc3, _ := NewLocation(5, 5)
	
	distance1, _ := NewDistance(loc1, loc3)
	distance2, _ := NewDistance(loc2, loc3)
	distance3, _ := NewDistance(loc1, loc2)
	
	if !distance1.Equals(distance2) {
		t.Error("distances with same value should be equal")
	}
	
	if distance1.Equals(distance3) {
		t.Error("distances with different values should not be equal")
	}
}

func Test_DistanceIsZero(t *testing.T) {
	sameLoc, _ := NewLocation(5, 5)
	zero, _ := NewDistance(sameLoc, sameLoc)
	if !zero.IsZero() {
		t.Error("distance between same locations should be zero")
	}
	
	loc1, _ := NewLocation(1, 1)
	loc2, _ := NewLocation(5, 5)
	nonZero, _ := NewDistance(loc1, loc2)
	if nonZero.IsZero() {
		t.Error("distance between different locations should not be zero")
	}
}

func Test_NewDistance(t *testing.T) {
	tests := []struct {
		name       string
		x1, y1     uint8
		x2, y2     uint8
		wantValue  int
	}{
		{"same point", 1, 1, 1, 1, 0},
		{"horizontal", 1, 1, 3, 1, 2},
		{"vertical", 1, 1, 1, 4, 3},
		{"diagonal", 1, 1, 3, 4, 5},
		{"reverse", 5, 5, 2, 2, 6},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc1, _ := NewLocation(tt.x1, tt.y1)
			loc2, _ := NewLocation(tt.x2, tt.y2)
			distance, err := NewDistance(loc1, loc2)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if distance.Value() != tt.wantValue {
				t.Errorf("distance = %d, want %d", distance.Value(), tt.wantValue)
			}
		})
	}
}

func Test_NewDistanceReturnsErrorWhenLocationEmpty(t *testing.T) {
	loc, _ := NewLocation(1, 1)
	empty := Location{}
	
	_, err := NewDistance(loc, empty)
	if err == nil {
		t.Error("expected error when target location is empty")
	}
	
	_, err = NewDistance(empty, loc)
	if err == nil {
		t.Error("expected error when source location is empty")
	}
}
