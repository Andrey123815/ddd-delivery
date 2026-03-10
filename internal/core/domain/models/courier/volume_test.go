package courier

import (
	"testing"
)

func Test_NewVolumeCreatedAtValidParams(t *testing.T) {
	v, err := NewVolume(10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if v.Volume() != 10 {
		t.Errorf("Volume() = %d, want 10", v.Volume())
	}

	v, err = NewVolume(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if v.Volume() != 1 {
		t.Errorf("Volume() = %d, want 1", v.Volume())
	}
}

func Test_NewVolumeReturnsErrorAtZeroOrNegative(t *testing.T) {
	_, err := NewVolume(0)
	if err == nil {
		t.Error("expected error for volume 0")
	}

	_, err = NewVolume(-5)
	if err == nil {
		t.Error("expected error for negative volume")
	}
}

func Test_VolumeEqualsReturnsCorrectResult(t *testing.T) {
	a, _ := NewVolume(5)
	b, _ := NewVolume(5)
	c, _ := NewVolume(10)

	if !a.Equals(b) {
		t.Error("same volume should be equal")
	}
	if a.Equals(c) {
		t.Error("different volumes should not be equal")
	}
	if !a.Equals(a) {
		t.Error("volume should equal itself")
	}
}
