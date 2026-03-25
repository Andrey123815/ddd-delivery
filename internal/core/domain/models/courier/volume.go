package courier

import "errors"

type Volume struct {
	value int
}

func NewVolume(value int) (Volume, error) {
	if value <= 0 {
		return Volume{}, errors.New("объем должен быть больше нуля")
	}
	return Volume{value: value}, nil
}

func (v Volume) Value() int {
	return v.value
}

func (v Volume) Equals(other Volume) bool {
	return v.value == other.value
}

func (v Volume) CanFit(required Volume) bool {
	return v.value >= required.value
}
