package courier

import "testing"

func Test_NewSpeedValidAtValidParams(t *testing.T) {
	speed, err := NewSpeed(5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if speed.Value() != 5 {
		t.Errorf("Value() = %d, want 5", speed.Value())
	}
}

func Test_NewSpeedReturnsErrorAtZero(t *testing.T) {
	_, err := NewSpeed(0)
	if err == nil {
		t.Error("expected error for speed = 0")
	}
}

func Test_NewSpeedReturnsErrorAtNegative(t *testing.T) {
	_, err := NewSpeed(-1)
	if err == nil {
		t.Error("expected error for negative speed")
	}
	_, err = NewSpeed(-100)
	if err == nil {
		t.Error("expected error for negative speed")
	}
}

func Test_SpeedEquals(t *testing.T) {
	speed1, _ := NewSpeed(5)
	speed2, _ := NewSpeed(5)
	speed3, _ := NewSpeed(10)
	
	if !speed1.Equals(speed2) {
		t.Error("speeds with same value should be equal")
	}
	
	if speed1.Equals(speed3) {
		t.Error("speeds with different values should not be equal")
	}
}
