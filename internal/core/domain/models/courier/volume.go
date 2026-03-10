package courier

import (
	"errors"
)

type Volume struct {
	volume int
}

func NewVolume(volume int) (*Volume, error) {
	if volume <= 0 {
		return nil, errors.New("Объем не может быть меньше или равным 0")
	}

	return &Volume{
		volume: volume,
	}, nil
}

func (v *Volume) Equals(other *Volume) bool {
	return v.volume == other.volume
}

func (v *Volume) Volume() int {
	return v.volume
}

func (v *Volume) GreaterThanOrEqual(other *Volume) bool {
	return v.volume >= other.volume
}
