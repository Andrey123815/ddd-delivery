package courier

import "testing"

func Test_NewVolumeValidAtValidParams(t *testing.T) {
	volume, err := NewVolume(50)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if volume.Value() != 50 {
		t.Errorf("Value() = %d, want 50", volume.Value())
	}
}

func Test_NewVolumeReturnsErrorAtZero(t *testing.T) {
	_, err := NewVolume(0)
	if err == nil {
		t.Error("expected error for volume = 0")
	}
}

func Test_NewVolumeReturnsErrorAtNegative(t *testing.T) {
	_, err := NewVolume(-1)
	if err == nil {
		t.Error("expected error for negative volume")
	}
	_, err = NewVolume(-100)
	if err == nil {
		t.Error("expected error for negative volume")
	}
}

func Test_VolumeEquals(t *testing.T) {
	volume1, _ := NewVolume(50)
	volume2, _ := NewVolume(50)
	volume3, _ := NewVolume(100)
	
	if !volume1.Equals(volume2) {
		t.Error("volumes with same value should be equal")
	}
	
	if volume1.Equals(volume3) {
		t.Error("volumes with different values should not be equal")
	}
}

func Test_VolumeCanFit(t *testing.T) {
	large, _ := NewVolume(100)
	small, _ := NewVolume(50)
	exact, _ := NewVolume(100)
	
	if !large.CanFit(small) {
		t.Error("larger volume should fit smaller")
	}
	
	if !large.CanFit(exact) {
		t.Error("volume should fit exact match")
	}
	
	if small.CanFit(large) {
		t.Error("smaller volume should not fit larger")
	}
}
